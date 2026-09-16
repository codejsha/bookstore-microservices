package support

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/logging"

	"github.com/codejsha/bookstore-microservices/customer/generated/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/customer/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/httpx"
)

const (
	readyCheckTimeout = 3 * time.Second
	requestTimeout    = 15 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = requestTimeout
	writeTimeout      = requestTimeout + 5*time.Second
	idleTimeout       = 60 * time.Second
)

type GinServer struct {
	engine      *gin.Engine
	server      *http.Server
	shutdowner  fx.Shutdowner
	serverCfg   *config.ServerConfig
	logHelper   *logging.LogHelper
	readiness   *ReadinessDataSource
	readyCheck  func(ctx context.Context) error
	draining    atomic.Bool
	customerAPI openapi.CustomerApi
	orderAPI    openapi.OrderApi
	paymentAPI  openapi.PaymentApi
	deliveryAPI openapi.DeliveryApi
	pointAPI    openapi.PointApi
	reviewAPI   openapi.ReviewApi
	wishlistAPI openapi.WishlistApi
}

func NewGinServer(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	serverCfg *config.ServerConfig,
	logHelper *logging.LogHelper,
	readiness *ReadinessDataSource,
	customerAPI openapi.CustomerApi,
	orderAPI openapi.OrderApi,
	paymentAPI openapi.PaymentApi,
	deliveryAPI openapi.DeliveryApi,
	pointAPI openapi.PointApi,
	reviewAPI openapi.ReviewApi,
	wishlistAPI openapi.WishlistApi,
) *GinServer {
	s := &GinServer{
		shutdowner:  shutdowner,
		serverCfg:   serverCfg,
		logHelper:   logHelper,
		readiness:   readiness,
		customerAPI: customerAPI,
		orderAPI:    orderAPI,
		paymentAPI:  paymentAPI,
		deliveryAPI: deliveryAPI,
		pointAPI:    pointAPI,
		reviewAPI:   reviewAPI,
		wishlistAPI: wishlistAPI,
	}
	s.readyCheck = func(ctx context.Context) error {
		return s.readiness.Ping(ctx)
	}
	s.InitializeEngine()
	s.RegisterRoutes()

	addr := net.JoinHostPort(s.serverCfg.Host, s.serverCfg.Port)
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.engine,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen http on %s: %w", addr, err)
			}
			go s.serve(listener)
			return nil
		},
	})

	return s
}

func (s *GinServer) BeginDrain() {
	s.draining.Store(true)
}

func (s *GinServer) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err == nil || ctx.Err() == nil {
		return err
	}
	if closeErr := s.server.Close(); closeErr != nil {
		return errors.Join(err, closeErr)
	}
	return err
}

func (s *GinServer) serve(listener net.Listener) {
	if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Errorf("http server stopped: %v", err)
		s.signalShutdown()
	}
}

func (s *GinServer) signalShutdown() {
	if s.shutdowner == nil {
		return
	}
	if err := s.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
		logrus.Errorf("failed to signal application shutdown: %v", err)
	}
}

func (s *GinServer) InitializeEngine() {
	gin.SetMode(s.serverCfg.Mode)
	gin.DefaultWriter = s.logHelper.Logger.WriterLevel(logrus.InfoLevel)
	gin.DefaultErrorWriter = s.logHelper.Logger.WriterLevel(logrus.ErrorLevel)
	RegisterBindingTagNames()
	s.engine = gin.New()
	s.engine.Use(GinAccessLogMiddleware(s.logHelper))
	s.engine.Use(gin.Recovery())
	s.engine.Use(otelgin.Middleware(string(constant.TracerNameGinServer)))
	s.engine.Use(GinRequestDeadlineMiddleware(requestTimeout))
	s.engine.Use(httpx.GinResponseMapping())
	s.engine.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	s.engine.GET("/health/ready", func(c *gin.Context) {
		if s.draining.Load() {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), readyCheckTimeout)
		defer cancel()
		if err := s.readyCheck(ctx); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	})
	s.engine.Use(GinPrincipalMiddleware())
	s.engine.Use(GinOwnershipMiddleware())
}

func GinRequestDeadlineMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinAccessLogMiddleware(logHelper *logging.LogHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		if strings.HasPrefix(c.Request.URL.Path, "/health") && status < http.StatusBadRequest {
			return
		}

		fields := logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"latency_us": time.Since(start).Microseconds(),
			"client_ip":  c.ClientIP(),
		}
		if raw := c.Request.URL.RawQuery; raw != "" {
			fields["query"] = raw
		}
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}
		if sc := trace.SpanContextFromContext(c.Request.Context()); sc.HasTraceID() {
			fields["trace_id"] = sc.TraceID().String()
			fields["span_id"] = sc.SpanID().String()
		}
		if p := PrincipalFromContext(c); p != nil {
			fields["user_uid"] = p.Sub
		}

		entry := logHelper.Logger.WithFields(fields)
		switch {
		case status >= http.StatusInternalServerError:
			entry.Error()
		case status >= http.StatusBadRequest:
			entry.Warn()
		default:
			entry.Info()
		}
	}
}

func (s *GinServer) RegisterRoutes() {
	openapi.NewCustomerApiHandler(s.customerAPI).RegisterRoutes(s.engine)
	openapi.NewOrderApiHandler(s.orderAPI).RegisterRoutes(s.engine)
	openapi.NewPaymentApiHandler(s.paymentAPI).RegisterRoutes(s.engine)
	openapi.NewDeliveryApiHandler(s.deliveryAPI).RegisterRoutes(s.engine)
	openapi.NewPointApiHandler(s.pointAPI).RegisterRoutes(s.engine)
	openapi.NewReviewApiHandler(s.reviewAPI).RegisterRoutes(s.engine)
	openapi.NewWishlistApiHandler(s.wishlistAPI).RegisterRoutes(s.engine)
}

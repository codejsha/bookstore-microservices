package support

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
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

type GinServer struct {
	engine      *gin.Engine
	server      *http.Server
	serverCfg   *config.ServerConfig
	logHelper   *logging.LogHelper
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
	serverCfg *config.ServerConfig,
	logHelper *logging.LogHelper,
	customerAPI openapi.CustomerApi,
	orderAPI openapi.OrderApi,
	paymentAPI openapi.PaymentApi,
	deliveryAPI openapi.DeliveryApi,
	pointAPI openapi.PointApi,
	reviewAPI openapi.ReviewApi,
	wishlistAPI openapi.WishlistApi,
) *GinServer {
	s := &GinServer{
		serverCfg:   serverCfg,
		logHelper:   logHelper,
		customerAPI: customerAPI,
		orderAPI:    orderAPI,
		paymentAPI:  paymentAPI,
		deliveryAPI: deliveryAPI,
		pointAPI:    pointAPI,
		reviewAPI:   reviewAPI,
		wishlistAPI: wishlistAPI,
	}
	s.InitializeEngine()
	s.RegisterRoutes()

	addr := net.JoinHostPort(s.serverCfg.Host, s.serverCfg.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen http on %s: %w", addr, err)
			}
			go func() {
				if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logrus.Errorf("http server stopped: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.server.Shutdown(ctx)
		},
	})

	return s
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
	s.engine.Use(httpx.GinResponseMapping())
	s.engine.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	s.engine.Use(GinPrincipalMiddleware())
	s.engine.Use(GinOwnershipMiddleware())
}

func GinAccessLogMiddleware(logHelper *logging.LogHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		if c.Request.URL.Path == "/health" && status < http.StatusBadRequest {
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

package support

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/customer/internal/application/port/openapi"
)

func NewGinServer(
	serverCfg *config.ServerConfig,
	logHelper *logging.LogHelper,
	customerAPI openapi.CustomerAPI,
	orderAPI openapi.OrderAPI,
	paymentAPI openapi.PaymentAPI,
	pointAPI openapi.PointAPI,
	wishlistAPI openapi.WishlistAPI,
) *GinServer {
	server := &GinServer{
		serverCfg:   serverCfg,
		logHelper:   logHelper,
		customerAPI: customerAPI,
		orderAPI:    orderAPI,
		paymentAPI:  paymentAPI,
		pointAPI:    pointAPI,
		wishlistAPI: wishlistAPI,
	}
	server.InitializeEngine()
	server.RegisterRoutes()
	return server
}

type GinServer struct {
	engine      *gin.Engine
	serverCfg   *config.ServerConfig
	logHelper   *logging.LogHelper
	customerAPI openapi.CustomerAPI
	orderAPI    openapi.OrderAPI
	paymentAPI  openapi.PaymentAPI
	pointAPI    openapi.PointAPI
	wishlistAPI openapi.WishlistAPI
}

func (s *GinServer) InitializeEngine() {
	gin.SetMode(s.serverCfg.Mode)
	gin.DefaultWriter = s.logHelper.Logger.WriterLevel(logrus.InfoLevel)
	gin.DefaultErrorWriter = s.logHelper.Logger.WriterLevel(logrus.ErrorLevel)
	s.engine = gin.Default()
}

func (s *GinServer) Run() {
	addr := net.JoinHostPort(s.serverCfg.Host, s.serverCfg.Port)
	err := s.engine.Run(addr)
	if err != nil {
		panic(err)
	}
}

func (s *GinServer) RegisterRoutes() {
	api := s.engine.Group("")
	{
		api.GET("/api/v1/customers", s.customerAPI.ApiV1CustomersGet)
		api.GET("/api/v1/customers/:id", s.customerAPI.ApiV1CustomersIdGet)
		api.PUT("/api/v1/customers/:id", s.customerAPI.ApiV1CustomersIdPut)
		api.DELETE("/api/v1/customers/:id", s.customerAPI.ApiV1CustomersIdDelete)

		api.GET("/api/v1/customers/wishlist", s.wishlistAPI.ApiV1CustomersIdWishlistGet)
		api.POST("/api/v1/customers/wishlist", s.wishlistAPI.ApiV1CustomersIdWishlistPost)
		api.PUT("/api/v1/customers/wishlist", s.wishlistAPI.ApiV1CustomersIdWishlistPut)
		api.DELETE("/api/v1/customers/wishlist", s.wishlistAPI.ApiV1CustomersIdWishlistDelete)

		api.GET("/api/v1/customers/points", s.pointAPI.ApiV1CustomersIdPointGet)
		api.PUT("/api/v1/customers/points", s.pointAPI.ApiV1CustomersIdPointPut)

		api.GET("/api/v1/customers/orders", s.orderAPI.ApiV1CustomersIdOrdersGet)
		api.GET("/api/v1/customers/orders/:id", s.orderAPI.ApiV1CustomersIdOrdersOrderIdGet)

		api.GET("/api/v1/customers/payments", s.paymentAPI.ApiV1CustomersIdPaymentsGet)
		api.GET("/api/v1/customers/payments/:id", s.paymentAPI.ApiV1CustomersIdPaymentsPaymentIdGet)
	}
}

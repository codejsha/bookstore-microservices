package support

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/inventory/internal/application/port/openapi"
)

func NewGinServer(
	serverCfg *config.ServerConfig,
	logHelper *logging.LogHelper,
	stocksAPI openapi.StocksAPI,
	warehousesAPI openapi.WarehouseAPI,
) *GinServer {
	server := &GinServer{
		serverCfg:     serverCfg,
		logHelper:     logHelper,
		stocksAPI:     stocksAPI,
		warehousesAPI: warehousesAPI,
	}
	server.InitializeEngine()
	server.RegisterRoutes()
	return server
}

type GinServer struct {
	engine        *gin.Engine
	serverCfg     *config.ServerConfig
	logHelper     *logging.LogHelper
	stocksAPI     openapi.StocksAPI
	warehousesAPI openapi.WarehouseAPI
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
		api.GET("/api/v1/stocks", s.stocksAPI.ApiV1StocksGet)
		api.POST("/api/v1/stocks", s.stocksAPI.ApiV1StocksPost)
		api.GET("/api/v1/stocks/:id", s.stocksAPI.ApiV1StocksIdGet)
		api.PUT("/api/v1/stocks/:id", s.stocksAPI.ApiV1StocksIdPut)
		api.DELETE("/api/v1/stocks/:id", s.stocksAPI.ApiV1StocksIdDelete)

		api.GET("/api/v1/warehouses", s.warehousesAPI.ApiV1WarehousesGet)
		api.POST("/api/v1/warehouses", s.warehousesAPI.ApiV1WarehousesPost)
		api.GET("/api/v1/warehouses/:id", s.warehousesAPI.ApiV1WarehousesIdGet)
		api.PUT("/api/v1/warehouses/:id", s.warehousesAPI.ApiV1WarehousesIdPut)
		api.DELETE("/api/v1/warehouses/:id", s.warehousesAPI.ApiV1WarehousesIdDelete)
	}
}

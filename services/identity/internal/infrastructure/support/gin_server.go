package support

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/openapi"
)

func NewGinServer(
	serverCfg *config.ServerConfig,
	logHelper *logging.LogHelper,
	identityAPI openapi.IdentityAPI,
) *GinServer {
	server := &GinServer{
		serverCfg:          serverCfg,
		logHelper:          logHelper,
		identityController: identityAPI,
	}
	server.InitializeEngine()
	server.RegisterRoutes()
	return server
}

type GinServer struct {
	engine             *gin.Engine
	serverCfg          *config.ServerConfig
	logHelper          *logging.LogHelper
	identityController openapi.IdentityAPI
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
		api.POST("/api/v1/identity/signin", s.identityController.ApiV1UsersSigninPost)
		api.POST("/api/v1/identity/signout", s.identityController.ApiV1UsersSignoutPost)
		api.POST("/api/v1/identity/signup", s.identityController.ApiV1UsersSignoutPost)
		api.POST("/api/v1/identity/token/exchange", s.identityController.ApiV1UsersTokensExchangePost)
	}
}

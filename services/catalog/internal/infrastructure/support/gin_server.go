package support

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/codejsha/bookstore-microservices/catalog/internal/application/port/openapi"
	"github.com/codejsha/bookstore-microservices/catalog/internal/domain/constant"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/logging"
)

func NewGinServer(
	serverCfg *config.ServerConfig,
	logHelper *logging.LogHelper,
	authorAPI openapi.AuthorAPI,
	bookAPI openapi.BookAPI,
	publisherAPI openapi.PublisherAPI,
	categoryAPI openapi.CategoryAPI,
) *GinServer {
	server := &GinServer{
		serverCfg:    serverCfg,
		logHelper:    logHelper,
		authorAPI:    authorAPI,
		bookAPI:      bookAPI,
		publisherAPI: publisherAPI,
		categoryAPI:  categoryAPI,
	}
	server.InitializeEngine()
	server.RegisterRoutes()
	return server
}

type GinServer struct {
	engine       *gin.Engine
	serverCfg    *config.ServerConfig
	logHelper    *logging.LogHelper
	authorAPI    openapi.AuthorAPI
	bookAPI      openapi.BookAPI
	publisherAPI openapi.PublisherAPI
	categoryAPI  openapi.CategoryAPI
}

func (s *GinServer) Run() {
	addr := net.JoinHostPort(s.serverCfg.Host, s.serverCfg.Port)
	err := s.engine.Run(addr)
	if err != nil {
		panic(err)
	}
}

func (s *GinServer) InitializeEngine() {
	gin.SetMode(s.serverCfg.Mode)
	gin.DefaultWriter = s.logHelper.Logger.WriterLevel(logrus.InfoLevel)
	gin.DefaultErrorWriter = s.logHelper.Logger.WriterLevel(logrus.ErrorLevel)
	s.engine = gin.Default()
	s.engine.Use(otelgin.Middleware(string(constant.TracerNameGinServer)))
}

func (s *GinServer) RegisterRoutes() {
	api := s.engine.Group("")
	{
		api.GET("/api/v1/authors", s.authorAPI.ApiV1AuthorsGet)
		api.POST("/api/v1/authors", s.authorAPI.ApiV1AuthorsPost)
		api.GET("/api/v1/authors/:id", s.authorAPI.ApiV1AuthorsIdGet)
		api.PUT("/api/v1/author/:id", s.authorAPI.ApiV1AuthorsIdPut)
		api.DELETE("/api/v1/author/:id", s.authorAPI.ApiV1AuthorsIdDelete)

		api.GET("/api/v1/books", s.bookAPI.ApiV1BooksGet)
		api.POST("/api/v1/books", s.bookAPI.ApiV1BooksPost)
		api.GET("/api/v1/books/:id", s.bookAPI.ApiV1BooksIdGet)
		api.PUT("/api/v1/books/:id", s.bookAPI.ApiV1BooksIdPut)
		api.DELETE("/api/v1/books/:id", s.bookAPI.ApiV1BooksIdDelete)

		api.GET("/api/v1/publishers", s.publisherAPI.ApiV1PublishersGet)
		api.POST("/api/v1/publishers", s.publisherAPI.ApiV1PublishersPost)
		api.GET("/api/v1/publishers/:id", s.publisherAPI.ApiV1PublishersIdGet)
		api.PUT("/api/v1/publisher/:id", s.publisherAPI.ApiV1PublishersIdPut)
		api.DELETE("/api/v1/publisher/:id", s.publisherAPI.ApiV1PublishersIdDelete)

		api.GET("/api/v1/categories", s.categoryAPI.ApiV1CategoriesGet)
	}
}

package gin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/madevara24/go-common/errors"
	"github.com/madevara24/go-common/response"
	"github.com/madevara24/go-common/server"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type ginHttpServer struct {
	config server.Config
	router *gin.Engine
}

type Option func(*ginHttpServer)

func NewGinHttpServer(config server.Config, opts ...Option) (*ginHttpServer, error) {
	handler := ginHttpServer{
		router: gin.New(),
		config: config,
	}

	switch config.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	handler.setMiddleware()

	for _, opt := range opts {
		opt(&handler)
	}

	return &handler, nil
}

func (g *ginHttpServer) setMiddleware() {
	defaultOrigins := []string{"http://localhost:3000"}
	if len(g.config.AllowedOrigins) > 0 {
		urls := strings.Split(g.config.AllowedOrigins, ",")
		defaultOrigins = append(defaultOrigins, urls...)
	}

	g.router.Use(cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowedOrigins:   defaultOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "content-type", "Origin", "Accept", "Authorization"},
	}))

	g.router.Use(requestid.New())
	g.router.Use(gzip.Gzip(gzip.DefaultCompression))

	g.router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		c.JSON(http.StatusInternalServerError, response.BasePayload{
			Error:   errors.ErrorCodeGeneralError,
			Success: false,
			Message: fmt.Sprintf("panic error: %s", err),
		})
	}))
}

func (g *ginHttpServer) GetRouter() *gin.Engine {
	return g.router
}

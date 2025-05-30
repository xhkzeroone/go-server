package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type RouteConfig struct {
	Path       string
	Method     string
	Handler    gin.HandlerFunc
	Middleware []gin.HandlerFunc
}

// Server implements Server for Gin.
type Server struct {
	*ProxyRouter
	*ProviderRouter
	engine     *gin.Engine
	rootGroup  *gin.RouterGroup
	config     *Config
	httpServer *http.Server
}

func NewServer(configs ...*Config) *Server {
	cfg := GetConfig(configs...)
	gin.SetMode(cfg.Mode)
	engine := gin.New()
	rootGroup := engine.Group(cfg.RootPath)
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &Server{
		ProxyRouter:    &ProxyRouter{},
		ProviderRouter: &ProviderRouter{},
		engine:         engine,
		rootGroup:      rootGroup,
		config:         cfg,
	}
}

func (s *Server) Start() error {
	addr := s.config.GetAddr()
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	//add api from @Api tag
	s.LoadRouter()
	s.Routes(s.ProxyRouter.Routes)

	//add api from provider route
	s.Routes(s.ProviderRouter.Routes)

	log.Printf("Server is running at %s", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Routes(routes []RouteConfig) {
	for _, r := range routes {
		s.Add(r.Method, r.Path, r.Handler, r.Middleware...)
	}
}

func (s *Server) Add(method, path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) {
	fullHandlers := append(middleware, handler)
	switch method {
	case http.MethodGet:
		s.rootGroup.GET(path, fullHandlers...)
	case http.MethodPost:
		s.rootGroup.POST(path, fullHandlers...)
	case http.MethodPut:
		s.rootGroup.PUT(path, fullHandlers...)
	case http.MethodDelete:
		s.rootGroup.DELETE(path, fullHandlers...)
	case http.MethodPatch:
		s.rootGroup.PATCH(path, fullHandlers...)
	default:
		log.Printf("Unsupported method: %s for path: %s", method, path)
	}
}

func (s *Server) HealthCheck() {
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	s.engine.GET("/liveness", func(c *gin.Context) {
		c.JSON(StatusOK, gin.H{"status": "alive"})
	})

	s.engine.GET("/readiness", func(c *gin.Context) {
		// Bạn có thể kiểm tra kết nối DB, Redis, etc. tại đây
		c.JSON(StatusOK, gin.H{"status": "ready"})
	})

	s.engine.POST("/terminate", func(c *gin.Context) {
		go func() {
			time.Sleep(1 * time.Second)
			_ = s.Shutdown(context.Background())
		}()
		c.JSON(StatusOK, gin.H{"status": "terminating"})
	})
}

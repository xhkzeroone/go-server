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
	*gin.Engine
	*ProxyRouter
	*ProviderRouter
	rootGroup  *gin.RouterGroup
	config     *Config
	httpServer *http.Server
}

func New(configs ...*Config) *Server {
	cfg := GetConfig(configs...)
	gin.SetMode(cfg.Mode)
	engine := gin.Default()
	rootGroup := engine.Group(cfg.RootPath)

	return &Server{
		Engine:         engine,
		ProxyRouter:    &ProxyRouter{},
		ProviderRouter: &ProviderRouter{},
		rootGroup:      rootGroup,
		config:         cfg,
	}
}

func (s *Server) Start() error {
	addr := s.config.GetAddr()
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.Engine,
	}

	//add api from @Api tag
	s.LoadRouter()
	s.Routes(s.ProxyRouter.Routes)

	//add api from provider route
	s.Routes(s.ProviderRouter.Routes)

	log.Printf("Server is running at %s", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
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

func (s *Server) Use(middleware ...gin.HandlerFunc) {
	s.rootGroup.Use(middleware...)
}

func (s *Server) Group(relativePath string, middleware ...gin.HandlerFunc) *gin.RouterGroup {
	return s.rootGroup.Group(relativePath, middleware...)
}

func (s *Server) Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.rootGroup.Handle(httpMethod, relativePath, handlers...)
}

func (s *Server) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodGet, relativePath, handlers...)
}

func (s *Server) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodPost, relativePath, handlers...)
}

func (s *Server) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodPut, relativePath, handlers...)
}

func (s *Server) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodDelete, relativePath, handlers...)
}

func (s *Server) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodPatch, relativePath, handlers...)
}

func (s *Server) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodOptions, relativePath, handlers...)
}

func (s *Server) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.Handle(http.MethodOptions, relativePath, handlers...)
}

func (s *Server) Add(method, path string, handler gin.HandlerFunc, middleware ...gin.HandlerFunc) {
	handlers := append(middleware, handler)
	s.Handle(method, path, handlers...)
}

func (s *Server) HealthCheck() {
	s.Engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	s.Engine.GET("/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	s.Engine.GET("/readiness", func(c *gin.Context) {
		// Bạn có thể kiểm tra kết nối DB, Redis, etc. tại đây
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	s.Engine.POST("/terminate", func(c *gin.Context) {
		go func() {
			time.Sleep(1 * time.Second)
			_ = s.Stop(context.Background())
		}()
		c.JSON(http.StatusOK, gin.H{"status": "terminating"})
	})
}

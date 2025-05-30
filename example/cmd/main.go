package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xhkzeroone/go-server/example/internal/api"
	"github.com/xhkzeroone/go-server/server"
	"log"
	"net/http"
)

func main() {
	// Tạo cấu hình cho serv
	cfg := &server.Config{
		Host:     "localhost",
		Port:     "8081",
		Mode:     "debug",
		RootPath: "/api/v1/",
	}
	// Khởi tạo serv với cấu hình
	serv := server.New(cfg)

	funcHandler := []server.RouteConfig{
		{
			Path:   "/users/index",
			Method: http.MethodGet,
			Handler: func(c *gin.Context) {
				c.JSON(200, map[string]string{
					"message": "/users/index",
				})
			},
			Middleware: []gin.HandlerFunc{func(c *gin.Context) {
				log.Println("Test /users/index")
				c.Next()
			}},
		},
	}
	serv.Routes(funcHandler)

	serv.RegisterHandlersWithTags(&api.MyApiHandler{})

	serv.Group("/admin").GET("/index", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "/admin/index",
		})
	})

	serv.HealthCheck()

	// Bắt đầu chạy serv
	if err := serv.Start(); err != nil {
		log.Fatalf("Error starting serv: %v", err)
	}
}

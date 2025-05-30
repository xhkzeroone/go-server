package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xhkzeroone/go-server/server"
)

// MyApiHandler
// @BaseUrl /iloveu/
type MyApiHandler struct {
}

// SayHi API
// @Api GET /say-hi
func (h *MyApiHandler) SayHi(c *gin.Context) {
	c.JSON(200, map[string]string{
		"message": "SayHi Api",
	})
}

// Print555 API
// @Api GET /print/555
func (h *MyApiHandler) Print555(c *gin.Context) {
	c.JSON(200, map[string]string{
		"message": "Print555 Api",
	})
}

func (h *MyApiHandler) SayHello(c *gin.Context) {
	c.JSON(200, map[string]string{
		"message": "SayHello",
	})
}

func (h *MyApiHandler) Routes() []server.RouteConfig {
	return []server.RouteConfig{
		{
			Method:  "GET",
			Path:    "/say-hello",
			Handler: h.SayHello,
		},
	}
}

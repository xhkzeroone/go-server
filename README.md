# go-server

## Giới thiệu

`go-server` là một module HTTP server sử dụng framework Gin, hỗ trợ cấu hình động, đăng ký route tự động qua tag, và dễ dàng mở rộng cho các dịch vụ backend hiện đại.  
Module này phù hợp để xây dựng các microservice hoặc backend API với khả năng mở rộng, bảo trì và tích hợp dễ dàng.

---

## Tính năng nổi bật

- **Cấu hình linh hoạt**: Dễ dàng cấu hình qua biến môi trường hoặc file cấu hình (hỗ trợ bởi [viper](https://github.com/spf13/viper)).
- **Đăng ký route động**: Tự động quét và đăng ký các route dựa trên tag `@Api` trong code.
- **Health check endpoints**: Tích hợp sẵn các API kiểm tra tình trạng server (`/ping`, `/liveness`, `/readiness`, `/terminate`).
- **Middleware & Group route**: Hỗ trợ thêm middleware và group route linh hoạt.
- **ProviderRouter & ProxyRouter**: Dễ dàng mở rộng, tích hợp các router động hoặc provider bên ngoài.

---

## Cấu trúc thư mục

```
go-server/
├── server/
│   ├── config.go            # Định nghĩa và load cấu hình server
│   ├── dynamic-router.go    # Đăng ký route động qua tag @Api
│   ├── provider-router.go   # Đăng ký route qua interface Route
│   └── server.go            # Khởi tạo, chạy, stop server và các API health check
├── go.mod
├── go.sum
└── README.md
```

---

## Hướng dẫn sử dụng

### 1. Cài đặt phụ thuộc

```bash
go mod tidy
```

### 2. Cấu hình

Cấu hình server qua biến môi trường hoặc file cấu hình (hỗ trợ bởi viper):

- `SERVER_HOST` (default: `localhost`)
- `SERVER_PORT` (default: `8080`)
- `SERVER_MODE` (default: `debug`)
- `SERVER_ROOTPATH` (default: `""`)

Ví dụ, để chạy ở cổng 9000:

```bash
export SERVER_PORT=9000
```

### 3. Khởi động server cơ bản

Tạo file `main.go`:

```go
package main

import (
    "go-server/server"
)

func main() {
    cfg := server.DefaultConfig()
    srv := server.New(cfg)
    srv.HealthCheck() // Đăng ký các API health check
    if err := srv.Start(); err != nil {
        panic(err)
    }
}
```

Chạy server:

```bash
go run main.go
```

---

## 4. Các API mặc định

- `GET /ping` → `{ "message": "pong" }`  
  Kiểm tra server còn sống.
- `GET /liveness` → `{ "status": "alive" }`  
  Kiểm tra server có đang hoạt động.
- `GET /readiness` → `{ "status": "ready" }`  
  Kiểm tra server đã sẵn sàng nhận request (có thể kiểm tra DB, cache, ...).
- `POST /terminate` → `{ "status": "terminating" }`  
  Yêu cầu server tắt an toàn.

---

## 5. Đăng ký route động với ProxyRouter

Bạn có thể đăng ký các handler với tag `@Api` để tự động tạo route.

**Ví dụ:**

```go
// user_handler.go

package handler

import (
    "github.com/gin-gonic/gin"
)

// @Api UserGet GET /user/:id
func (h *UserHandler) UserGet(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
}
```

Trong `main.go`:

```go
package main

import (
    "go-server/server"
    "your-module/handler"
)

func main() {
    cfg := server.DefaultConfig()
    srv := server.New(cfg)

    // Đăng ký handler động
    userHandler := &handler.UserHandler{}
    srv.ProxyRouter.RegisterHandlersWithTags(userHandler)
    srv.LoadRouter() // Quét và đăng ký các route động

    srv.HealthCheck()
    if err := srv.Start(); err != nil {
        panic(err)
    }
}
```

---

## 6. Đăng ký route qua ProviderRouter

Nếu bạn muốn đăng ký route qua interface `Route`:

```go
type MyProvider struct{}

func (p *MyProvider) Routes() []server.RouteConfig {
    return []server.RouteConfig{
        {
            Method:  "GET",
            Path:    "/hello",
            Handler: func(c *gin.Context) { c.String(200, "Hello from provider!") },
        },
    }
}
```

Trong `main.go`:

```go
provider := &MyProvider{}
srv.ProviderRouter.RegisterHandlers(provider)
srv.Routes(srv.ProviderRouter.Routes)
```

---

## 7. Thêm middleware & group route

**Thêm middleware:**

```go
srv.Use(gin.Logger(), gin.Recovery())
```

**Tạo group route:**

```go
apiGroup := srv.Group("/api")
apiGroup.GET("/status", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})
```

---

## 8. Mở rộng & tích hợp

- Có thể tích hợp thêm các middleware xác thực, logging, CORS, ...
- Dễ dàng mở rộng thêm các router động hoặc provider mới.
- Có thể tích hợp với các hệ thống khác như database, cache, message queue, ...

---

## 9. Đóng góp

Mọi đóng góp, issue hoặc pull request đều được hoan nghênh!  
Vui lòng tạo issue hoặc PR trên repository này.

---

## License

MIT

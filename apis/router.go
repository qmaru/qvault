package apis

import (
	"log"
	"os"

	"qkms/apis/health"
	"qkms/apis/secret"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Run() error {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")
	api.GET("/health", health.Health)

	secretGroup := api.Group("/secrets", AuthMiddleware)
	secretGroup.GET("", secret.ListKeys)
	secretGroup.GET("/:key", secret.GetKey)
	secretGroup.PUT("/:key", secret.PutKey)
	secretGroup.DELETE("/:key", secret.DeleteKey)

	listenAddr := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
	log.Println("Listening on", listenAddr)

	return e.Start(listenAddr)
}

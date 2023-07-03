package main

import (
	"os"

	"github.com/IDEA/SERVER/conf"
	"github.com/IDEA/SERVER/pkg/gateway"
	"github.com/IDEA/SERVER/pkg/handler"
	"github.com/IDEA/SERVER/pkg/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	conf.NewEnv()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FRONTNED_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))
	dg := gateway.NewDiscordGateway()
	ns := service.NewNotifyService(dg)
	h := handler.NewHandler(ns)
	e.POST("/application", h.HandleApplication)
	e.Logger.Fatal(e.Start(":8080"))
}

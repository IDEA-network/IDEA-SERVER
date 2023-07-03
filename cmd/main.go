package main

import (
	"github.com/IDEA/SERVER/conf"
	"github.com/IDEA/SERVER/pkg/gateway"
	"github.com/IDEA/SERVER/pkg/handler"
	"github.com/IDEA/SERVER/pkg/service"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	conf.NewEnv()
	dg := gateway.NewDiscordGateway()
	ns := service.NewNotifyService(dg)
	h := handler.NewHandler(ns)
	e.GET("/application", h.HandleApplication)
	e.Logger.Fatal(e.Start(":8080"))
}

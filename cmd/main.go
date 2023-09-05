package main

import (
	"os"

	"github.com/IDEA/SERVER/conf"
	"github.com/IDEA/SERVER/pkg/gateway"
	"github.com/IDEA/SERVER/pkg/handler"
	"github.com/IDEA/SERVER/pkg/repository"
	"github.com/IDEA/SERVER/pkg/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	conf.NewEnv()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FRONTNED_URL"), "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))
	dg := gateway.NewDiscordGateway()
	rc := gateway.NewRedisCilent()
	cr := repository.NewCacheRepo(rc)
	gg := gateway.NewGoogleOAuthGateway(cr)
	ns := service.NewNotifyService(dg)
	es := service.NewEmailService(dg, gg)
	ms := service.NewManageMemberService(gg)
	h := handler.NewHandler(ns, es, ms)
	e.POST("/application", h.HandleApplication)
	e.Logger.Fatal(e.Start(":8080"))
}

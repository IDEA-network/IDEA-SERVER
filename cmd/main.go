package main

import (
	"github.com/IDEA/SERVER/pkg/usecase"
	"net/http"
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
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))
	dg := gateway.NewDiscordGateway()
	rc := gateway.NewRedisCilent()
	cr := repository.NewCacheRepo(rc)
	gg := gateway.NewGoogleOAuthGateway(cr)
	ns := service.NewNotifyService(dg)
	es := service.NewEmailService(dg, gg)
	ms := service.NewManageMemberService(gg)
	entry := usecase.NewEntryUsecase(dg, ns)
	h := handler.NewHandler(ns, es, ms, entry)
	e.POST("/application", h.HandleApplication)
	e.POST("/contact", h.HandleContact)
	e.POST("/entry", h.HandleEntry)
	e.Logger.Fatal(e.Start(":8080"))
}

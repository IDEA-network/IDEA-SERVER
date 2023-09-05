package handler

import (
	"github.com/IDEA/SERVER/pkg/service"
)

type Handler struct {
	ns service.NotifyService
	es *service.EmailService
	ms *service.ManageMemberService
}

func NewHandler(ns service.NotifyService, es *service.EmailService, ms *service.ManageMemberService) *Handler {
	return &Handler{ns: ns, es: es, ms: ms}
}

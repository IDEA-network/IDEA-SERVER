package handler

import (
	"github.com/IDEA/SERVER/pkg/service"
	"github.com/IDEA/SERVER/pkg/usecase"
)

type Handler struct {
	ns    service.NotifyService
	es    *service.EmailService
	ms    *service.ManageMemberService
	entry *usecase.EntryUsecase
}

func NewHandler(ns service.NotifyService, es *service.EmailService, ms *service.ManageMemberService, entry *usecase.EntryUsecase) *Handler {
	return &Handler{ns: ns, es: es, ms: ms, entry: entry}
}

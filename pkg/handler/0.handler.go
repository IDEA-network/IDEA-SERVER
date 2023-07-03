package handler

import "github.com/IDEA/SERVER/pkg/service"

type Handler struct {
	ns service.NotifyService
}

func NewHandler(ns service.NotifyService) *Handler {
	return &Handler{ns: ns}
}

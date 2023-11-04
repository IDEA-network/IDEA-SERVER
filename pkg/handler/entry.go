package handler

import (
	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/labstack/echo"
)

type EntryResponse struct {
	InviteURL   *dto.InviteURL   `json:"invite_url"`
	Application *dto.Application `json:"application"`
}

func (h *Handler) HandleEntry(c echo.Context) error {
	var application dto.Application
	if err := c.Bind(&application); err != nil {
		return c.JSON(403, err.Error())
	}
	res, err := h.entry.AcceptEntryApplication(&application)
	if err != nil {
		return c.JSON(400, err)
	}
	return c.JSON(200, &EntryResponse{
		Application: res.Application, InviteURL: res.InviteURL,
	})
}

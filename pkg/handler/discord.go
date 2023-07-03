package handler

import (
	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/labstack/echo"
)

func (h *Handler) HandleApplication(c echo.Context) error {
	var application dto.Application
	if err := c.Bind(&application); err != nil {
		return c.JSON(403, err.Error())
	}
	if err := h.ns.NotifyApplication(&application); err != nil {
		return c.JSON(500, err.Error())
	}
	return c.JSON(200, "application success")
}

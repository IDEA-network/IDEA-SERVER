package handler

import (
	"context"
	"github.com/IDEA/SERVER/pkg/service"
	"time"

	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/util"
	"github.com/labstack/echo"
)

func (h *Handler) HandleApplication(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var application dto.Application
	var limiter = util.NewRateLimiter(5, 10)
	if err := limiter.Limit(ctx, func() error {
		if err := c.Bind(&application); err != nil {
			return c.JSON(403, err.Error())
		}
		req := &service.NotifyApplicationRequest{
			Application:     &application,
			InviteDiscordID: "",
		}
		if err := h.ns.NotifyApplicationToDiscord(req); err != nil {
			return c.JSON(500, err.Error())
		}
		return nil
	}); err != nil {
		return c.JSON(401, err.Error())
	}
	go func() {
		if err := h.es.SendInviteEmail(&application); err != nil {
			h.ns.NotifyEmailErrorToDiscord(&application, err.Error())
		}
	}()
	go func() {
		if err := h.ms.StoreApplicationMemberData(&application); err != nil {
			h.ns.NotifyStoreApplicationErrorToDiscord(&application, err.Error())
		}
	}()
	return c.JSON(200, "application success")
}

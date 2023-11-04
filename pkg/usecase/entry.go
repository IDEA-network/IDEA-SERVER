package usecase

import (
	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/gateway"
	"github.com/IDEA/SERVER/pkg/service"
)

type EntryUsecase struct {
	discord gateway.DiscordGateway
	notify  service.NotifyService
}

func NewEntryUsecase(discord gateway.DiscordGateway, notify service.NotifyService) *EntryUsecase {
	return &EntryUsecase{discord: discord, notify: notify}
}

type AcceptEntryApplicationResponse struct {
	InviteURL   *dto.InviteURL
	Application *dto.Application
}

// 申請をDiscordに通知 & 招待URLを生成する
func (u *EntryUsecase) AcceptEntryApplication(application *dto.Application) (*AcceptEntryApplicationResponse, error) {
	inviteURL, err := u.discord.CreateIviteURL()
	if err != nil {
		return nil, err
	}
	req := &service.NotifyApplicationRequest{
		Application:     application,
		InviteDiscordID: inviteURL.ID(),
	}
	if err := u.notify.NotifyApplicationToDiscord(req); err != nil {
		return nil, nil
	}
	return &AcceptEntryApplicationResponse{
		InviteURL:   &inviteURL,
		Application: application,
	}, nil
}

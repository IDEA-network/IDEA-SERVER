package service

import (
	"fmt"
	"os"

	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/gateway"
)

type NotifyService interface {
	NotifyApplication(application *dto.Application) error
}

type notifyService struct {
	dg gateway.DiscordGateway
}

func NewNotifyService(dg gateway.DiscordGateway) NotifyService {
	return &notifyService{dg: dg}
}

func (s *notifyService) NotifyApplication(application *dto.Application) error {
	var message string
	message = fmt.Sprintf("氏名: %s\n連絡先: %s\n", application.Name, application.Email)
	for i, v := range application.Surveys {
		index := i + 1
		message += fmt.Sprintf("質問[%d]: %s\n回答[%d]: %s\n\n", index, v.Question, index, v.Answer)
	}
	payload := dto.DiscordPayload{
		Username:  "入会申請のお知らせ📢",
		AvatarUrl: "https://img.benesse-cms.jp/pet-cat/item/image/normal/f3978ebc-9030-49e7-aa5e-4a370a955e1b.jpg?w=1200&h=1200&resize_type=cover&resize_mode=force",
		Content:   message,
	}
	webhookURL := os.Getenv("APPLICATION_WEBHOOK")
	if err := s.dg.SendMessage(webhookURL, payload); err != nil {
		return err
	}
	return nil
}

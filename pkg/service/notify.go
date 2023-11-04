package service

import (
	"fmt"
	"log"
	"os"

	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/gateway"
)

type NotifyService interface {
	NotifyApplicationToDiscord(req *NotifyApplicationRequest) error
	NotifyContactToDiscord(application *dto.Application) error
	NotifyEmailErrorToDiscord(application *dto.Application, errorMessage string)
	NotifyStoreApplicationErrorToDiscord(application *dto.Application, errorMessage string)
}

type notifyService struct {
	dg gateway.DiscordGateway
}

func NewNotifyService(dg gateway.DiscordGateway) NotifyService {
	return &notifyService{dg: dg}
}

func (s *notifyService) NotifyContactToDiscord(application *dto.Application) error {
	var message string
	message = fmt.Sprintf("氏名: %s\n連絡先: %s\n", application.Name, application.Email)
	for i, v := range application.Surveys {
		index := i + 1
		message += fmt.Sprintf("質問[%d]: %s\n回答[%d]: %s\n\n", index, v.Question, index, v.Answer)
	}
	payload := dto.DiscordPayload{
		Username:  "お問合せのお知らせ📢",
		AvatarUrl: "https://img.benesse-cms.jp/pet-cat/item/image/normal/f3978ebc-9030-49e7-aa5e-4a370a955e1b.jpg?w=1200&h=1200&resize_type=cover&resize_mode=force",
		Content:   message,
	}
	webhookURL := "https://discord.com/api/webhooks/1159104306067280003/g53zBA-iqtTVuP0R7qXQ5NcOmYYaWVYDdjH10n_cbkoRNlMgYJmVHhSUegarGvxTEZjc"
	if err := s.dg.SendMessage(webhookURL, payload); err != nil {
		return err
	}
	return nil
}

type NotifyApplicationRequest struct {
	Application     *dto.Application
	InviteDiscordID string
}

func (s *notifyService) NotifyApplicationToDiscord(req *NotifyApplicationRequest) error {
	var message string
	message += fmt.Sprintf("[招待ID]: %s\n", req.InviteDiscordID)
	message += fmt.Sprintf("[氏名]: %s\n[連絡先]: %s\n\n", req.Application.Name, req.Application.Email)
	for _, v := range req.Application.Surveys {
		message += fmt.Sprintf("[%s]: %s\n\n", v.Question, v.Answer)
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

func (s *notifyService) NotifyEmailErrorToDiscord(application *dto.Application, errorMessage string) {
	msg := fmt.Sprintf("%s(%s)さんへの招待メールの送信に失敗しました\n (Details: %s)", application.Name, application.Email, errorMessage)
	s.sendErrorNotificationToDiscord(application, msg)
}

func (s *notifyService) NotifyStoreApplicationErrorToDiscord(application *dto.Application, errorMessage string) {
	msg := fmt.Sprintf("%s(%s)さんの申請情報の保存に失敗しました\n (Details: %s)", application.Name, application.Email, errorMessage)
	s.sendErrorNotificationToDiscord(application, msg)
}

func (s *notifyService) sendErrorNotificationToDiscord(application *dto.Application, message string) {
	payload := dto.DiscordPayload{
		Username:  "システムエラーのお知らせ",
		AvatarUrl: "https://img.benesse-cms.jp/pet-cat/item/image/normal/f3978ebc-9030-49e7-aa5e-4a370a955e1b.jpg?w=1200&h=1200&resize_type=cover&resize_mode=force",
		Content:   message,
	}
	webhookURL := os.Getenv("APPLICATION_WEBHOOK")
	if err := s.dg.SendMessage(webhookURL, payload); err != nil {
		log.Fatal(err.Error())
	}
}

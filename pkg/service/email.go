package service

import (
	"bytes"
	"text/template"

	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/gateway"
)

type EmailService struct {
	dg gateway.DiscordGateway
	gg gateway.GoogleOAuthGateway
}

func NewEmailService(dg gateway.DiscordGateway, gg gateway.GoogleOAuthGateway) *EmailService {
	return &EmailService{dg: dg, gg: gg}
}

func (s *EmailService) SendInviteEmail(application *dto.Application) error {
	inviteURL, err := s.dg.CreateIviteURL()
	content, err := generateContent(application.Name, inviteURL)
	if err != nil {
		return err
	}
	emailPayload := &dto.EmailPayload{
		ToAddress: application.Email,
		Subject:   `【IDEA】入会申請ありがとうございます !`,
		Content:   content,
	}
	if err := s.gg.SendEmailByGmail(emailPayload); err != nil {
		return err
	}
	return nil
}

func generateContent(name, inviteURL string) (string, error) {
	tmpl := `
{{.Name}}さん

お世話になっております。学生団体 IDEA です。

{{.Name}}さんとは是非一緒に活動していきたいと思いますので
以下のリンクよりDiscordサーバーに入っていただき、自己紹介チャンネルにて自己紹介をお願いいたします！

招待リンク → {{.InviteURL}}

招待リンクは7日後に期限切れとなりますのでお早めにお願いいたします。

それでは、今後ともよろしくお願いいたします。

学生団体 IDEA
	`
	type MailData struct {
		Name      string
		InviteURL string
	}
	data := MailData{
		Name:      name,
		InviteURL: inviteURL,
	}

	t, err := template.New("mail").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

package gateway_test

import (
	"os"
	"testing"

	"github.com/IDEA/SERVER/conf"
	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/gateway"
)

func Test_Discord(t *testing.T) {
	conf.NewEnv()
	dg := gateway.NewDiscordGateway()
	payload := dto.DiscordPayload{
		Username:  "入会申請のお知らせ📢",
		AvatarUrl: "https://img.benesse-cms.jp/pet-cat/item/image/normal/f3978ebc-9030-49e7-aa5e-4a370a955e1b.jpg?w=1200&h=1200&resize_type=cover&resize_mode=force",
		Content:   "テストメッセージ",
	}
	webhookURL := os.Getenv("APPLICATION_WEBHOOK")
	err := dg.SendMessage(webhookURL, payload)
	if err != nil {
		t.Error(err.Error())
	}
}

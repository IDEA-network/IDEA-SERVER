package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IDEA/SERVER/pkg/dto"
)

type DiscordGateway interface {
	SendMessage(webhookURL string, payload dto.DiscordPayload) error
}

type discordGateway struct{}

func NewDiscordGateway() DiscordGateway {
	return &discordGateway{}
}

func (dg *discordGateway) SendMessage(webhookURL string, payload dto.DiscordPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling json", err)
		return err
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending message", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

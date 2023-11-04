package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/bwmarrin/discordgo"
)

const (
	INVITE_CHANNEL_ID = "1048672304177098848"
)

type DiscordGateway interface {
	SendMessage(webhookURL string, payload dto.DiscordPayload) error
	CreateIviteURL() (dto.InviteURL, error)
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

func (dg *discordGateway) CreateIviteURL() (dto.InviteURL, error) {
	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	session, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		return "", err
	}
	invite, err := session.ChannelInviteCreate(INVITE_CHANNEL_ID, discordgo.Invite{
		MaxAge:    60 * 60 * 24 * 7,
		MaxUses:   1,
		Temporary: false,
		Unique:    true,
	})
	if err != nil {
		return "", err
	}
	inviteURL := fmt.Sprintf("https://discord.gg/%s", invite.Code)
	return dto.InviteURL(inviteURL), nil
}

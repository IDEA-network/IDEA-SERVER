package dto

import "strings"

type DiscordPayload struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
}

type InviteURL string

func (i InviteURL) ID() string {
	parts := strings.Split(i.String(), "/")
	lastPart := parts[len(parts)-1]

	return lastPart
}

func (i InviteURL) String() string {
	return string(i)
}

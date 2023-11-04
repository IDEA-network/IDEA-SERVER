package dto

import "strings"

type DiscordPayload struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
}

type InviteURL string

func (i InviteURL) ID() string {
	lowercaseURL := strings.ToLower(string(i))
	// URLをスラッシュで分割し、最後の要素を取得
	parts := strings.Split(lowercaseURL, "/")
	lastPart := parts[len(parts)-1]

	return lastPart
}

func (i InviteURL) String() string {
	return string(i)
}

package dto

type DiscordPayload struct {
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Content   string `json:"content"`
}

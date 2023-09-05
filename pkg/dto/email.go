package dto

type EmailPayload struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

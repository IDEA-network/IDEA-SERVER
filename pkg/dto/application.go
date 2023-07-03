package dto

type Application struct {
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Surveys []Surveys `json:"surveys"`
}
type Surveys struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

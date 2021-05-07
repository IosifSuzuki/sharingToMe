package models

type Notification struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Error       Error  `json:"error"`
}

package models

type Subscription struct {
	ChatID int64    `json:"ChatId"`
	Groups []string `json:"Groups"`
}

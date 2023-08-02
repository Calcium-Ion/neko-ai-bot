package model

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	ChatId   string `json:"chat_id"`
	Balance  int    `json:"balance"`
}

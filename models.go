package main

type Config struct {
	BotToken, BotUsername string
}

type User struct {
	ID              int `gorm:"primary_key"`
	Username, Token string
	RefCount        int
}
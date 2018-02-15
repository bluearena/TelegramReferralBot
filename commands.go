package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"strconv"
)

func start(message *tgbotapi.Message){
	fields := strings.Fields(message.Text)
	if len(fields) == 1{
		var user User
		db.First(&user, "id = ?", message.From.ID)
		if user == (User{}){
			user := User{message.From.ID, message.From.FirstName, generateToken(), 0}
			db.Create(&user)
			sendMessage(message.Chat.ID, phrases[0] + " " + phrases[2] +
				"\nt.me/" + configuration.BotUsername + "?start=" + user.Token + "\n" + phrases[5], nil)
		}else {
			sendMessage(message.Chat.ID, phrases[3], nil)
		}
	}else if len(fields) == 2{
		var user User
		db.Find(&user, "token = ?", fields[1])
		if user == (User{}){
			user = User{}
			db.First(&user, "id = ?", message.From.ID)
			if user != (User{}){
				sendMessage(message.Chat.ID, phrases[3], nil)
				return
			}else {
				user := User{message.From.ID, message.From.FirstName, generateToken(), 0}
				db.Create(&user)
				sendMessage(message.Chat.ID, phrases[0] + " " + phrases[2] +
					"t.me/" + configuration.BotUsername + "?start=" + user.Token + "\n" + phrases[5], nil)
			}
		}else {
			user2 := User{}
			db.First(&user2, "id = ?", message.From.ID)
			if user2 != (User{}){
				sendMessage(message.Chat.ID, phrases[3], nil)
				return
			}else {
				user.RefCount++
				db.Save(&user)

				user2 = User{message.From.ID, message.From.FirstName, generateToken(), 0}
				db.Create(&user2)
				sendMessage(message.Chat.ID, phrases[0] + phrases[1] + user.Username + " " + phrases[2] +
					"t.me/" + configuration.BotUsername + "?start=" + user2.Token + "\n" + phrases[5], nil)
			}
		}
	}
}

func refs(message *tgbotapi.Message){
	var user User
	db.First(&user, "id = ?", message.From.ID)
	sendMessage(message.Chat.ID, phrases[2] + "\nt.me/" + configuration.BotUsername + "?start=" + user.Token + "\n\n" +
		phrases[4] + strconv.Itoa(user.RefCount), nil)
}
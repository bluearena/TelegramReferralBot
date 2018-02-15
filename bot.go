package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"io"
)

var (
	bot           *tgbotapi.BotAPI
	configuration Config
	phrases       map[int]string
	db            *gorm.DB
)

func main() {
	initLog()
	initConfig()
	initStrings()
	initDB()

	var err error
	bot, err = tgbotapi.NewBotAPI(configuration.BotToken)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			log.Print("It is command")
			command := update.Message.Command()
			switch command {
			case "start":
				start(update.Message)
			case "refs":
				refs(update.Message)
			}
		}
	}
}

func initLog() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func initConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	log.Print("First 10 bytes from config.json")
	log.Print(body[:10])
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	log.Print("First 10 bytes after trim")
	reader := bytes.NewReader(body)
	log.Print(body[:10])
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&configuration)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}

}

func initStrings() {
	file, err := os.Open("strings.json")
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	log.Print("First 10 bytes from strings.json")
	log.Print(body[:10])
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	log.Print("First 10 bytes after trim")
	reader := bytes.NewReader(body)
	log.Print(body[:10])
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&phrases)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
}

func initDB() {
	var err error
	db, err = gorm.Open("sqlite3", "data.db")
	if err != nil {
		log.Print("********** ERROR: ")
		log.Panic("Failed to connect database")
	} else {
		log.Print("Opened DB")
	}
	db.LogMode(true)
	log.Print("Set LogMode")
	db.AutoMigrate(&User{})
	log.Print("Migrated")
}

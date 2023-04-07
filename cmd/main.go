package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sadreddinov/tgbot/pkg/telegram"
	"github.com/sirupsen/logrus"
)
func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())	
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		logrus.Fatalf("error connecting bot: %s", err.Error())
	}
	bot.Debug = true
	telegrambot := telegram.NewBot(bot)
	if err := telegrambot.Start();err != nil {
		logrus.Fatalf("error starting bot: %s", err.Error())
	} 
	
}
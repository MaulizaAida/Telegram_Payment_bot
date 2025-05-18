package main

import (
	"telegram-history-bot/bot"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	bot.StartBot()
}

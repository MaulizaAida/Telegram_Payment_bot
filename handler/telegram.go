package handler

import (
	"log"
	"os"
	"strings"

	"telegram-history-bot/ai"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func StartTelegramBot() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Gagal membuat bot:", err)
	}

	bot.Debug = false
	log.Printf("Bot %s telah aktif", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			userInput := update.Message.Text
			log.Printf("Pesan dari %s: %s", update.Message.From.UserName, userInput)

			prompt := ai.BuildPrompt(userInput)
			response, err := ai.CallOllama(prompt)
			if err != nil {
				response = "❌ AI tidak merespons. Pastikan Ollama aktif."
				log.Println("AI Error:", err)
			}

			// Remove any markdown formatting
			escapedResponse := "```json\n" + escapeMarkdownV2(response) + "\n```"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, escapedResponse)
			msg.ParseMode = "MarkdownV2"

			bot.Send(msg)
		}
	}
}

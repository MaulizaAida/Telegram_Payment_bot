package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"telegram-history-bot/ai"     // Updated import path
	"telegram-history-bot/sheets" // Updated import path

	"context"
	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	sheetsapi "google.golang.org/api/sheets/v4"
)

// Add Transaction type if not defined in sheets package
type Transaction struct {
	Amount        int    `json:"amount"`
	Description   string `json:"description"`
	PaymentMethod string `json:"payment_method"`
	Category      string `json:"category"`
	Type          string `json:"type"`
}

func StartBot() {
	idStr := os.Getenv("AUTHORIZED_USER_ID")
	authorizedUserID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatal("❌ AUTHORIZED_USER_ID tidak valid:", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("❌ BOT_TOKEN belum diset di .env")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// ✅ Validasi hanya user tertentu
		if update.Message.From.ID != authorizedUserID {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Maaf, kamu tidak diizinkan menggunakan bot ini.")
			bot.Send(msg)
			continue
		}

		handleMessage(bot, update)
	}
}

// Add this helper function
func initSheetsService() (*sheetsapi.Service, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(b, sheetsapi.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	client := conf.Client(ctx)
	srv, err := sheetsapi.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %v", err)
	}

	return srv, nil
}

// Update handleMessage to use the helper function
func formatAmount(amount int) string {
	str := strconv.Itoa(amount)
	var result []rune
	for i := len(str) - 1; i >= 0; i-- {
		if i != len(str)-1 && (len(str)-1-i)%3 == 0 {
			result = append([]rune{'.'}, result...)
		}
		result = append([]rune{rune(str[i])}, result...)
	}
	return string(result)
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Initialize sheets service once at the beginning
	srv, err := initSheetsService()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal membuat koneksi sheets: "+err.Error()))
		return
	}

	botUsername := "@" + bot.Self.UserName
	text := update.Message.Text

	// Add new command handlers
	switch {
	case text == "/summary" || text == "/summary"+botUsername:
		// Use existing srv instead of creating new one
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 Mengambil ringkasan...")
		statusMsg, _ := bot.Send(msg)

		summary, err := sheets.GetSummary(srv, os.Getenv("SHEET_ID"))
		if err != nil {
			errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal mengambil ringkasan: "+err.Error())
			bot.Send(errMsg)
			return
		}

		formattedSummary := sheets.FormatSummary(summary)

		// Update the "loading" message with the actual summary
		edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, statusMsg.MessageID, formattedSummary)
		edit.ParseMode = "MarkdownV2" // Enable markdown formatting
		bot.Send(edit)
		return
	
	case text == "/today" || text == "/today"+botUsername:
		handleHistory(bot, update, srv, "day")
		return
	
	case text == "/week" || text == "/week"+botUsername:
		handleHistory(bot, update, srv, "week")
		return
	
	case text == "/month" || text == "/month"+botUsername:
		handleHistory(bot, update, srv, "month")
		return
	
	case text == "/help" || text == "/help"+botUsername:
	    helpText := `🤖 *Daftar Perintah*

	/summary \- Ringkasan keuangan & 5 transaksi terakhir
	/today \- Riwayat transaksi hari ini
	/week \- Riwayat transaksi 7 hari terakhir
	/month \- Riwayat transaksi 30 hari terakhir
	/help \- Menampilkan bantuan ini

	Untuk mencatat transaksi, cukup ketik seperti:
	\- beli makan 25rb via gopay
	\- terima gaji 5jt via bca
	\- bayar listrik 200k cash`

	    msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
	    msg.ParseMode = "MarkdownV2"
	    bot.Send(msg)
	    return
	
	default:
		// Handle regular transaction input
		input := update.Message.Text
		prompt := ai.BuildPrompt(input)
		result, err := ai.CallOllama(prompt)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			return
		}

		username := update.Message.From.UserName
		if username == "" {
			username = fmt.Sprintf("id:%d", update.Message.From.ID)
		}

		// Use existing srv instead of creating new one
		if err := sheets.SaveToSheet(srv, result, username); err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal menyimpan: "+err.Error()))
			return
		}

		// Send confirmation with transaction details
		var tx Transaction
		if err := json.Unmarshal([]byte(result), &tx); err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal membaca hasil: "+err.Error()))
			return
		}

		// Use existing srv instead of creating new one
		summary, err := sheets.GetSummary(srv, os.Getenv("SHEET_ID"))
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal mengambil saldo: "+err.Error()))
			return
		}

		// In handleMessage function, update the message formatting
		emoji := "💸"
		displayType := "Pengeluaran"
		if tx.Type == "income" {
			emoji = "💰"
			displayType = "Pemasukan"
		}
    msg := fmt.Sprintf("%s Transaksi tercatat:\n%s %s\nRp %s via %s\n\n💵 Saldo saat ini: Rp %s",
        emoji, displayType, tx.Description,
        formatAmount(tx.Amount),
        tx.PaymentMethod,
        formatAmount(summary.Balance))
    bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
}
}

// Now handleHistory function can be properly defined
func handleHistory(bot *tgbotapi.BotAPI, update tgbotapi.Update, srv *sheetsapi.Service, period string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 Mengambil riwayat...")
	statusMsg, _ := bot.Send(msg)

	history, err := sheets.GetHistory(srv, os.Getenv("SHEET_ID"), period)
	if err != nil {
		errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal mengambil riwayat: "+err.Error())
		bot.Send(errMsg)
		return
	}

	formattedHistory := sheets.FormatHistory(history, period)

	edit := tgbotapi.NewEditMessageText(update.Message.Chat.ID, statusMsg.MessageID, formattedHistory)
	bot.Send(edit)
}
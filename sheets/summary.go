package sheets

import (
	"fmt"
	"strconv"
)

type Summary struct {
	TotalIncome        int
	TotalExpense       int
	Balance            int
	RecentTransactions []Transaction
}

func FormatSummary(summary *Summary) string {
	result := "📊 Ringkasan Keuangan:\n\n"
	result += fmt.Sprintf("💰 Total Pemasukan: Rp %s\n", formatAmount(summary.TotalIncome))
	result += fmt.Sprintf("💸 Total Pengeluaran: Rp %s\n", formatAmount(summary.TotalExpense))
	result += fmt.Sprintf("💵 Saldo: Rp %s\n\n", formatAmount(summary.Balance))

	result += "📝 5 Transaksi Terakhir:\n"
	for _, tx := range summary.RecentTransactions {
		emoji := "💸"
		txType := "Pengeluaran"
		if tx.Type == "income" {
			emoji = "💰"
			txType = "Pemasukan"
		}
		result += fmt.Sprintf("%s %s %s: Rp %s (%s)\n",
			emoji, txType, tx.Description, formatAmount(tx.Amount), tx.PaymentMethod)
	}

	return result
}

// Add new format function for history
func FormatHistory(history *TransactionHistory, period string) string {
	var timeRange string
	switch period {
	case "day":
		timeRange = "Hari Ini"
	case "week":
		timeRange = "7 Hari Terakhir"
	case "month":
		timeRange = "30 Hari Terakhir"
	}

	result := fmt.Sprintf("📊 Riwayat Transaksi %s:\n\n", timeRange)
	result += fmt.Sprintf("💰 Total Pemasukan: Rp %s\n", formatAmount(history.TotalIncome))
	result += fmt.Sprintf("💸 Total Pengeluaran: Rp %s\n", formatAmount(history.TotalExpense))
	result += fmt.Sprintf("💵 Saldo Period Ini: Rp %s\n\n", formatAmount(history.Balance))

	result += "📝 Daftar Transaksi:\n"
	for _, tx := range history.Transactions {
		emoji := "💸"
		txType := "Pengeluaran"
		if tx.Type == "income" {
			emoji = "💰"
			txType = "Pemasukan"
		}
		result += fmt.Sprintf("%s %s %s: Rp %s (%s)\n",
			emoji, txType, tx.Description, formatAmount(tx.Amount), tx.PaymentMethod)
	}

	return result
}

// Add helper function for amount formatting
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

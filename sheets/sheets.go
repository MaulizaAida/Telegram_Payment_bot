package sheets

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	//
	sheets "google.golang.org/api/sheets/v4"
)

type Transaction struct {
	Amount        int    `json:"amount"`
	Description   string `json:"description"`
	PaymentMethod string `json:"payment_method"`
	Category      string `json:"category"`
	Type          string `json:"type"` // "expense" or "income"
}

var lastID int // Package-level variable to track the last ID

// Update SaveToSheet to accept the service
func SaveToSheet(srv *sheets.Service, jsonStr string, username string) error {
	var tx Transaction
	err := json.Unmarshal([]byte(jsonStr), &tx)
	if err != nil {
		return fmt.Errorf("❌ JSON dari AI tidak valid: %v", err)
	}

	lastID++
	idString := fmt.Sprintf("%04d", lastID)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	spreadsheetId := os.Getenv("SHEET_ID")
	writeRange := "Sheet1!A:H" // Updated to include Type column

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{
		idString,
		timestamp,
		username,
		tx.Description,
		tx.PaymentMethod,
		tx.Category,
		tx.Amount,
		tx.Type, // Make sure this is exactly "income" or "expense"
	})

	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &vr).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("Gagal kirim ke Google Sheet: %v", err)
	}

	return nil
}

// Update GetSummary to use the correct type
func GetSummary(srv *sheets.Service, spreadsheetId string) (*Summary, error) {
	readRange := "Sheet1!A2:H" // Include type column
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca sheet: %v", err)
	}

	summary := &Summary{
		TotalIncome:  0,
		TotalExpense: 0,
		Balance:      0,
	}

	for _, row := range resp.Values {
		if len(row) < 7 { // Skip invalid rows
			continue
		}

		amount := 0
		switch v := row[6].(type) {
		case float64:
			amount = int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				amount = n
			}
		}

		txType := "expense" // default type
		if len(row) >= 8 {
			if t, ok := row[7].(string); ok {
				txType = t
			}
		}

		if txType == "income" {
			summary.TotalIncome += amount
		} else {
			summary.TotalExpense += amount
		}

		// Add to recent transactions (last 5)
		if len(summary.RecentTransactions) < 5 {
			tx := Transaction{
				Amount:        amount,
				Description:   fmt.Sprint(row[3]),
				PaymentMethod: fmt.Sprint(row[4]),
				Category:      fmt.Sprint(row[5]),
				Type:          txType,
			}
			summary.RecentTransactions = append([]Transaction{tx}, summary.RecentTransactions...)
		}
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}

// Add new struct for filtered transactions
type TransactionHistory struct {
	Transactions []Transaction
	TotalIncome  int
	TotalExpense int
	Balance      int
}

func GetHistory(srv *sheets.Service, spreadsheetId string, period string) (*TransactionHistory, error) {
	readRange := "Sheet1!A2:H"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca sheet: %v", err)
	}

	history := &TransactionHistory{
		Transactions: []Transaction{},
	}

	now := time.Now()
	var startTime time.Time

	switch period {
	case "day":
		startTime = now.Truncate(24 * time.Hour)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	default:
		startTime = now.AddDate(0, 0, -1) // default 1 day
	}

	for _, row := range resp.Values {
		if len(row) < 8 {
			continue
		}

		// Parse timestamp from column B
		timestamp, err := time.Parse("2006-01-02 15:04:05", fmt.Sprint(row[1]))
		if err != nil {
			continue
		}

		// Skip if transaction is before start time
		if timestamp.Before(startTime) {
			continue
		}

		amount := 0
		switch v := row[6].(type) {
		case float64:
			amount = int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				amount = n
			}
		}

		txType := fmt.Sprint(row[7])
		tx := Transaction{
			Amount:        amount,
			Description:   fmt.Sprint(row[3]),
			PaymentMethod: fmt.Sprint(row[4]),
			Category:      fmt.Sprint(row[5]),
			Type:          txType,
		}

		if txType == "income" {
			history.TotalIncome += amount
		} else {
			history.TotalExpense += amount
		}

		history.Transactions = append([]Transaction{tx}, history.Transactions...)
	}

	history.Balance = history.TotalIncome - history.TotalExpense
	return history, nil
}

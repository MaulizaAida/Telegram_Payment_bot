package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var paymentMethods = "[BCA, Jago, ShopeePay, Gopay, Cash]"
var expenseCategories = "[Makanan, Bahan Makanan, Transportasi, Belanja Harian, Belanja Online, Tagihan, Hiburan, Buah, Kesehatan]"
var incomeCategories = "[Gaji, Bonus, Investasi, Penjualan, Hadiah, Lainnya]"
var defaultMethod = "Cash"

func BuildPrompt(userInput string) string {
	return fmt.Sprintf(`
Kamu adalah asisten keuangan. Analisis input pengguna dan tentukan apakah itu pengeluaran atau pemasukan.
Kembalikan hasilnya dalam format JSON dengan kolom berikut:

- amount: jumlah dalam Rupiah (bilangan bulat) tanpa tanda mata uang
- description: deskripsi singkat transaksi
- payment_method: salah satu dari %s, jika tidak disebutkan gunakan %s
- category: 
    Untuk pengeluaran: %s
    Untuk pemasukan: %s
- type: "expense" untuk pengeluaran, "income" untuk pemasukan

Contoh Pengeluaran:
Input: "25K nasi goreng via ShopeePay"
Output:
{
  "amount": 25000,
  "description": "nasi goreng",
  "payment_method": "ShopeePay",
  "category": "Makanan",
  "type": "expense"
}

Contoh Pemasukan:
Input: "terima gaji 5jt via BCA"
Output:
{
  "amount": 5000000,
  "description": "gaji bulanan",
  "payment_method": "BCA",
  "category": "Gaji",
  "type": "income"
}

Sekarang, analisis input berikut dan balas hanya dengan JSON tanpa kode blok:
"%s"
`, paymentMethods, defaultMethod, expenseCategories, incomeCategories, userInput)
}

func CallOllama(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model":  "gemma3:1b",
		"prompt": prompt,
		"stream": false,
	}

	body, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 120 * time.Second}

	resp, err := client.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("❌ Gagal POST ke Ollama: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
	}

	resBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("❌ Ollama gagal (status %d): %s", resp.StatusCode, resBody)
	}

	if err := json.Unmarshal(resBody, &result); err != nil {
		return "", fmt.Errorf("❌ Gagal parse JSON dari Ollama: %s", resBody)
	}

	cleaned := cleanJSONResponse(result.Response)

	if cleaned == "" {
		return "", errors.New("❌ Response dari AI kosong")
	}

	return cleaned, nil
}

func PingOllama() error {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return errors.New("❌ Gagal terhubung ke Ollama. Pastikan Ollama aktif di localhost:11434")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("❌ Ollama merespons status bukan 200")
	}
	return nil
}

func cleanJSONResponse(raw string) string {
	// Hilangkan ```json ... ``` atau ``` ... ```
	re := regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")
	match := re.FindStringSubmatch(raw)
	if len(match) >= 2 {
		return match[1] // Ambil hanya isi JSON-nya
	}
	return raw // Kalau tidak dibungkus, kembalikan apa adanya
}

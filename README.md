Berikut adalah versi README.md yang sudah saya edit agar lebih rapi, terstruktur dengan jelas, dan enak dibaca saat di GitHub atau editor markdown lain:

````markdown
# Telegram Finance Bot

Bot Telegram untuk mencatat dan melacak keuangan pribadi dengan integrasi Google Sheets.  
Memudahkan pencatatan transaksi dengan bahasa natural, menyimpan data otomatis ke Google Sheets, serta menampilkan ringkasan dan riwayat transaksi dengan mudah.

---

## Fitur

- 📝 Pencatatan transaksi keuangan (pemasukan & pengeluaran) dengan bahasa natural  
- 📊 Ringkasan keuangan real-time (saldo, total pemasukan, total pengeluaran)  
- 📅 Riwayat transaksi dengan filter periode (hari ini, minggu ini, bulan ini)  
- 🤖 Pemrosesan input AI menggunakan Ollama untuk memahami bahasa natural  
- 📈 Integrasi Google Sheets untuk penyimpanan data yang aman dan mudah diakses  

---

## Perintah Bot

| Perintah | Deskripsi                                |
| -------- | --------------------------------------- |
| `/summary` | Menampilkan ringkasan keuangan & 5 transaksi terakhir |
| `/today`   | Menampilkan transaksi hari ini           |
| `/week`    | Menampilkan transaksi 7 hari terakhir    |
| `/month`   | Menampilkan transaksi 30 hari terakhir   |
| `/help`    | Menampilkan panduan penggunaan bot       |

---

## Format Input Transaksi

Bot memahami input transaksi dalam bahasa natural, contoh:

- `beli makan siang 25rb via gopay`
- `terima gaji 5jt via bca`
- `bayar listrik 200k cash`

### Format Nominal

| Notasi    | Contoh | Nilai           |
| --------- | -------|-----------------|
| `rb`/`ribu` | `25rb` | Rp 25.000       |
| `k`       | `100k` | Rp 100.000      |
| `jt`/`juta`| `5jt`  | Rp 5.000.000    |

### Metode Pembayaran

- BCA (Transfer Bank)  
- Jago (Digital Banking)  
- ShopeePay (E-wallet)  
- Gopay (E-wallet)  
- Cash (default jika tidak disebutkan)  

---

## Kategori Transaksi

### Pengeluaran

- Makanan (Restaurant)  
- Bahan Makanan (Grocery)  
- Transportasi (Bensin/Parkir/Online)  
- Belanja Harian (Kebutuhan Sehari-hari)  
- Belanja Online (E-commerce)  
- Tagihan (Listrik/Air/Internet)  
- Hiburan (Film/Game/Hobi)  
- Buah (Buah-buahan Segar)  
- Kesehatan (Obat/Vitamin)  

### Pemasukan

- Gaji  
- Bonus  
- Investasi  
- Penjualan  
- Hadiah  
- Lainnya  

---

## Prasyarat

- Go 1.23.4 atau lebih tinggi  
- Ollama untuk pemrosesan AI  
- Google Cloud Project dengan Sheets API aktif  
- Bot Telegram dengan token akses  

---

## Setup

1. **Clone repository:**

   ```bash
   git clone https://github.com/yourusername/telegram-history-bot.git
   cd telegram-history-bot
````

2. **Buat file `.env` dan isi:**

   ```env
   TELEGRAM_BOT_TOKEN=your_bot_token
   AUTHORIZED_USER_ID=your_telegram_user_id
   SHEET_ID=your_google_sheet_id
   ```

3. **Setup Google Cloud:**

   * Buat project baru di Google Cloud Console
   * Enable Google Sheets API
   * Buat Service Account dan download `credentials.json`
   * Share Google Sheet dengan email service account

4. **Install dependencies:**

   ```bash
   go mod download
   ```

5. **Setup Ollama:**

   * Download & install dari [https://ollama.ai](https://ollama.ai)
   * Jalankan service Ollama
   * Pull model yang dibutuhkan:

     ```bash
     ollama pull gemma:3b
     ```

6. **Jalankan bot:**

   ```bash
   go run main.go
   ```

---

## Struktur Google Sheet

Header sheet harus sesuai:

| Kolom          | Keterangan                               |
| -------------- | ---------------------------------------- |
| ID             | Auto-generated (format: 0001, 0002, dst) |
| Timestamp      | Format: `YYYY-MM-DD HH:mm:ss`            |
| Username       | Username Telegram atau ID                |
| Description    | Deskripsi transaksi                      |
| Payment Method | Metode pembayaran                        |
| Category       | Kategori transaksi                       |
| Amount         | Nominal dalam Rupiah                     |
| Type           | `income` / `expense`                     |

---

## Cara Kerja AI

Bot memakai model Gemma dari Ollama untuk:

1. Mengekstrak nominal transaksi dari teks input
2. Mengidentifikasi metode pembayaran
3. Mengkategorikan transaksi otomatis
4. Menentukan tipe transaksi (pemasukan atau pengeluaran)

---

Kalau mau aku buatkan juga contoh input dan output atau penjelasan tambahan, tinggal bilang ya!

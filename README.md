# Telegram Finance Bot

Bot Telegram untuk mencatat dan melacak keuangan pribadi dengan integrasi Google Sheets. Bot ini memungkinkan pengguna untuk mencatat transaksi keuangan dengan bahasa natural, menyimpan data secara otomatis ke Google Sheets, dan melihat ringkasan serta riwayat transaksi dengan mudah.

## Fitur

- 📝 Pencatatan transaksi keuangan (pemasukan & pengeluaran) dengan bahasa natural
- 📊 Ringkasan keuangan real-time (saldo, total pemasukan, total pengeluaran)
- 📅 Riwayat transaksi dengan filter periode (hari ini, minggu ini, bulan ini)
- 🤖 AI-powered input processing menggunakan Ollama untuk pemahaman bahasa natural
- 📈 Integrasi dengan Google Sheets untuk penyimpanan data yang aman dan mudah diakses

## Perintah Bot

- `/summary` - Menampilkan ringkasan keuangan & 5 transaksi terakhir untuk melihat status keuangan terkini
- `/today` - Menampilkan transaksi hari ini untuk tracking pengeluaran harian
- `/week` - Menampilkan transaksi 7 hari terakhir untuk analisis mingguan
- `/month` - Menampilkan transaksi 30 hari terakhir untuk evaluasi bulanan
- `/help` - Menampilkan panduan lengkap penggunaan bot

## Format Input Transaksi

Bot memahami input dalam bahasa natural untuk kemudahan pencatatan, contoh:
- "beli makan siang 25rb via gopay"
- "terima gaji 5jt via bca"
- "bayar listrik 200k cash"

### Format Nominal
Bot memahami beberapa format nominal untuk fleksibilitas input:
- rb/ribu: "25rb" = Rp 25.000
- k: "100k" = Rp 100.000
- jt/juta: "5jt" = Rp 5.000.000

### Metode Pembayaran
Tersedia beberapa opsi pembayaran untuk mencakup berbagai jenis transaksi:
- BCA (Transfer Bank)
- Jago (Digital Banking)
- ShopeePay (E-wallet)
- Gopay (E-wallet)
- Cash (default jika tidak disebutkan)

### Kategori
#### Pengeluaran:
- Makanan (Makanan jadi/Restaurant)
- Bahan Makanan (Grocery/Bahan Masak)
- Transportasi (Bensin/Parkir/Transportasi Online)
- Belanja Harian (Kebutuhan Sehari-hari)
- Belanja Online (E-commerce/Online Shopping)
- Tagihan (Listrik/Air/Internet)
- Hiburan (Film/Game/Hobi)
- Buah (Buah-buahan Segar)
- Kesehatan (Obat/Vitamin/Medical)

#### Pemasukan:
- Gaji (Pendapatan Bulanan)
- Bonus (Pendapatan Tambahan)
- Investasi (Return Investasi)
- Penjualan (Hasil Jualan)
- Hadiah (Gift/Bonus)
- Lainnya (Pendapatan Lain)

## Prerequisite

- Go 1.23.4 atau lebih tinggi untuk menjalankan bot
- Ollama untuk pemrosesan AI dan natural language
- Google Cloud Project dengan Sheets API enabled untuk penyimpanan data
- Bot Telegram untuk interface pengguna

## Setup

1. Clone repository:
   ```bash
   git clone https://github.com/yourusername/telegram-history-bot.git
   cd telegram-history-bot
   ```

2. Buat file `.env` dan isi dengan konfigurasi yang diperlukan:
   ```env
   TELEGRAM_BOT_TOKEN=your_bot_token
   AUTHORIZED_USER_ID=your_telegram_user_id
   SHEET_ID=your_google_sheet_id
   ```

3. Setup Google Cloud:
   - Buat project baru di Google Cloud Console
   - Enable Google Sheets API untuk akses spreadsheet
   - Buat Service Account untuk autentikasi
   - Download credentials.json untuk akses API
   - Share Google Sheet dengan email service account

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Setup Ollama:
   - Download dari https://ollama.ai
   - Install dan jalankan service
   - Pull model yang dibutuhkan:
     ```bash
     ollama pull gemma:3b
     ```

6. Jalankan bot:
   ```bash
   go run main.go
   ```

## Struktur Google Sheet

Sheet harus memiliki header berikut di baris pertama untuk tracking yang terorganisir:
- ID (Auto-generated, format: 0001, 0002, dst)
- Timestamp (Format: YYYY-MM-DD HH:mm:ss)
- Username (Username Telegram atau ID)
- Description (Deskripsi detail transaksi)
- Payment Method (Metode pembayaran yang digunakan)
- Category (Kategori untuk analisis pengeluaran)
- Amount (Nominal dalam Rupiah)
- Type (income/expense untuk klasifikasi)

## Cara Kerja AI

Bot menggunakan Ollama dengan model Gemma untuk pemrosesan bahasa natural:
1. Mengekstrak nominal transaksi dari input teks
2. Mengidentifikasi metode pembayaran yang digunakan
3. Mengkategorikan transaksi secara otomatis
4. Menentukan tipe transaksi (pemasukan/pengeluaran)

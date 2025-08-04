# Toko Biru - Backend E-commerce dengan Go & MongoDB

![Toko Biru Banner](https://placehold.co/1200x630/3B82F6/FFFFFF?text=Toko+Biru+Backend&font=inter)

**Toko Biru** adalah sebuah proyek yang berfokus pada pengembangan **backend REST API** yang tangguh dan kaya fitur untuk platform e-commerce modern. Dibangun menggunakan **Golang** dan **MongoDB**, backend ini menyediakan semua fungsionalitas inti yang dibutuhkan oleh sebuah toko online, mulai dari manajemen produk hingga integrasi dengan AI.

Proyek ini juga menyertakan **frontend demonstrasi** yang dibuat dengan **Vue.js** (via CDN) untuk menunjukkan bagaimana API dapat digunakan dalam aplikasi nyata. Namun, fokus utama dari repositori ini adalah pada arsitektur, logika bisnis, dan performa di sisi server.

Seluruh tumpukan aplikasi dikemas dalam kontainer **Docker** untuk kemudahan pengembangan dan deployment.

---

## ✨ Fitur Utama (Backend API)

### Untuk Pelanggan (Customer)
- **Registrasi & Login**: Sistem autentikasi aman menggunakan JWT.
- **Katalog Produk**: Endpoint untuk menampilkan semua produk dengan gambar dan harga.
- **Pencarian Produk**: API mendukung filter produk secara dinamis berdasarkan nama.
- **Halaman Detail Produk**: Endpoint untuk mengambil deskripsi lengkap, spesifikasi, dan FAQ produk.
- **Keranjang Belanja**: API untuk menambah, mengurangi, dan menghapus item di keranjang.
- **Proses Checkout**: Alur API untuk memproses pesanan dari keranjang.
- **Riwayat Pesanan**: API untuk melihat daftar semua transaksi yang pernah dilakukan beserta statusnya.
- **Chatbot AI**: Endpoint yang terintegrasi dengan Google Gemini untuk menjawab pertanyaan seputar produk.

### Untuk Administrator (Admin)
- **Dashboard Admin**: Kumpulan endpoint khusus untuk manajemen toko.
- **Manajemen Produk (CRUD)**: API untuk menambah, melihat, mengedit, dan menghapus produk.
- **Laporan Penjualan**: Endpoint agregasi untuk menghasilkan ringkasan performa toko, termasuk total pendapatan, jumlah pesanan, dan produk terlaris.
- **Manajemen Pesanan**: API untuk melihat semua pesanan dari pelanggan dan mengubah statusnya (misal: dari "baru" menjadi "dikirim").

---

## 🚀 Teknologi yang Digunakan

| Kategori | Teknologi |
| :--- | :--- |
| **Backend** | Golang (Go) 1.23+ dengan framework Gin |
| **Frontend (Demo)** | Vue.js 3 (via CDN) & Tailwind CSS |
| **Database** | MongoDB |
| **Kontainerisasi** | Docker & Docker Compose |
| **AI Chatbot** | Google Gemini API |
| **Autentikasi** | JSON Web Tokens (JWT) |

---

## ⚙️ Cara Menjalankan Proyek

Pastikan **Docker** dan **Docker Compose** sudah terpasang di sistem Anda.

### 1. Clone Repositori
```bash
git clone [https://github.com/URL_REPO_ANDA/tokobiru.git](https://github.com/URL_REPO_ANDA/tokobiru.git)
cd tokobiru
```

### 2. Konfigurasi Environment
Salin file konfigurasi contoh dan isi variabel yang diperlukan.
```bash
cp .env.example .env
```
Buka file `.env` yang baru dibuat dan masukkan **API Key Google Gemini** Anda:
```
GEMINI_API_KEY=MASUKKAN_API_KEY_ANDA_DI_SINI
```

### 3. Jalankan dengan Docker Compose
Perintah ini akan membangun image untuk backend Go, menarik image MongoDB, dan menjalankan semuanya.
```bash
docker-compose up --build
```
Biarkan terminal ini berjalan. Server backend akan aktif di `http://localhost:8080`.

### 4. Isi Data Awal (Seeder)
Buka **terminal baru**, masuk ke direktori proyek, dan jalankan perintah ini untuk mengisi database dengan data produk dan akun admin awal.
```bash
docker-compose exec go-app ./seeder
```

### 5. Buka Frontend (Demo)
Buka file `index.html` langsung di browser Anda. Aplikasi sekarang siap digunakan untuk berinteraksi dengan backend.

---

## 📂 Struktur Proyek (Backend)

```
.
├── controllers/    # Logika untuk menangani request HTTP
├── database/       # Koneksi ke MongoDB
├── middlewares/    # Middleware untuk autentikasi & otorisasi
├── models/         # Struct untuk data (User, Product, dll.)
├── routes/         # Definisi semua endpoint API
├── seed/           # Skrip untuk data awal
├── services/       # Logika bisnis (termasuk AI Chatbot)
├── .env            # File konfigurasi (diabaikan oleh Git)
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── main.go         # Entry point aplikasi Go

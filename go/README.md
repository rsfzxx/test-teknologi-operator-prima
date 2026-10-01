# Test Teknologi Operator Prima

# Database hanya untuk mempermudah keperluan testing/review, saya memahami apabila menyimpan database url dan upload ke github sangat tidak di anjurkan
Database URL = postgresql://neondb_owner:npg_4XtCeYmsy6fv@ep-polished-scene-b32c0s7x-pooler.c-4.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require

Berikut adalah panduan untuk melakukan test pada GO (Go + PostgreSQL).

## 1. Instalasi dan Menjalankan Aplikasi

### Prasyarat

- Go 1.22 atau lebih baru (cek dengan `go version`)
- Git
- Koneksi internet (database berada di Neon, cloud PostgreSQL)
- Postman atau tools sejenis untuk menguji API

### Langkah-langkah

**1.1 Clone repository dan masuk ke folder `go`**

```bash
git clone https://github.com/rsfzxx/test-teknologi-operator-prima.git
cd test-teknologi-operator-prima/go
```

**1.2 Buat file `.env`**

Salin `.env.example` menjadi `.env` dan masukan Database URL yang sudah saya berikan diatas

**1.3 Unduh dependency**

```bash
go mod tidy
```

**1.4 Jalankan aplikasi**
Setiap kali aplikasi dijalankan, aplikasi akan melakukan pengecekan migrasi database secara otomatis. Aplikasi memeriksa apakah terdapat migrasi baru (misalnya tabel baru, stored procedure, atau data awal) yang belum diterapkan, yang dicatat pada tabel schema_migrations, kemudian hanya menjalankan migrasi yang belum diterapkan tersebut. Apabila database sudah pernah dijalankan sebelumnya, migrasi yang telah diterapkan tidak dijalankan kembali sehingga data yang sudah ada tidak berubah. Dengan demikian, eksekusi SQL secara manual tidak diperlukan.

```bash
go run ./cmd/api
```

**1.5. Jalankan aplikasi menggunakan Air** (Optional)
Aplikasi dapat dijalankan menggunakan Air (live reload), sehingga aplikasi otomatis dimuat ulang setiap kali terdapat perubahan pada kode.
Install Air
```bash
go install github.com/air-verse/air@latest
```

**1.6 Pastikan aplikasi berjalan**

Buka `http://localhost:8080/health`. Jika berhasil, responsenya:

```json
{ "success": true, "message": "Service is healthy", "data": { "status": "healthy", "database": "up" } }
```

Aplikasi siap diuji lewat Postman dengan base URL `http://localhost:8080/api/v1`.

---

## 2. Daftar API

Base URL: `http://localhost:8080/api/v1`

Semua response berformat JSON.

- Sukses: `{ "success": true, "message": "...", "data": ... }`. Endpoint list juga memuat `meta` (`page`, `per_page`, `total`, `total_pages`).
- Gagal: `{ "success": false, "message": "...", "errors": [ { "field": "...", "message": "..." } ] }`.

### Health

| Method | Endpoint | Penjelasan |
|---|---|---|
| GET | `/health` | Memeriksa aplikasi dan koneksi database. Path ini berada di luar `/api/v1`. |

### Bank

| Method | Endpoint | Penjelasan |
|---|---|---|
| GET | `/banks` | Daftar bank dengan pagination. Query opsional: `page`, `per_page`, `search` (nama atau kode), `status` (`Active`/`Inactive`), `type`, `sort_by`, `order` (`asc`/`desc`). |
| GET | `/banks/{id}` | Detail satu bank berdasarkan UUID. |
| PATCH | `/banks/{id}/status` | Mengubah status bank. Body: `{ "status": "Active" }` atau `{ "status": "Inactive" }`. |

### Bank Account

| Method | Endpoint | Penjelasan |
|---|---|---|
| GET | `/bank-accounts` | Daftar rekening dengan pagination. Query opsional: `page`, `per_page`, `search`, `status` (`Accepted`/`Review`/`Rejected`), `bank_code`, `sort_by`, `order`. Data yang sudah di-soft delete tidak ditampilkan. |
| GET | `/bank-accounts/{id}` | Detail satu rekening berdasarkan UUID. Data yang sudah di-soft delete menghasilkan 404. |
| POST | `/bank-accounts` | Membuat rekening baru. |
| PUT | `/bank-accounts/{id}` | Mengubah seluruh data rekening. |
| DELETE | `/bank-accounts/{id}` | Menghapus rekening (soft delete). |

**POST `/bank-accounts`**

Body:

```json
{
  "bank": "<UUID dari tabel banks>",
  "account_number": "1234567890",
  "account_name": "Risman Muhammad"
}
```

- `bank` harus UUID bank yang valid. `bank_name` dan `bank_code` terisi otomatis dari bank tersebut.
- `account_number` hanya boleh angka.
- `status` otomatis `Review`, `reason` otomatis `null`, `created_at` dan `updated_at` otomatis waktu saat ini.
- Kombinasi `account_number` dan `bank_name` harus unik. Jika hanya salah satunya yang sama, tetap diperbolehkan. Jika kombinasinya sudah ada, response `409`.
- Jika kombinasinya sudah ada tetapi sudah ter-soft delete, tidak dibuat baris baru. Hanya `deleted_at` yang dikosongkan, response `200`.

**PUT `/bank-accounts/{id}`**

Body:

```json
{
  "bank": "<UUID dari tabel banks>",
  "account_number": "1234567890",
  "account_name": "Budi Santoso",
  "status": "Rejected",
  "reason": "Account name does not match registered identity"
}
```

- Semua field wajib, kecuali `reason`.
- `status` harus `Accepted`, `Review`, atau `Rejected`.
- `reason` wajib diisi jika `status` = `Rejected`, dan tidak boleh diisi jika `Accepted` atau `Review`.
- Mengganti `bank` otomatis mengganti `bank_name` dan `bank_code`.
- Aturan unik `account_number` + `bank_name` sama seperti POST (`409` jika bentrok dengan rekening lain).
- Rekening yang sudah ter-soft delete tidak bisa diubah: response `409` dengan pesan *There is an issue with your data; please contact our customer service to resolve the problem.*

**DELETE `/bank-accounts/{id}`**

- Hanya mengisi `deleted_at`. Data tetap ada di database tetapi tidak muncul di API.
- Rekening yang tidak ada atau sudah dihapus menghasilkan `404`.

### Kode status HTTP

| Kode | Arti |
|---|---|
| 200 | Berhasil (termasuk restore data soft delete) |
| 201 | Data berhasil dibuat |
| 400 | Request tidak valid: UUID salah, JSON rusak, atau query parameter tidak valid |
| 404 | Data atau endpoint tidak ditemukan |
| 405 | Method tidak didukung untuk endpoint tersebut |
| 409 | Data duplikat, atau data ter-soft delete tidak boleh diubah |
| 422 | Nilai field tidak valid (misalnya `status` salah atau `reason` tidak sesuai aturan) |
| 500 | Kesalahan server |
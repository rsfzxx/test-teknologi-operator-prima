# Test Teknologi Operator Prima — PHP Slim + MongoDB

REST API CRUD menggunakan **PHP Slim Framework 4** dan **MongoDB**.

## 1. Instalasi dan Menjalankan Aplikasi

### Prasyarat

- PHP 8.2 atau lebih baru (cek dengan `php -v`)
- Ekstensi PHP `mongodb` versi 2.x
- Composer
- Koneksi internet (database berada di MongoDB Atlas)
- MongoDB Compass (untuk melihat collection dan pipeline)
- Postman atau tools sejenis untuk menguji API

### Langkah-langkah

**1.1 Clone repository dan masuk ke folder `php-slim`**

```bash
git clone https://github.com/rsfzxx/test-teknologi-operator-prima.git
cd test-teknologi-operator-prima/php-slim
```

**1.2 Buat file `.env`**

Salin `.env.example` menjadi `.env`, lalu isi connection string MongoDB:

```env
APP_DEBUG=false
MONGODB_URI=mongodb+srv://rismanmuhammadhafidz21_db_user:TestingTop@test-top.ohbk3g2.mongodb.net/?appName=Test-top
MONGODB_DATABASE=operator_prima
```

**1.3 Unduh dependency**

```bash
composer install
```

**1.4 Jalankan migrasi**

```bash
composer migrate
```

Migrasi membuat collection `stores` beserta validator `$jsonSchema`, membuat index (termasuk unique index `store_uuid` dan `email`), lalu mengisi data awal. Migrasi yang sudah dijalankan dicatat di collection `migrations`, sehingga perintah ini aman dijalankan berkali-kali: hanya migrasi baru yang akan dieksekusi.

**1.5 Jalankan aplikasi**

```bash
composer start
```

Perintah ini menjalankan server bawaan PHP di `http://localhost:8080`

**1.6 Pastikan aplikasi berjalan**

Buka `http://localhost:8080/health`. Jika berhasil, responsenya:

```json
{ "success": true, "message": "Service is healthy", "data": { "status": "healthy", "database": "up" } }
```

---

## 2. Daftar API

Base URL: `http://localhost:8080/api/v1`

Setiap response memiliki header `X-Request-ID` untuk keperluan tracing log.

### Health

| Method | Endpoint | Penjelasan |
|---|---|---|
| GET | `/health` | Memeriksa aplikasi dan koneksi MongoDB. Path ini berada di luar `/api/v1`. |

### Store

| Method | Endpoint | Penjelasan |
|---|---|---|
| GET | `/stores` | Daftar store dengan pagination. Query opsional: `page`, `per_page` (maks. 100), `search` (nama, email, telepon, atau alamat), `status` (`Active`/`Inactive`), `sort_by` (`name`, `email`, `status`, `created_at`, `updated_at`), `order` (`asc`/`desc`). Default: terbaru di atas. |

| GET | `/stores/summary` | Ringkasan jumlah store per status, hasil aggregation pipeline. |

| GET | `/stores/{id}` | Detail satu store berdasarkan `store_uuid`. |

| POST | `/stores` | Membuat store baru. |

| PUT | `/stores/{id}` | Mengubah seluruh data store. |

| PATCH | `/stores/{id}/status` | Mengubah status store saja. |

| DELETE | `/stores/{id}` | Menghapus store (hard delete). |

**POST `/stores`** dan **PUT `/stores/{id}`**

Body:

```json
{
  "name": "Toko Sumber Rejeki",
  "address": "Jl. Kebon Jeruk Raya No. 12, Jakarta Barat",
  "email": "sumberrejeki@example.com",
  "phone": "081234567801",
  "status": "Active",
  "notes": "Buka setiap hari 08.00 - 21.00"
}
```

- Semua field wajib, kecuali `notes`.
- `store_uuid`, `created_at`, dan `updated_at` diisi otomatis.
- Pada PUT, seluruh data diganti. Jika `notes` tidak dikirim, nilainya menjadi `null`.
- Email harus unik. Jika sudah dipakai store lain, response `409`.
- Field di luar daftar di atas ditolak dengan response `400`.

**PATCH `/stores/{id}/status`**

Body:

```json
{ "status": "Inactive" }
```

**GET `/stores/summary`**

Contoh response:

```json
{
  "success": true,
  "message": "Store summary retrieved successfully",
  "data": {
    "total": 10,
    "by_status": [
      { "status": "Active", "total": 8, "percentage": 80, "latest_created_at": "2026-10-02T08:15:30.123+00:00" },
      { "status": "Inactive", "total": 2, "percentage": 20, "latest_created_at": "2026-10-02T08:15:30.120+00:00" }
    ]
  }
}
```

## 4. Aggregation Pipeline MongoDB

MongoDB tidak memiliki stored procedure seperti PostgreSQL. Sebagai ganti nya adalah **aggregation pipeline**.

Pipeline tersimpan di [`pipelines/store_summary_by_status.json`](pipelines/store_summary_by_status.json). File yang sama dibaca oleh endpoint `GET /api/v1/stores/summary`, sehingga pipeline yang dilampirkan identik dengan yang dijalankan aplikasi.

Tahapan pipeline:

1. **`$facet`** menjalankan dua sub-pipeline sekaligus:
   - `by_status`: `$group` per `status` untuk menghitung jumlah store dan `created_at` terbaru, lalu `$sort` berdasarkan status.
   - `overall`: `$count` total seluruh store.
2. **`$set`** mengambil angka total dari hasil `overall` (0 jika collection kosong).
3. **`$project`** dengan `$map` membentuk hasil akhir per status, termasuk persentase (`$divide`, `$multiply`, `$round`).

Cara menjalankannya di MongoDB Compass: buka collection `stores`, masuk ke tab **Aggregations**, pilih mode teks, lalu tempel isi file JSON tersebut.

---
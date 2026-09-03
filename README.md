# Rakamin Evermos: Backend E-Commerce API Service

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean_Architecture-blue?style=flat)
![Database](https://img.shields.io/badge/Database-MySQL-4479A1?style=flat&logo=mysql)
![Coverage](https://img.shields.io/badge/Tests-100%25_Passing-brightgreen?style=flat)

Backend Service E-Commerce berbasis **Go (Golang)** dengan penerapan **Clean Architecture**, dibangun sebagai tugas akhir Rakamin Academy × Evermos Virtual Internship.

Aplikasi ini mencakup fitur lengkap mulai dari Autentikasi, Pembuatan Toko Otomatis, Manajemen Produk (Multipart Upload), Transaksi Atomik (Row Locking & Immutable Product Logs), hingga Pengelolaan Alamat Standar Wilayah Indonesia.

---

## Daftar Isi
- [Tech Stack](#tech-stack)
- [Struktur Project](#struktur-project)
- [Fitur & Arsitektur](#fitur--arsitektur)
- [Database Schema (ERD)](#database-schema-erd)
- [Testing](#testing)
- [API Endpoints](#api-endpoints)
- [Build & Run Locally](#️-build-and-run-locally)
- [Dokumentasi Lengkap](#dokumentasi-lengkap)

---

## Tech Stack

| Komponen | Teknologi |
| :--- | :--- |
| **Language** | Go 1.25+ |
| **HTTP Framework** | Go Fiber v2 (`github.com/gofiber/fiber/v2`) |
| **Database** | MySQL 8.0 (Docker & Production) |
| **ORM** | GORM (`gorm.io/gorm` + `gorm.io/driver/mysql`) |
| **Auth & Security** | JWT (`golang-jwt/jwt/v5`) + Bcrypt Password Hashing |
| **API Docs** | Interactive Swagger UI / OpenAPI 3.0.3 (`/swagger`) |
| **API Client Docs**| Postman Collection v2.1.0 (Lengkap Bahasa Indonesia) |
| **Testing** | `testify/assert` + `testify/require` + Go Race Detector |
| **Administrative Standard** | Standar Kode & Wilayah Indonesia (Provinsi, Kabupaten, Kecamatan, Kelurahan) |
| **External API** | [EMSifa Wilayah Indonesia](https://emsifa.github.io/api-wilayah-indonesia/) (HTTP Client + In-Memory Cache) |

---

## Struktur Project

Penerapan **Clean Architecture** memisahkan *Business Logic* dengan kerangka eksternal secara konsisten:

```text
.
├── app/
│   └── main.go                     # Entry point + bootstrap server
├── bin/
│   └── evermos-api                 # Standalone compiled production binary
├── docs/
│   └── swagger.json                # Spesifikasi OpenAPI 3.0.3
├── internal/
│   ├── helper/                     # Formatter respon JSON terpadu & pagination
│   ├── infrastructure/
│   │   ├── container/              # Dependency Injection (DI) Container
│   │   └── mysql/                  # Connection pooling & GORM AutoMigrate
│   ├── pkg/
│   │   ├── entity/                 # GORM Domain Entities (7 model basis data)
│   │   ├── model/                  # Data Transfer Objects (DTOs Request & Response)
│   │   ├── repository/             # Abstraksi database & GORM persistence
│   │   └── usecase/                # Business logic & isolasi multi-tenant
│   ├── server/
│   │   └── http/
│   │       ├── handler/            # Fiber HTTP Controllers murni
│   │       ├── middleware/         # JWT Auth & Admin RBAC Middlewares
│   │       └── httproute.go        # Registrasi endpoint & rute
│   └── utils/                      # Bcrypt password hashing & JWT lifecycle
├── test/
│   └── e2e/                        # Folder khusus pengujian E2E & integrasi
│       ├── test_helper.go          # Centralized test runner & database auto-migrate
│       ├── auth_e2e_test.go
│       ├── user_store_e2e_test.go
│       ├── address_e2e_test.go
│       ├── category_e2e_test.go
│       ├── product_e2e_test.go
│       ├── transaction_e2e_test.go
│       ├── region_e2e_test.go
│       └── swagger_e2e_test.go
├── uploads/
│   └── products/                   # Penyimpanan fisik file gambar produk
├── Evermos_API.postman_collection.json # Koleksi Postman lengkap dengan dokumentasi ID
├── docker-compose.yaml             # Layanan MySQL 8.0 otomatis
├── .env.example                    # Template konfigurasi environment
├── go.mod & go.sum
└── README.md
```

> [!NOTE]
> **Clean Architecture Boundary:** Folder `internal/server/http/handler/` murni hanya menangani request/response HTTP. Seluruh validasi kepemilikan dan aturan bisnis berada di `internal/pkg/usecase/`, dan seluruh test integrasi terisolasi di folder `test/e2e/`.

---

## Fitur & Arsitektur

### Modul 1: Autentikasi & Registrasi
- **Register & Auto-Create Toko (Invarian #1)**: Dalam satu *Database Transaction* atomik, akun pengguna dibuat bersamaan dengan toko merchant bernama `"{UserName}'s Store"` (relasi 1:1).
- **Login JWT**: Kata sandi di-hash menggunakan `bcrypt`. Token JWT di-generate dengan masa berlaku terkonfigurasi.
- **Validasi Keunikan**: `email` dan `phone` wajib unik dan divalidasi ganda di level database dan aplikasi.

### Modul 2: Wilayah & Alamat Pengiriman
- **EMSifa External API Integration**: Mengambil master data resmi wilayah Indonesia secara langsung dari [EMSifa Wilayah Indonesia API](https://emsifa.github.io/api-wilayah-indonesia/) (`provinces`, `regencies`, `districts`, `villages`).
- **In-Memory Caching & Resilient Timeout**: Respon dari API eksternal disimpan di memory cache dengan proteksi `sync.RWMutex` untuk performa tinggi sub-milidetik, serta dilengkapi proteksi timeout (5 detik).
- **Standar Hierarki Wilayah Indonesia**: Menyimpan data resmi `provinsi`, `provinsi_id`, `kabupaten`, `kabupaten_id`, `kecamatan`, `kecamatan_id`, `kelurahan`, `kelurahan_id`.
- **Atomic Default Rebalancing**: Menandai sebuah alamat sebagai default (`is_default = true`) secara otomatis menonaktifkan status default pada alamat-alamat lama dalam satu transaksi database.
- **Isolasi Alamat**: Pengguna hanya dapat melihat, mengedit, dan menghapus alamat miliknya sendiri (`404 Not Found` jika mengakses alamat orang lain).

### Modul 3: Produk & Kategori
- **Multipart Upload**: Upload foto produk asli melalui `multipart/form-data` yang disimpan ke folder server `uploads/products/{userID}_{timestamp}_{filename}`.
- **Admin RBAC Kategori**: Manajemen kategori (`POST`, `PATCH`, `DELETE` `/categories`) diproteksi ketat hanya untuk Administrator (`is_admin == true`). Pengguna reguler ditolak dengan `403 Forbidden`.
- **Katalog Publik & Filtering**: Penelusuran produk publik dengan filter substring nama (`search`), filter kategori (`category_id`), filter toko (`store_id`), pengurutan harga (`sort=price_asc|price_desc|newest`), dan paginasi terpadu.
- **Isolasi Merchant**: Penjual hanya dapat memperbarui atau menghapus produk milik tokonya sendiri (`403 Forbidden` jika melanggar).

### Modul 4: Transaksi Atomik & Checkout Engine
- **Pessimistic Row Locking (`FOR UPDATE`)**: Mencegah *race condition* dan *overselling* stok pada saat lonjakan pesanan bersamaan.
- **Deadlock Prevention**: Pengurutan ID produk (*ascending*) sebelum penguncian baris database.
- **Validasi & Deduct Stok Otomatis**: Memastikan kuantitas stok mencukupi sebelum memotong stok secara instan di inventaris.
- **Invoice Number Generator**: Format unik otomatis `INV-{userID}-{timestamp}`.
- **Immutable Historical Snapshots (`product_logs`)**: Snapshot harga satuan dan kuantitas produk saat pembelian disalin secara permanen ke tabel `product_logs`. Perubahan harga produk di kemudian hari tidak merusak integritas invoice lama.

### Modul 5: Security & Multi-Tenancy
- **Zero-Trust Multi-Tenancy**: Data antar-pengguna terisolasi penuh. Seluruh operasi mengekstrak `userID` dari token JWT terverifikasi.
- **Soft Delete Terpusat**: Seluruh entitas menggunakan `gorm.DeletedAt` sehingga data historis tetap terlindungi.
- **Connection Pooling Teroptimasi**: Konfigurasi `MaxOpenConns`, `MaxIdleConns`, dan `ConnMaxLifetime` untuk efisiensi koneksi database.

---

## Database Schema (ERD)

```mermaid
erDiagram
    USERS ||--|| STORES : "owns (1:1)"
    USERS ||--o{ ADDRESSES : "has many"
    USERS ||--o{ TRANSACTIONS : "orders"
    STORES ||--o{ PRODUCTS : "sells"
    CATEGORIES ||--o{ PRODUCTS : "classifies"
    ADDRESSES ||--o{ TRANSACTIONS : "ships to"
    TRANSACTIONS ||--o{ PRODUCT_LOGS : "records snapshot"
    PRODUCTS ||--o{ PRODUCT_LOGS : "referenced by"

    USERS {
        uint id PK
        string email UK
        string phone UK
        string name
        string password
        bool is_admin
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    STORES {
        uint id PK
        uint user_id FK,UK
        string name
        string address
        string phone
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    ADDRESSES {
        uint id PK
        uint user_id FK
        string judul_alamat
        string penerima_nama
        string penerima_phone
        string provinsi
        string provinsi_id
        string kabupaten
        string kabupaten_id
        string kecamatan
        string kecamatan_id
        string kelurahan
        string kelurahan_id
        string detail_alamat
        bool is_default
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    CATEGORIES {
        uint id PK
        string name
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    PRODUCTS {
        uint id PK
        uint store_id FK
        uint category_id FK
        string name
        string description
        double price
        int quantity
        string image_url
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    TRANSACTIONS {
        uint id PK
        uint user_id FK
        uint address_id FK
        string invoice_number UK
        double total_price
        string status
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    PRODUCT_LOGS {
        uint id PK
        uint transaction_id FK
        uint product_id FK
        string product_name
        double price_at_purchase
        int quantity
        datetime created_at
    }
```

---

## Testing

Pengujian dilakukan secara komprehensif mencakup **Unit Testing** pada lapis domain & utilitas, serta **Integration / E2E Testing** di folder khusus `test/e2e/` menggunakan instance MySQL aktif.

```bash
# 1. Menjalankan seluruh pengujian di folder khusus test/e2e
go test -v ./test/e2e/...

# 2. Menjalankan seluruh test suite di proyek (Unit + E2E)
go test -v ./...

# 3. Menjalankan pengujian dengan Race Condition Detector
go test -v -race ./...

# 4. Memeriksa coverage pengujian
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## API Endpoints

Seluruh endpoint terpusat pada base path `/api/v1`:

<details>
<summary><b>Lihat Ringkasan Daftar Endpoint</b></summary>

### Public Routes
| Method | Rute | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Pendaftaran User & Auto-Create Toko |
| `POST` | `/api/v1/auth/login` | Mendapatkan JWT Bearer Token |
| `GET` | `/api/v1/stores` | List Toko Terdaftar (Paginasi) |
| `GET` | `/api/v1/stores/:id` | Detail Toko berdasarkan ID |
| `GET` | `/api/v1/categories` | List Kategori Produk (Paginasi) |
| `GET` | `/api/v1/categories/:id` | Detail Kategori berdasarkan ID |
| `GET` | `/api/v1/products` | Katalog Produk (Search, Filter, Sort) |
| `GET` | `/api/v1/products/:id` | Detail Produk berdasarkan ID |
| `GET` | `/api/v1/health` | Health Check Server |
| `GET` | `/swagger` | Antarmuka Interaktif Swagger UI |
| `GET` | `/swagger/doc.json` | Spesifikasi OpenAPI 3.0 (JSON) |

### Protected Routes (Bearer Token)
| Method | Rute | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/users/me` | Lihat profil sendiri |
| `PATCH`| `/api/v1/users/me` | Perbarui profil sendiri |
| `DELETE`| `/api/v1/users/me` | Hapus akun sendiri (*soft delete*) |
| `GET` | `/api/v1/users/:id` | Detail user berdasarkan ID |
| `GET` | `/api/v1/stores/me` | Lihat toko merchant sendiri |
| `PATCH`| `/api/v1/stores/:id` | Perbarui toko sendiri (Zero-Trust) |
| `DELETE`| `/api/v1/stores/:id` | Hapus toko sendiri |
| `POST` | `/api/v1/addresses` | Tambah alamat pengiriman baru |
| `GET` | `/api/v1/addresses` | Daftar alamat milik sendiri |
| `GET` | `/api/v1/addresses/:id` | Detail alamat sendiri |
| `PATCH`| `/api/v1/addresses/:id` | Perbarui alamat sendiri |
| `DELETE`| `/api/v1/addresses/:id` | Hapus alamat sendiri |
| `POST` | `/api/v1/products` | Tambah produk (Multipart Form Upload) |
| `GET` | `/api/v1/products/me` | Daftar produk milik toko sendiri |
| `PATCH`| `/api/v1/products/:id` | Perbarui produk sendiri |
| `DELETE`| `/api/v1/products/:id` | Hapus produk sendiri |
| `POST` | `/api/v1/transactions` | Checkout pesanan (ACID Transaction) |
| `GET` | `/api/v1/transactions/me` | Riwayat pesanan sendiri |
| `GET` | `/api/v1/transactions/me/:id`| Detail pesanan & log produk historis |
| `PATCH`| `/api/v1/transactions/me/:id`| Perbarui status pesanan (`completed`/`cancelled`) |

### Wilayah Indonesia (EMSifa External API)
| Method | Rute | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/regions/provinces` | List seluruh provinsi di Indonesia |
| `GET` | `/api/v1/regions/regencies/:province_id` | List kabupaten/kota berdasarkan ID provinsi |
| `GET` | `/api/v1/regions/districts/:regency_id` | List kecamatan berdasarkan ID kabupaten/kota |
| `GET` | `/api/v1/regions/villages/:district_id` | List kelurahan/desa berdasarkan ID kecamatan |

### Admin Routes (`is_admin == true`)
| Method | Rute | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/users` | List seluruh user terdaftar (Admin Only) |
| `POST` | `/api/v1/categories` | Tambah kategori baru (Admin Only) |
| `PATCH`| `/api/v1/categories/:id` | Perbarui kategori (Admin Only) |
| `DELETE`| `/api/v1/categories/:id` | Hapus kategori (Admin Only) |
| `GET` | `/api/v1/transactions` | List seluruh transaksi platform (Admin Only) |

</details>

---

## ⚙️ Build and Run Locally

Ikuti panduan langkah demi langkah berikut untuk mengompilasi dan menjalankan backend service di lingkungan lokal Anda:

### 1. Prasyarat Sistem (Prerequisites)
Sebelum memulai, pastikan perangkat lokal Anda telah memenuhi persyaratan berikut:
- **Go (Golang)**: Versi `1.25+` (`go version`)
- **Docker & Docker Compose**: Untuk menjalankan instance database MySQL 8.0
- **Git**: Untuk kloning repositori dan version control

### 2. Kloning Repositori & Masuk ke Folder Proyek
```bash
git clone https://github.com/VUXXE/Rakamim-Evermost.git
cd Rakamim-Evermost
```

### 3. Konfigurasi Environment Variables
Salin template konfigurasi `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

Contoh konfigurasi standar pada `.env`:
```ini
APP_NAME=Evermos-Ecommerce-API
APP_PORT=8080
APP_ENV=development

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=evermos

JWT_SECRET=super_secret_jwt_key_evermos_2026_safe_and_long!
JWT_EXP_HOURS=72
```

### 4. Setup Database MySQL (Docker Compose)
Jalankan layanan MySQL 8.0 secara otomatis di background:
```bash
docker compose up -d
```
> [!NOTE]
> Seluruh tabel basis data (`users`, `stores`, `addresses`, `categories`, `products`, `transactions`, `product_logs`) akan otomatis dibuat oleh **GORM AutoMigrate** saat aplikasi pertama kali dijalankan.

### 5. Unduh & Sinkronisasi Dependensi Go
```bash
go mod download
go mod tidy
```

### 6. Menjalankan Aplikasi (Development Mode)
Untuk menjalankan aplikasi secara langsung dari kode sumber:
```bash
go run ./app/main.go
```
*Server aktif dan melayani request HTTP pada: `http://localhost:8080`.*

### 7. Mengompilasi Binary Produksi (Build Standalone Binary)
Untuk mengompilasi backend menjadi single executable binary mandiri:
```bash
# Kompilasi kode Go ke dalam folder bin/
go build -o bin/evermos-api ./app/main.go

# Berikan izin eksekusi & jalankan binary
chmod +x bin/evermos-api
./bin/evermos-api
```

### 8. Memverifikasi Kesehatan Server (Health Check)
Verifikasi bahwa server berjalan optimal dengan memeriksa endpoint health:
```bash
curl -i http://localhost:8080/api/v1/health
```
Contoh respon HTTP 200 OK:
```json
{
  "code": 200,
  "message": "server is healthy and running",
  "data": {
    "status": "UP"
  }
}
```

### 9. Eksplorasi API via Swagger UI (Browser)
Buka browser dan akses antarmuka Swagger UI interaktif:
👉 **[http://localhost:8080/swagger](http://localhost:8080/swagger)** (atau `http://localhost:8080/docs`)

1. Eksekusi endpoint `POST /api/v1/auth/login` (atau `register`) untuk memperoleh token JWT.
2. Klik tombol hijau **Authorize 🔓** di pojok kanan atas Swagger UI.
3. Masukkan token dengan format: `Bearer {TOKEN}`.
4. Anda dapat menguji seluruh endpoint langsung dari browser!

### 10. Eksplorasi API via Postman Collection
Import file koleksi Postman yang telah disediakan:
- File koleksi: **`Evermos_API.postman_collection.json`**
- Buka Postman -> Klik **Import** -> Pilih file tersebut.
- Sudah terstruktur dalam 9 modul dengan dokumentasi lengkap berbahasa Indonesia pada tab **Docs / Documentation**.
- Variabel lingkungan `{{base_url}}`, `{{token}}`, dan `{{admin_token}}` telah terkonfigurasi otomatis.

---

## Dokumentasi Lengkap

| Dokumen | Lokasi File / URL |
| :--- | :--- |
| **Spesifikasi OpenAPI 3.0 (JSON)** | `docs/swagger.json` atau `/swagger/doc.json` |
| **Koleksi Postman v2.1.0** | `Evermos_API.postman_collection.json` |
| **Antarmuka Swagger UI** | `http://localhost:8080/swagger` |
| **Pengujian Integrasi E2E** | `test/e2e/` |
| **Entity Database Models** | `internal/pkg/entity/` |

---

## License

Didistribusikan di bawah lisensi MIT. Lihat file `LICENSE` untuk informasi lebih lanjut.

# Fullstack Starter — Go + Fiber, Next.js, PostgreSQL

Kerangka project fullstack monorepo sebagai fondasi aplikasi:

- **Backend** — REST API dengan Go + Fiber v3
- **Frontend** — Next.js 16 (App Router, TypeScript, Tailwind CSS)
- **Database** — PostgreSQL 18
- **Tooling** — Docker & Docker Compose

Project ini menyertakan satu contoh fitur CRUD (`items`) untuk membuktikan bahwa
frontend, backend, dan database sudah terhubung dengan benar. Fitur bisnis
sebenarnya dikerjakan di issue terpisah.

---

## Tech stack

| Bagian   | Teknologi / library                                      |
| -------- | -------------------------------------------------------- |
| Backend  | Go 1.26, Fiber v3, pgx v5, golang-migrate, log/slog      |
| Frontend | Next.js 16, React 19, TypeScript, Tailwind CSS 4         |
| Database | PostgreSQL 18                                            |
| Tooling  | Docker, Docker Compose, Make (opsional)                  |

---

## Struktur folder

```text
project-root/
├── backend/
│   ├── cmd/
│   │   ├── api/                 # entry point HTTP server
│   │   └── migrate/             # CLI migration (up / version)
│   ├── internal/
│   │   ├── config/              # load konfigurasi dari environment
│   │   ├── database/            # koneksi PostgreSQL + runner migration
│   │   ├── handler/             # HTTP handler, middleware, response envelope
│   │   ├── model/               # struct/entity
│   │   ├── repository/          # akses data (SQL)
│   │   └── service/             # logika bisnis + validasi
│   ├── migrations/              # file SQL migration (di-embed ke binary)
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── app/                     # App Router (layout, page, globals.css)
│   ├── components/              # komponen UI (items-manager)
│   ├── lib/                     # API client terpusat (lib/api.ts)
│   ├── public/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── .env.example
│   ├── package.json
│   └── next.config.ts
├── docker-compose.yml           # service db, backend, frontend
├── .env.example                 # template environment variable
├── .gitignore
├── Makefile                     # shortcut perintah umum (opsional)
├── README.md
└── issue.md
```

---

## Prasyarat

### Menjalankan dengan Docker Compose (disarankan)

- Docker Engine 24+ dan Docker Compose v2 (`docker compose version`)
- Tidak perlu memasang Go, Node.js, atau PostgreSQL di host

### Menjalankan secara lokal tanpa Docker

- Go 1.26+
- Node.js 24+ dan npm 11+
- PostgreSQL 18 yang sedang berjalan

---

## Menjalankan dengan Docker Compose

1. Salin template environment dan sesuaikan bila perlu:

   ```bash
   cp .env.example .env
   ```

   Windows PowerShell:

   ```powershell
   Copy-Item .env.example .env
   ```

2. Jalankan seluruh service:

   ```bash
   docker compose up --build
   ```

   Tunggu sampai healthcheck `db` dan `backend` hijau. Migration dijalankan
   **otomatis** saat backend start (`DB_AUTO_MIGRATE=true`).

3. Buka:

   - Frontend: <http://localhost:3000>
   - Health check API: <http://localhost:8080/api/health>
   - PostgreSQL: `localhost:5432` (kredensial dari `.env`)

4. Hentikan service:

   ```bash
   docker compose down
   ```

   Data PostgreSQL tetap tersimpan di volume `postgres_data`. Untuk menghapus
   data juga (destruktif):

   ```bash
   docker compose down -v
   ```

### Perintah Docker Compose yang sering dipakai

```bash
docker compose up -d --build     # start di background + rebuild image
docker compose logs -f backend   # lihat log backend
docker compose ps                # status container
docker compose down              # stop dan hapus container (data tetap)
docker compose run --rm --build backend /app/migrate up       # migration manual
docker compose run --rm backend /app/migrate version          # cek versi migration
```

Jika memakai `make`, lihat target yang tersedia dengan `make help`.

> **Catatan:** `NEXT_PUBLIC_API_URL` di-*inline* ke bundle browser saat build.
> Jika nilainya diubah, jalankan ulang `docker compose up --build` agar image
> frontend dibuat ulang.

---

## Menjalankan secara lokal tanpa Docker

### 1. Siapkan PostgreSQL

Pastikan PostgreSQL berjalan, lalu buat role dan database yang sesuai dengan
`.env` (default di `.env.example`):

```sql
CREATE ROLE app WITH LOGIN PASSWORD 'change-me';
CREATE DATABASE app OWNER app;
```

Atau ubah `DB_USER`, `DB_PASSWORD`, dan `DB_NAME` di `.env` sesuai database
yang sudah ada. `DB_HOST=localhost` dipakai untuk mode lokal ini (Docker
Compose akan meng-override menjadi `db`).

### 2. Jalankan backend

```bash
cd backend
go mod download
go run ./cmd/migrate up   # opsional; backend juga auto-migrate saat start
go run ./cmd/api
```

Backend membaca environment variable dan otomatis memuat `.env` dari folder
`backend/` atau root project. Environment variable yang sudah diset di shell
selalu menang atas isi `.env`.

API berjalan di <http://localhost:8080>. Cek dengan:

```bash
curl http://localhost:8080/api/health
```

### 3. Jalankan frontend

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env.local
```

Next.js hanya membaca `.env*` di folder `frontend/`, jadi `NEXT_PUBLIC_API_URL`
perlu ada di `frontend/.env.local` (bukan hanya di `.env` root).

Buka <http://localhost:3000>. Pastikan `CORS_ALLOWED_ORIGINS` di `.env` memuat
origin frontend (`http://localhost:3000` secara default).

---

## Environment variable

Semua konfigurasi Docker Compose dibaca dari file `.env` di root project.
`.env` sudah di-ignore oleh git; hanya `.env.example` yang di-commit.

| Variable                 | Default                 | Digunakan oleh        | Keterangan                                                                 |
| ------------------------ | ----------------------- | --------------------- | -------------------------------------------------------------------------- |
| `DB_HOST`                | `localhost`             | Backend               | Host PostgreSQL untuk mode lokal. Compose meng-override menjadi `db`.      |
| `DB_PORT`                | `5432`                  | Backend, host DB port | Port PostgreSQL. Di Compose juga menjadi port publik service `db`.         |
| `DB_NAME`                | `app`                   | Backend, DB container | Nama database.                                                             |
| `DB_USER`                | `app`                   | Backend, DB container | User database.                                                             |
| `DB_PASSWORD`            | `change-me`             | Backend, DB container | Password database. **Ganti sebelum dipakai di luar lokal.**                |
| `DB_SSLMODE`             | `disable`               | Backend               | Mode SSL PostgreSQL (`disable`, `require`, dll).                           |
| `DB_MAX_OPEN_CONNS`      | `25`                    | Backend               | Maksimum koneksi terbuka di pool.                                          |
| `DB_MAX_IDLE_CONNS`      | `25`                    | Backend               | Maksimum koneksi idle di pool.                                             |
| `DB_CONN_MAX_LIFETIME`   | `5m`                    | Backend               | Umur maksimum satu koneksi (format durasi Go).                             |
| `APP_ENV`                | `development`           | Backend               | `production` menyembunyikan startup message Fiber.                         |
| `APP_PORT`               | `8080`                  | Backend               | Port HTTP API.                                                             |
| `CORS_ALLOWED_ORIGINS`   | `http://localhost:3000` | Backend               | Daftar origin dipisah koma yang boleh memanggil API dari browser.          |
| `DB_AUTO_MIGRATE`        | `true`                  | Backend               | Jalankan migration otomatis saat backend start.                            |
| `NEXT_PUBLIC_API_URL`    | `http://localhost:8080` | Frontend (build time) | Base URL API yang dipanggil browser. Harus bisa diakses dari browser.      |
| `FRONTEND_PORT`          | `3000`                  | Docker Compose        | Port publik frontend.                                                      |

Jika `APP_PORT` atau `FRONTEND_PORT` diubah, sesuaikan juga
`NEXT_PUBLIC_API_URL` dan `CORS_ALLOWED_ORIGINS`, lalu rebuild image frontend.

---

## Endpoint API

Base URL: `http://localhost:8080`

| Method   | Path               | Deskripsi                        | Sukses |
| -------- | ------------------ | -------------------------------- | ------ |
| `GET`    | `/api/health`      | Status API + koneksi database    | `200`  |
| `GET`    | `/api/items`       | Daftar semua item (terbaru dulu) | `200`  |
| `GET`    | `/api/items/:id`   | Detail satu item                 | `200`  |
| `POST`   | `/api/items`       | Buat item baru                   | `201`  |
| `PUT`    | `/api/items/:id`   | Update item (full replace)       | `200`  |
| `DELETE` | `/api/items/:id`   | Hapus item                       | `200`  |

### Format response

Semua response memakai envelope JSON yang konsisten.

Sukses:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Contoh item",
    "description": "Deskripsi singkat",
    "created_at": "2026-01-01T10:00:00Z",
    "updated_at": "2026-01-01T10:00:00Z"
  }
}
```

Error:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "validation failed",
    "details": {
      "name": "name is required"
    }
  }
}
```

Kode error yang dipakai:

| HTTP | Code                   | Kapan terjadi                                  |
| ---- | ---------------------- | ---------------------------------------------- |
| 400  | `INVALID_REQUEST`      | ID bukan bilangan positif / body bukan JSON valid |
| 404  | `NOT_FOUND`            | Item atau route tidak ditemukan                |
| 405  | `METHOD_NOT_ALLOWED`   | Method tidak didukung untuk route tersebut     |
| 422  | `VALIDATION_ERROR`     | Input tidak lolos validasi                     |
| 500  | `INTERNAL_ERROR`       | Kesalahan tak terduga (detail hanya di log)    |
| 503  | `DATABASE_UNAVAILABLE` | Health check gagal menghubungi database        |

### Aturan validasi `items`

- `name` wajib diisi, maksimal 255 karakter.
- `description` opsional, maksimal 2000 karakter.
- Spasi di awal/akhir `name` dan `description` dipangkas.

### Contoh request

```bash
# List
curl http://localhost:8080/api/items

# Create
curl -X POST http://localhost:8080/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Contoh item","description":"Dibuat lewat curl"}'

# Update
curl -X PUT http://localhost:8080/api/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Nama baru","description":"Deskripsi baru"}'

# Delete
curl -X DELETE http://localhost:8080/api/items/1
```

---

## Database & migration

Migration memakai [golang-migrate](https://github.com/golang-migrate/migrate)
dan file SQL di `backend/migrations/` di-*embed* ke binary backend
(`//go:embed`), sehingga tidak perlu mengirim folder migration terpisah ke
container.

- File migration mengikuti pola `NNNNNN_nama.up.sql` dan `NNNNNN_nama.down.sql`.
- Saat backend start dan `DB_AUTO_MIGRATE=true`, semua migration yang belum
  dijalankan akan diterapkan secara berurutan.
- Riwayat versi disimpan di tabel `schema_migrations`.
- Migration manual:

  ```bash
  # Docker
  docker compose run --rm --build backend /app/migrate up
  docker compose run --rm backend /app/migrate version

  # Lokal
  cd backend
  go run ./cmd/migrate up
  go run ./cmd/migrate version
  ```

### Menambah migration baru

1. Buat dua file berurutan, misalnya:

   ```text
   backend/migrations/000002_add_status_to_items.up.sql
   backend/migrations/000002_add_status_to_items.down.sql
   ```

2. Tulis perubahan di `.up.sql` dan rollback-nya di `.down.sql`.
3. Rebuild/restart backend. Migration baru otomatis diterapkan.

Jika migration gagal, backend berhenti dengan error. Perbaiki SQL-nya, lalu
jalankan ulang. Untuk mengulang dari database kosong:

```bash
docker compose down -v
docker compose up --build
```

> Migration pertama (`000001_create_items`) membuat tabel `items` dengan kolom
> `id`, `name`, `description`, `created_at`, dan `updated_at`.

---

## Keputusan teknis

- **Fiber v3** untuk HTTP framework, dengan middleware `recover`, `logger`, dan
  `cors`.
- **pgx v5 (stdlib driver)** sebagai driver PostgreSQL, dipakai lewat
  `database/sql`.
- **golang-migrate** dengan source `iofs` agar migration ter-embed di binary.
- **log/slog** untuk logging terstruktur (tanpa library tambahan).
- **godotenv** untuk memuat `.env` saat pengembangan lokal. Di container,
  environment variable dari Compose yang dipakai.
- **Validasi manual** di service layer agar dependensi tetap minimal dan pesan
  error konsisten.
- **Response envelope** `{ success, data, error }` untuk semua endpoint.
- **Graceful shutdown** memakai `signal.NotifyContext` + `GracefulContext`
  Fiber dengan timeout 10 detik.
- **Tailwind CSS 4** untuk styling sederhana; tidak ada state management atau
  data-fetching library tambahan.
- **Next.js `output: "standalone"`** agar image Docker runtime kecil; server
  dijalankan dengan `node server.js`.
- **`NEXT_PUBLIC_API_URL`** adalah variabel publik yang di-inline saat build,
  bukan dibaca saat runtime container.
- **Module path** backend memakai placeholder
  `github.com/example/fullstack-starter/backend`. Ganti dengan path repository
  sebenarnya bila project ini dipakai untuk produksi.

---

## Troubleshooting

**Port sudah dipakai**

Ubah `DB_PORT`, `APP_PORT`, atau `FRONTEND_PORT` di `.env`. Di Docker Compose,
`DB_PORT` hanya mengatur port publik di host; backend tetap terhubung ke
`db:5432` di jaringan internal Compose. Jika mengubah `APP_PORT`/`FRONTEND_PORT`,
sesuaikan juga `NEXT_PUBLIC_API_URL` dan `CORS_ALLOWED_ORIGINS`, lalu rebuild
frontend.

**Frontend menampilkan "Could not reach the API"**

- Pastikan backend berjalan: `curl http://localhost:8080/api/health`.
- Pastikan `NEXT_PUBLIC_API_URL` menunjuk ke URL yang bisa diakses browser.
- Setelah mengubah `NEXT_PUBLIC_API_URL`, rebuild frontend
  (`docker compose up --build` atau `npm run build`).

**CORS error di browser**

`CORS_ALLOWED_ORIGINS` harus sama persis dengan origin browser, termasuk
skema, host, dan port (misalnya `http://localhost:3000`). Beberapa origin
dipisahkan koma.

**Password database berubah tetapi container gagal login**

`POSTGRES_USER`, `POSTGRES_PASSWORD`, dan `POSTGRES_DB` hanya dipakai saat
volume database pertama kali diinisialisasi. Untuk lokal, opsi termudah adalah
`docker compose down -v` (menghapus data) lalu start ulang. Untuk data yang
harus dipertahankan, ubah password role lewat `psql`.

**`next start` memberi peringatan `output: standalone`**

Itu normal untuk Next.js dengan `output: "standalone"`. Untuk pengembangan
lokal gunakan `npm run dev`. Untuk menjalankan hasil build secara lokal:

```bash
cd frontend
npm run build
cp -r .next/static .next/standalone/.next/
cp -r public .next/standalone/
node .next/standalone/server.js
```

Di Docker, langkah tersebut sudah dilakukan oleh `frontend/Dockerfile`.

---

## Di luar scope

- Autentikasi dan otorisasi
- Fitur bisnis spesifik di luar contoh `items`
- Deployment production dan CI/CD
- Unit/integration test (dikerjakan di issue terpisah)

---

## Keamanan

- Jangan commit file `.env`; hanya `.env.example` yang masuk repository.
- Ganti `DB_PASSWORD` default sebelum menjalankan di lingkungan bersama.
- Endpoint API saat ini terbuka tanpa autentikasi (sesuai scope scaffold).

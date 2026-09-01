# 📖 Dokumentasi API — Expense Tracker (Manajemen Keuangan Mahasiswa)

> ⚠️ **Status OAuth**: Google OAuth **belum dipakai sementara**. Untuk dapat token, pakai endpoint `POST /auth/dev-login` (lihat bagian 2). Endpoint ini nanti **dihapus/dinonaktifkan** setelah Google OAuth selesai di-setup dan dites — jangan sampai kebawa ke production.

**Base URL**: `http://localhost:8080`

---

## 1. Format Response

Semua response API berbentuk sama:

**Sukses**
```json
{
  "success": true,
  "message": "berhasil mengambil data transaksi",
  "data": { ... }
}
```

**Gagal**
```json
{
  "success": false,
  "message": "penjelasan errornya"
}
```

Frontend cukup selalu cek `success` untuk tahu request berhasil atau tidak, tanpa perlu mengandalkan status code HTTP secara ketat (walau status code tetap dikirim dengan benar: 200/201 sukses, 400/401/404 gagal).

---

## 2. Autentikasi

### 2.1 Dev Login (SEMENTARA — dipakai sekarang)

```
POST /auth/dev-login
Content-Type: application/json
```
**Body:**
```json
{
  "email": "mahasiswa@example.com",
  "name": "Budi Santoso"
}
```
**Response 200:**
```json
{
  "success": true,
  "message": "dev login berhasil (mode testing, bukan Google OAuth)",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "mahasiswa@example.com",
      "name": "Budi Santoso",
      "avatar": "",
      "created_at": "2026-08-31T10:00:00Z",
      "updated_at": "2026-08-31T10:00:00Z"
    }
  }
}
```
Simpan `data.token` — dipakai untuk semua request ke endpoint yang butuh login.

Kalau email yang sama dipanggil lagi, user yang sama akan dipakai (bukan bikin baru) — jadi aman dipanggil berulang kali buat testing.

### 2.2 Google OAuth (belum dipakai, siap dipasang nanti)

```
GET /auth/google/login       → redirect ke consent screen Google (buka di browser, bukan Postman)
GET /auth/google/callback    → dipanggil otomatis oleh Google setelah user setuju
```
Hasil akhirnya sama: dapat `token` JWT. Setelah ini aktif, `dev-login` dihapus dan frontend tinggal arahkan tombol "Login with Google" ke `GET /auth/google/login`.

### 2.3 Cara pakai token

Setiap request ke endpoint yang butuh login **wajib** kirim header:
```
Authorization: Bearer <token>
```
Kalau tidak ada / salah / kedaluwarsa → response `401 Unauthorized`.

### 2.4 Get Profil Sendiri

```
GET /api/auth/me
Authorization: Bearer <token>
```
**Response 200:**
```json
{
  "success": true,
  "message": "berhasil mengambil profil",
  "data": {
    "id": 1,
    "email": "mahasiswa@example.com",
    "name": "Budi Santoso",
    "avatar": "",
    "created_at": "2026-08-31T10:00:00Z",
    "updated_at": "2026-08-31T10:00:00Z"
  }
}
```

---

## 3. Transaksi (Pemasukan / Pengeluaran)

Semua endpoint di bawah ini **wajib** header `Authorization: Bearer <token>`.

### 3.1 Tambah Transaksi

```
POST /api/transactions
Content-Type: application/json
```
**Body:**
```json
{
  "type": "expense",
  "category": "makan",
  "amount": 25000,
  "description": "makan siang warteg",
  "date": "2026-08-31T12:00:00Z"
}
```
| Field | Tipe | Wajib? | Keterangan |
|---|---|---|---|
| `type` | string | ✅ | `"income"` atau `"expense"` saja |
| `category` | string | ✅ | 2–50 karakter, bebas (mis. "makan", "transport", "uang saku") |
| `amount` | number | ✅ | harus > 0 |
| `description` | string | ❌ | maks 255 karakter |
| `date` | string (ISO 8601) | ❌ | kalau kosong, otomatis pakai waktu sekarang |

**Response 201:**
```json
{
  "success": true,
  "message": "transaksi berhasil dibuat",
  "data": {
    "id": 1,
    "user_id": 1,
    "type": "expense",
    "category": "makan",
    "amount": 25000,
    "description": "makan siang warteg",
    "date": "2026-08-31T12:00:00Z",
    "created_at": "2026-08-31T10:01:00Z",
    "updated_at": "2026-08-31T10:01:00Z"
  }
}
```
**Response 400** (validasi gagal), contoh kalau `amount` dikirim 0 atau negatif:
```json
{ "success": false, "message": "input tidak valid: Key: 'TransactionInput.Amount' Error:Field validation for 'Amount' failed on the 'gt' tag" }
```

### 3.2 Ambil Semua Transaksi (dengan filter & pagination)

```
GET /api/transactions?type=expense&category=makan&start_date=2026-08-01&end_date=2026-08-31&page=1&limit=10
```
Semua query param **opsional**:

| Param | Contoh | Keterangan |
|---|---|---|
| `type` | `income` / `expense` | filter jenis |
| `category` | `makan` | filter kategori persis |
| `start_date` | `2026-08-01` | transaksi >= tanggal ini |
| `end_date` | `2026-08-31` | transaksi <= tanggal ini |
| `page` | `1` | default 1 |
| `limit` | `10` | default 10, maks 100 |

**Response 200:**
```json
{
  "success": true,
  "message": "berhasil mengambil data transaksi",
  "data": {
    "transactions": [ { "id": 1, "type": "expense", "category": "makan", "amount": 25000, ... } ],
    "pagination": { "page": 1, "limit": 10, "total": 24 }
  }
}
```
`pagination.total` = jumlah total data (dipakai frontend untuk hitung jumlah halaman).

### 3.3 Ambil Satu Transaksi

```
GET /api/transactions/:id
```
**Response 200** → objek transaksi tunggal.
**Response 404** kalau id tidak ada / bukan milik user yang login.

### 3.4 Update Transaksi

```
PUT /api/transactions/:id
Content-Type: application/json
```
**Body** sama seperti create (semua field dikirim ulang, bukan partial update):
```json
{
  "type": "expense",
  "category": "makan",
  "amount": 30000,
  "description": "makan siang warteg (revisi harga)"
}
```
**Response 200** → objek transaksi setelah diupdate.
**Response 404** kalau id tidak ditemukan.

### 3.5 Hapus Transaksi

```
DELETE /api/transactions/:id
```
**Response 200:**
```json
{ "success": true, "message": "transaksi berhasil dihapus", "data": null }
```

### 3.6 Ringkasan Keuangan

```
GET /api/transactions/summary
```
**Response 200:**
```json
{
  "success": true,
  "message": "berhasil mengambil ringkasan",
  "data": {
    "total_income": 1000000,
    "total_expense": 275000,
    "balance": 725000
  }
}
```
> Catatan: summary ini total keseluruhan (belum bisa difilter per bulan). Kalau perlu per periode, tambahkan query param `start_date`/`end_date` seperti di endpoint list — ini salah satu latihan lanjutan yang bisa kamu tambahkan sendiri.

---

## 4. Testing dengan Postman

1. Import file **`postman_collection.json`** (Postman → Import → pilih file).
2. Collection sudah punya variable `base_url` (default `http://localhost:8080`) dan `token` (kosong).
3. Jalankan server: `go run main.go`.
4. Buka folder **Auth → Dev Login**, klik **Send**. Token otomatis kesimpan ke variable `{{token}}` (ada script kecil di tab *Tests* request ini).
5. Buka folder **Transactions**, semua request di situ sudah otomatis pakai `{{token}}` di header — tinggal klik Send satu-satu.
6. Request **Create Transaction** otomatis menyimpan `id` hasil create ke variable `{{transaction_id}}`, jadi request Get/Update/Delete by ID otomatis nyambung tanpa perlu copy-paste id manual.

Urutan testing yang disarankan: `Dev Login` → `Create Transaction` → `Get All Transactions` → `Get Summary` → `Update Transaction` → `Delete Transaction`.

---

## 5. Integrasi Frontend (JavaScript / React)

### 5.1 Setup client fetch sederhana

```js
// api.js
const BASE_URL = "http://localhost:8080";

function getToken() {
  return localStorage.getItem("token");
}

async function apiFetch(path, options = {}) {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(getToken() ? { Authorization: `Bearer ${getToken()}` } : {}),
      ...options.headers,
    },
  });

  const json = await res.json();
  if (!json.success) {
    throw new Error(json.message || "Terjadi kesalahan");
  }
  return json.data;
}

export default apiFetch;
```

### 5.2 Login (dev-login, sementara)

```js
import apiFetch from "./api";

async function login(email, name) {
  const data = await apiFetch("/auth/dev-login", {
    method: "POST",
    body: JSON.stringify({ email, name }),
  });
  localStorage.setItem("token", data.token);
  return data.user;
}
```

> Kalau nanti Google OAuth aktif: ganti fungsi `login()` di atas — tombol "Login with Google" cukup `window.location.href = "http://localhost:8080/auth/google/login"`. Setelah user selesai login, backend akan redirect balik ke `FRONTEND_URL` (diatur di `.env` backend) dengan token di query string, contoh: `http://localhost:3000/auth/callback?token=xxx`. Frontend tinggal baca `token` dari URL di halaman `/auth/callback`, simpan ke `localStorage`, redirect ke halaman utama.

### 5.3 CRUD Transaksi

```js
import apiFetch from "./api";

export const getTransactions = (params = {}) => {
  const query = new URLSearchParams(params).toString();
  return apiFetch(`/api/transactions?${query}`);
};

export const getSummary = () => apiFetch("/api/transactions/summary");

export const createTransaction = (payload) =>
  apiFetch("/api/transactions", {
    method: "POST",
    body: JSON.stringify(payload),
  });

export const updateTransaction = (id, payload) =>
  apiFetch(`/api/transactions/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });

export const deleteTransaction = (id) =>
  apiFetch(`/api/transactions/${id}`, { method: "DELETE" });
```

### 5.4 Contoh pemakaian di komponen

```js
async function loadDashboard() {
  const { transactions, pagination } = await getTransactions({ page: 1, limit: 10 });
  const summary = await getSummary();
  console.log(transactions, pagination, summary);
}

async function tambahPengeluaran() {
  await createTransaction({
    type: "expense",
    category: "transport",
    amount: 15000,
    description: "ojek online",
  });
}
```

### 5.5 CORS

Karena frontend (misal `localhost:3000`) beda origin dari backend (`localhost:8080`), browser akan **blokir** request kalau backend belum mengizinkan CORS. Kalau nanti muncul error CORS di console browser, tambahkan middleware CORS di `routes.go`:

```bash
go get github.com/gin-contrib/cors
```
```go
import "github.com/gin-contrib/cors"

r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
}))
```
Taruh ini di `routes.Setup()`, sebelum route lain didaftarkan.

---

## 6. Checklist sebelum aktifkan Google OAuth beneran

- [ ] Hapus/nonaktifkan route `POST /auth/dev-login` di `routes.go`
- [ ] Hapus function `DevLogin` di `auth_handler.go` (atau bungkus dengan pengecekan `if gin.Mode() == gin.DebugMode` kalau mau tetap disimpan untuk dev)
- [ ] Pastikan `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URL` di `.env` sudah benar dan cocok dengan yang didaftarkan di Google Cloud Console
- [ ] Update dokumen ini: pindahkan bagian 2.2 (Google OAuth) jadi cara login utama, hapus bagian 2.1

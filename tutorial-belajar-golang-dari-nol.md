# 🎓 Tutorial dari NOL: Golang → REST API Manajemen Keuangan Mahasiswa

Tutorial ini untuk kamu yang **baru pertama kali belajar Go**. Kita mulai dari instal, konsep dasar bahasa, sampai akhirnya jadi REST API CRUD dengan Gin + SQLite + Google OAuth. Setiap langkah **bisa langsung dijalankan** — jangan loncat, ketik ulang sendiri kodenya (jangan copy-paste) biar nempel di otak.

Kerjakan semua ini di Ubuntu kamu, buka terminal.

---

## BAB 0 — Instalasi

### 0.1 Install Go
```bash
# cek dulu siapa tau sudah ada
go version

# kalau belum ada, install versi terbaru dari situs resmi (lebih baik daripada apt, karena apt sering versi lama)
cd /tmp
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# tambahkan ke PATH permanen
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

go version
# harus muncul: go version go1.22.5 linux/amd64
```

### 0.2 Install tools pendukung
```bash
sudo apt update
sudo apt install sqlite3 sqlitebrowser -y
```

### 0.3 Editor
Pakai **VS Code**, install extension **"Go"** (by Google) — nanti otomatis nawarin install tools tambahan (gopls, dlv, dll), klik **Install All**.

---

## BAB 1 — Konsep Dasar Go (wajib paham sebelum lanjut)

Buat folder latihan terpisah dulu, jangan campur sama project utama:
```bash
mkdir -p ~/belajar-go/dasar
cd ~/belajar-go/dasar
go mod init belajar
```

`go mod init belajar` membuat file `go.mod` — ini semacam "KTP" project Go, isinya nama module & versi Go yang dipakai. Setiap project Go **wajib** punya ini.

### 1.1 Hello World
Buat file `main.go`:
```go
package main // setiap file yang bisa dieksekusi harus punya package main

import "fmt" // import package bawaan buat print ke layar

func main() { // main() adalah titik masuk program, wajib ada persis nama ini
	fmt.Println("Halo, aku belajar Go!")
}
```
Jalankan:
```bash
go run main.go
```
> `go run` = compile + langsung jalankan, cocok buat development. Nanti kalau mau bikin file binary jadi, pakai `go build`.

### 1.2 Variabel & Tipe Data
```go
package main

import "fmt"

func main() {
	var nama string = "Budi"     // deklarasi eksplisit
	umur := 20                   // ":=" short declaration, tipe otomatis dideteksi (int)
	var ipk float64 = 3.75
	aktif := true                // bool

	fmt.Println(nama, umur, ipk, aktif)
	fmt.Printf("Nama: %s, Umur: %d, IPK: %.2f\n", nama, umur, ipk)
}
```
Tipe data penting di Go: `string`, `int`, `float64`, `bool`. Go itu **statically typed** — sekali variabel dideklarasi tipe X, seumur hidup dia tipe X (beda dari JS/Python).

### 1.3 Function
```go
func tambah(a int, b int) int { // (parameter) tipe-return
	return a + b
}

func bagiInfo(nama string, umur int) (string, int) { // bisa return lebih dari satu nilai
	return nama, umur
}
```
Ini penting karena di Go, **pola return `(value, error)` dipakai di HAMPIR SEMUA tempat** — nanti kamu akan sering lihat `data, err := someFunction()`.

### 1.4 Error Handling — pola paling khas di Go
Go **tidak punya try/catch**. Errornya dikembalikan sebagai value biasa, lalu **wajib dicek manual**:
```go
package main

import (
	"errors"
	"fmt"
)

func bagi(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("tidak bisa bagi dengan nol")
	}
	return a / b, nil // nil = "tidak ada error"
}

func main() {
	hasil, err := bagi(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Hasil:", hasil)
}
```
Pola `if err != nil { ... }` ini akan kamu ketik **ratusan kali** selama belajar Go. Wajar, itu memang gayanya Go.

### 1.5 Struct — pengganti "class"
Go tidak punya class seperti Java/PHP. Sebagai gantinya pakai `struct` (kumpulan field) + function terpisah:
```go
type Mahasiswa struct {
	Nama string
	NIM  string
	IPK  float64
}

func main() {
	m := Mahasiswa{Nama: "Sari", NIM: "12345", IPK: 3.8}
	fmt.Println(m.Nama, m.IPK)
}
```
Nanti model `User` dan `Transaction` di project kita **adalah struct**.

### 1.6 Pointer (`*` dan `&`) — sering bikin bingung pemula
```go
func naikkanUmur(u *int) { // *int artinya "pointer ke int", bukan int biasa
	*u = *u + 1 // ubah nilai yang DITUNJUK pointer
}

func main() {
	umur := 20
	naikkanUmur(&umur) // & artinya "ambil alamat memori variabel ini"
	fmt.Println(umur)  // hasilnya 21, bukan 20!
}
```
Analoginya: `umur` itu **rumah**, `&umur` itu **alamat rumahnya**, `*` itu **"buka pintu rumah di alamat ini"**. Kalau kirim pointer, function bisa **mengubah** data aslinya. Kalau kirim value biasa, function cuma dapat **salinan**.

Kenapa penting? Karena di GORM nanti kita akan sering nulis `&user`, `&trx` — supaya GORM bisa langsung mengisi struct kita dengan data dari database.

### 1.7 Slice (array dinamis) & Map
```go
kategori := []string{"makan", "transport", "uang saku"} // slice
kategori = append(kategori, "beasiswa")                 // nambah elemen

harga := map[string]int{"makan": 25000, "bensin": 15000} // map = key-value, mirip dict/object
fmt.Println(harga["makan"])
```
`[]models.Transaction` yang nanti muncul di kode = **slice berisi banyak struct Transaction**, ini akan sering kamu lihat sebagai hasil query "GET semua data".

### 1.8 Package & Import project sendiri
Satu hal yang beda dari bahasa lain: di Go, kita **import folder kita sendiri** pakai path lengkap `module_name/nama_folder`. Ini nanti kepakai pas kita bikin folder `models`, `handlers`, dll — semua saling import pakai nama module di `go.mod`.

**Cukup segini dulu dasar Go-nya.** Sisanya (interface, goroutine, channel) akan kamu temukan lebih dalam, tapi untuk REST API sederhana ini belum kepakai — kita fokus praktik dulu.

---

## BAB 2 — Server HTTP Pertama dengan Gin

Sekarang kita mulai project sungguhan, terpisah dari folder latihan tadi.

```bash
mkdir -p ~/expense-tracker-api
cd ~/expense-tracker-api
go mod init expense-tracker-api
go get github.com/gin-gonic/gin
```
> `go get` mengunduh library dan mencatatnya otomatis di `go.mod` + `go.sum` (file checksum keamanan).

Buat `main.go`:
```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default() // bikin router Gin dengan logger & recovery middleware bawaan

	r.GET("/ping", func(c *gin.Context) { // c *gin.Context = objek yang bawa data request & tempat nulis response
		c.JSON(200, gin.H{"message": "pong"}) // gin.H itu shortcut untuk map[string]interface{}
	})

	r.Run(":8080") // jalankan server di port 8080 (blocking, program stuck di sini selama server hidup)
}
```
Jalankan:
```bash
go run main.go
```
Buka terminal **baru** (biarkan server tetap jalan), test:
```bash
curl http://localhost:8080/ping
# {"message":"pong"}
```
🎉 Server pertamamu jalan. Konsep penting: **route** = pasangan (method HTTP + path) → handler function.

---

## BAB 3 — Koneksi ke SQLite pakai GORM

```bash
go get gorm.io/gorm
go get github.com/glebarez/sqlite   # driver sqlite pure-Go, tanpa perlu compiler C
```

GORM adalah **ORM** (Object-Relational Mapping) — kamu nulis struct Go, GORM yang bikinkan tabel SQL-nya, tanpa kamu nulis `CREATE TABLE` manual.

Buat `main.go` versi baru (kita masih single-file dulu biar simpel, nanti dirapikan di BAB 8):
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Transaction = model kita. Tag `gorm:"..."` mengatur perilaku kolom di database.
type Transaction struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

var db *gorm.DB // variabel global sederhana untuk koneksi database (nanti kita rapikan)

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("expense.db"), &gorm.Config{})
	if err != nil {
		panic("gagal konek database: " + err.Error()) // panic = hentikan program paksa (dipakai untuk error fatal)
	}

	db.AutoMigrate(&Transaction{}) // GORM baca struct, bikin tabel "transactions" kalau belum ada

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")
}
```
Jalankan `go run main.go`, lalu di terminal lain cek:
```bash
sqlite3 expense.db
.tables
.schema transactions
.quit
```
Kamu akan lihat tabel `transactions` sudah otomatis terbentuk sesuai struct Go tadi — inilah kekuatan AutoMigrate.

---

## BAB 4 — CREATE: Insert Data Pertama

Tambahkan endpoint POST. Update `main.go`, tambahkan sebelum `r.Run(":8080")`:
```go
	r.POST("/transactions", func(c *gin.Context) {
		var input Transaction
		if err := c.ShouldBindJSON(&input); err != nil { // baca body JSON, isi ke struct input
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Create(&input) // INSERT INTO transactions ...
		c.JSON(201, input) // 201 = Created
	})
```
Test:
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"category": "makan", "amount": 25000}'
```
**Sekarang tracking database-nya** (ini kebiasaan yang harus kamu bangun tiap develop fitur baru):
```bash
sqlite3 expense.db "SELECT * FROM transactions;"
```
Harus muncul 1 baris data yang barusan kamu insert. Kalau kamu mau lihat lebih rapi:
```bash
sqlite3 -header -column expense.db "SELECT * FROM transactions;"
```
Atau buka `sqlitebrowser expense.db` kalau mau lihat visual.

> **Kebiasaan yang bagus**: setiap habis test endpoint lewat curl/Postman, langsung cek database-nya. Ini cara paling cepat mendeteksi bug (misal field kosong, tipe salah, dsb) sebelum lanjut ke fitur berikutnya.

---

## BAB 5 — READ: Ambil Semua & Ambil Satu

```go
	r.GET("/transactions", func(c *gin.Context) {
		var transactions []Transaction // slice kosong, nanti diisi GORM
		db.Find(&transactions)          // SELECT * FROM transactions
		c.JSON(200, transactions)
	})

	r.GET("/transactions/:id", func(c *gin.Context) {
		id := c.Param("id") // ambil ":id" dari URL
		var trx Transaction
		result := db.First(&trx, id) // SELECT * FROM transactions WHERE id = ? LIMIT 1
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "data tidak ditemukan"})
			return
		}
		c.JSON(200, trx)
	})
```
Test:
```bash
curl http://localhost:8080/transactions
curl http://localhost:8080/transactions/1
```

---

## BAB 6 — UPDATE & DELETE

```go
	r.PUT("/transactions/:id", func(c *gin.Context) {
		id := c.Param("id")
		var trx Transaction
		if err := db.First(&trx, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "data tidak ditemukan"})
			return
		}

		var input Transaction
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		trx.Category = input.Category
		trx.Amount = input.Amount
		db.Save(&trx) // UPDATE transactions SET ... WHERE id = ?
		c.JSON(200, trx)
	})

	r.DELETE("/transactions/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Delete(&Transaction{}, id) // DELETE FROM transactions WHERE id = ?
		c.JSON(200, gin.H{"message": "berhasil dihapus"})
	})
```
Test lalu **cek lagi pakai sqlite3** setiap habis update/delete — pastikan datanya benar berubah/hilang di database, jangan cuma percaya response JSON-nya saja.

Sampai sini kamu **sudah punya CRUD lengkap** dalam satu file `main.go` (~70 baris). Coba jalankan ulang semua endpoint dari awal (POST → GET → PUT → GET → DELETE → GET) sambil selalu cross-check dengan `sqlite3 expense.db`.

---

## BAB 7 — Query Aggregate (Summary)

```go
	r.GET("/transactions/summary", func(c *gin.Context) {
		var total float64
		db.Model(&Transaction{}).Select("COALESCE(SUM(amount), 0)").Scan(&total)
		c.JSON(200, gin.H{"total": total})
	})
```
> ⚠️ Taruh route `/transactions/summary` **sebelum** `/transactions/:id` kalau kamu urutkan berbeda — Gin mencocokkan route dari atas ke bawah, kalau `:id` didaftarkan duluan, request ke `/summary` akan dikira `id="summary"`.

Cek langsung pakai SQL manual buat bandingin hasilnya benar apa nggak:
```bash
sqlite3 expense.db "SELECT SUM(amount) FROM transactions;"
```

---

## BAB 8 — Merapikan Struktur Folder

Satu file `main.go` isi semuanya itu tidak scalable. Sekarang kita pecah jadi folder rapi — **ini persis struktur di project final yang sudah aku kirim sebelumnya** (`expense-tracker-api.zip`):

```
expense-tracker-api/
├── main.go
├── config/config.go
└── internal/
    ├── database/database.go
    ├── models/{user,transaction}.go
    ├── utils/{jwt,response}.go
    ├── middleware/auth_middleware.go
    ├── handlers/{auth,transaction}_handler.go
    └── routes/routes.go
```

Kenapa dipisah begini?
- **models** → definisi struct/tabel, tidak ada logic
- **database** → cuma urusan koneksi & migrasi
- **handlers** → logic tiap endpoint (yang tadi kamu tulis inline di `main.go`)
- **middleware** → kode yang jalan **sebelum** handler (misal cek token)
- **routes** → daftar "path → handler mana", satu tempat biar gampang dilihat semua endpoint
- **utils** → helper kecil yang dipakai di banyak tempat (JWT, format response)
- **config** → baca `.env`

Silakan **pindahkan** kode CRUD yang barusan kamu tulis manual di BAB 4–7 ke struktur ini satu per satu, sambil bandingkan dengan file yang sudah aku buat di zip sebelumnya. Ini latihan refactoring paling bagus buat pemula — kamu paham *kenapa* dipisah karena kamu baru saja ngerasain versi satu-file-nya.

---

## BAB 9 — Autentikasi: Konsep JWT Dulu (Sebelum Masuk Google OAuth)

Sebelum masuk OAuth (yang agak kompleks karena melibatkan pihak ketiga/Google), pahami dulu **JWT (JSON Web Token)** dengan versi paling sederhana: login pakai email biasa, generate token.

```bash
go get github.com/golang-jwt/jwt/v5
```

Konsep JWT: server bikin "tiket" berisi data user (misal `user_id`) + tanda tangan digital pakai *secret key*. Tiket ini dikirim ke client, client kirim balik tiket ini di setiap request lewat header `Authorization: Bearer <tiket>`. Server tinggal **verifikasi tanda tangannya** tanpa perlu simpan session di database — makanya disebut **stateless**.

Contoh generate & verifikasi (baca-baca dulu, versi lengkapnya sudah ada di `internal/utils/jwt.go` pada project final):
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"user_id": 1,
	"exp":     time.Now().Add(72 * time.Hour).Unix(),
})
signed, _ := token.SignedString([]byte("secret-rahasia"))
fmt.Println(signed) // ini yang dikirim ke client
```
Middleware nanti tugasnya: baca header `Authorization`, verifikasi token, kalau valid → lanjut ke handler; kalau tidak → return 401.

---

## BAB 10 — Google OAuth (Menggabungkan Semua)

Sekarang kamu sudah paham fondasinya: Gin, GORM, CRUD, JWT. Google OAuth cuma **satu cara alternatif untuk mendapatkan identitas user** (menggantikan form login email/password manual), setelah itu **alurnya sama**: generate JWT, pakai JWT untuk proteksi endpoint transaksi.

Alurnya (baca lagi biar nempel):
1. `GET /auth/google/login` → redirect user ke Google
2. User login/setuju di Google
3. Google redirect balik ke `GET /auth/google/callback?code=...`
4. Backend tukar `code` → dapat profil user (email, nama)
5. Backend cari/buat user di SQLite berdasarkan `google_id`
6. Backend generate JWT sendiri (bukan token Google), lempar ke frontend

Karena kode lengkapnya cukup panjang (OAuth config, HTTP call ke Google, dsb), pakai file `internal/handlers/auth_handler.go` dari project final yang sudah aku kirim — jangan ditulis ulang dari nol, tapi **baca baris per baris**, cocokkan dengan penjelasan alur di atas. Setup kredensial Google-nya ada di README project (`expense-tracker-api.zip` → `README.md` bagian 2a).

---

## BAB 11 — Kebiasaan Tracking Database Selama Development

Rangkuman workflow yang disarankan tiap kali develop fitur baru:

```bash
# 1. jalankan server
go run main.go

# 2. di terminal lain, test endpoint
curl -X POST http://localhost:8080/api/transactions -H "Authorization: Bearer $TOKEN" -d '...'

# 3. cek langsung ke database, jangan cuma percaya response JSON
sqlite3 -header -column expense.db "SELECT * FROM transactions ORDER BY id DESC LIMIT 5;"

# 4. kalau mau reset total waktu development (misal skema berubah)
rm expense.db
go run main.go   # AutoMigrate bikin ulang dari nol
```
Tips tambahan: install `litecli` (`pip install litecli`) kalau mau autocomplete SQL yang lebih enak dari `sqlite3` biasa — opsional, tidak wajib.

---

## Urutan Belajar yang Disarankan

1. ✅ Selesaikan BAB 0–7 dengan **mengetik ulang sendiri**, bukan copy-paste, sampai CRUD single-file jalan
2. ✅ Refactor ke struktur folder (BAB 8) sambil bandingkan ke project final
3. ✅ Pahami JWT sederhana (BAB 9) sebelum sentuh OAuth
4. ✅ Baru integrasikan Google OAuth (BAB 10) pakai kode dari project final
5. ✅ Kerjakan "Latihan Tambahan" di README project final (7 latihan lanjutan) untuk memperdalam

Kalau ada bagian yang error atau bingung di tengah jalan, kirim pesan error-nya ke aku — jangan loncat ke bab berikutnya sebelum paham bab sekarang.

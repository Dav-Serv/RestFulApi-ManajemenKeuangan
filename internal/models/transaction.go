package models

import "time"

type TransactionType string

const (
	TypeIncome	TransactionType = "income" // pemasukan
	TypeExpense	TransactionType = "expense" // pengeluaran
)

// Transaction adalah satu catatan pemasukan atau pengeluaran milik seorang user.
type Transaction struct {
	ID			uint				`gorm:"primaryKey" json:"id"`
	UserID		uint				`gorm:"index;not null" json:"user_id"`
	Type		TransactionType		`gorm:"type:varchar(10);not nul; index" json:"type"`
	Category	string				`gorm:"not null;index" json:"category"` // contoh: makan, transport, uang saku, beasiswa
	Amount		float64				`gorm:"not null" json:"amount"`
	Description	string				`json:"description"`
	Date		time.Time			`json:"date"` // tanggal transaksi (bisa beda dari created_at)
	CreatedAt	time.Time			`json:"created_at"`
	UpdatedAt	time.Time			`json:"updated_at"`
}

// TransactionInput dipakai untuk binding + validasi body request create/update.
type TransactionInput struct {
	Type		TransactionType		`json:"type" binding:"required,oneof=income expense"`
	Category	string				`json:"category" binding:"required,min=2,max=50"`
	Amount		float64				`json:"amount" binding:"required,gt=0"`
	Description	string				`json:"description" binding:"max=255"`
	Date		*time.Time			`json:"date"` // opsional, default waktu sekarang jika kosong
}

// SummaryResponse dipakai untuk endpoint ringkasan keuangan.
type SummaryResponse struct {
	TotalIncome		float64		`json:"total_income"`
	TotalExpense	float64		`json:"total_expense"`
	Balance			float64		`json:"balance"`
}
package models

import "time"

// User merepresentasikan mahasiswa yang login lewat Google OAuth.
type User struct {
	ID			uint		`gorm:"primaryKey" json:"id"`
	GoogleID	string		`gorm:"uniqueIndex;not null" json:"-"`
	Email		string		`gorm:"uniqueIndex;not null" json:"email"`
	Name		string		`json:"name"`
	Avatar		string		`json:"avatar"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`

	// Relasi 1-to-many ke Transaction, tidak selalu di-load (pakai Preload jika perlu)
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"-"`
}
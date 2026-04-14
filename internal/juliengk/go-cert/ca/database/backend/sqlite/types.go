package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Certificate struct {
	Model
	StatusFlag        string `gorm:"column:status_flag"`
	ExpirationDate    string `gorm:"column:expiration_date"`
	RevocationDate    string `gorm:"column:revocation_date"`
	RevocationReason  int    `gorm:"column:revocation_reason"`
	SerialNumber      int    `gorm:"unique;column:serial_number"`
	Filename          string
	DistinguishedName string `gorm:"column:distinguished_name"`
}

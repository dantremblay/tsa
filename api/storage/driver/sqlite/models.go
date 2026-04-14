package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type ServerConfig struct {
	Model

	Key   string
	Value string
}

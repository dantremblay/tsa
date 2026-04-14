package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Server struct {
	Model

	Name        string `gorm:"unique;"`
	TSAURL      string
	Description string
}

type Session struct {
	Model

	Server   Server `gorm:"unique;"`
	ServerID uint
	Active   bool
	Token    string
}

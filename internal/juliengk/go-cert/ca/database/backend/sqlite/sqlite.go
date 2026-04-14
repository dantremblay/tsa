package sqlite

import (
	"path"

	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/ca/database/backend"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	database.RegisterBackend("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(config string) (backend.Backender, error) {
	file := path.Join(config, "index.db")

	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Certificate{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	sqlDB, _ := c.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

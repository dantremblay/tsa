package sqlite

import (
	"path"

	"github.com/kassisol/tsa/cli/storage"
	"github.com/kassisol/tsa/cli/storage/driver"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	storage.RegisterDriver("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(config string) (driver.Storager, error) {
	dbFilePath := path.Join(config, "data.db")

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Server{}, &Session{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	sqlDB, _ := c.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

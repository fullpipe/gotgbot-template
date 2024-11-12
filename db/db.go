package db

import (
	"bm/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg config.Config) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Error,           // Log level
			IgnoreRecordNotFoundError: !cfg.Debug,             // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      !cfg.Debug,             // Don't include params in the SQL log
			Colorful:                  false,                  // Disable color
		},
	)

	if cfg.Debug {
		newLogger.LogMode(logger.Info)
	}

	d := postgres.Open(cfg.DSN)

	return gorm.Open(d, &gorm.Config{
		Logger: newLogger,
	})
}

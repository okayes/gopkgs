package sqlserver

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewSQLServerDatabase(dbConn string) *Database {
	db, err := gorm.Open(sqlserver.Open(dbConn), &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "", log.LstdFlags),
			logger.Config{
				SlowThreshold: 100 * time.Millisecond,
				LogLevel:      logger.Warn,
				Colorful:      false,
			},
		),
	})
	if err != nil {
		log.Panicf("NewSQLServerDatabase error: %v\n", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("db.DB error: %v\n", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{
		DB: db,
	}
}

func IsDBErrRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

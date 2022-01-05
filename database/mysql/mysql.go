package mysql

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewMySQLDatabase(dbConn string) *Database {
	db, err := gorm.Open(mysql.Open(dbConn), &gorm.Config{
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
		log.Panicf("NewMySQLDatabase error: %v\n", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("db.DB error: %v\n", err)
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{
		DB: db,
	}
}

func IsDBErrRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

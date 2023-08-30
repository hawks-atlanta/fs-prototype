package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func Default() (db *gorm.DB, err error) {
	const dsn = "host=127.0.0.1 user=sulcud password=sulcud dbname=sulcud port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

/*
TestBaseModel basic unit test to reduce the coverage footprint
*/
func testBaseModel_BeforeSafe(t *testing.T, db *gorm.DB) {
	db.AutoMigrate(&Model{})
	t.Run("Null UUID", func(t *testing.T) {
		assert.Nil(t, db.Create(&Model{}).Error)
	})
	t.Run("Set UUID", func(t *testing.T) {
		assert.Nil(t, db.Create(&Model{UUID: uuid.New()}).Error)
	})
}

func TestBaseModel_BeforeSafe(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		assertions := assert.New(t)
		db, err := database.Default()
		assertions.Nil(err)
		conn, _ := db.DB()
		defer conn.Close()
		testBaseModel_BeforeSafe(t, db)
	})
}

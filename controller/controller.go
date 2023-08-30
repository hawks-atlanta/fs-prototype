package controller

import (
	"github.com/hawks-atlanta/fs-prototype/database"
	"github.com/hawks-atlanta/fs-prototype/models"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func (c *Controller) Close() (err error) {
	db, err := c.DB.DB()
	if err == nil {
		err = db.Close()
	}
	return err
}

func New(db *gorm.DB) (c *Controller, err error) {
	err = db.AutoMigrate(
		&models.Archive{}, &models.File{}, &models.SharedFile{},
	)
	c = &Controller{db}
	return c, err
}

func Default() (c *Controller, err error) {
	db, err := database.Default()
	if err == nil {
		c, err = New(db)
	}
	return c, err
}

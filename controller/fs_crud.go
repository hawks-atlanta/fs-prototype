package controller

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateFile struct {
	Filename        string     `json:"filename"`
	OwnerUUID       uuid.UUID  `json:"ownerUUID"`
	Hash            string     `json:"hash,omitempty"`
	ParentDirectory *uuid.UUID `json:"parentDirectory,omitempty"`
	Size            uint64     `json:"size,omitempty"`
}

var checkFilename = regexp.MustCompile(`(?m)^\w+.+`)

func (cf *CreateFile) Check() (err error) {
	if !checkFilename.MatchString(cf.Filename) {
		err = fmt.Errorf("invalid file name provided, it should start with a number or character")
		return err
	}
	if cf.OwnerUUID == uuid.Nil {
		err = fmt.Errorf("no owner UUID provided")
		return err
	}
	if cf.Hash != "" && cf.Size == 0 {
		err = fmt.Errorf("no empty files allowed")
	}
	return err
}

// Creates a new file in the filesystem index
func (c *Controller) CreateFile(cf *CreateFile) (file models.File, err error) {
	err = cf.Check()
	if err != nil {
		err = fmt.Errorf("invalid file creation request: %w", err)
		return file, err
	}

	// Make sure current user is owner of the directory
	if cf.ParentDirectory != nil && *cf.ParentDirectory != uuid.Nil {
		var parentDirectory models.File
		err = c.DB.
			Where("uuid = ? AND owner_uuid = ?", *cf.ParentDirectory, cf.OwnerUUID).
			First(&parentDirectory).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = fmt.Errorf("user doesn't own directory: %w", err)
			} else {
				err = fmt.Errorf("something went wrong while checking ownership of parent directory: %w", err)
			}
			return file, err
		}
	}

	if cf.Size != 0 { // Create file
		err = c.DB.Transaction(func(tx *gorm.DB) (err error) {
			var archive models.Archive = models.Archive{
				Hash: cf.Hash,
				Size: cf.Size,
			}
			err = tx.
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(&archive).
				Error
			if err != nil {
				return err
			}
			archive = models.Archive{}
			err = tx.
				Where("hash = ? AND size = ?", cf.Hash, cf.Size).
				First(&archive).
				Error
			if err != nil {
				return err
			}
			file = models.File{
				OwnerUUID:   cf.OwnerUUID,
				ArchiveUUID: &archive.UUID,
				ParentUUID:  cf.ParentDirectory,
				Name:        cf.Filename,
			}
			err = tx.
				Create(&file).
				Error
			return err
		})
	} else { // Create directory
		err = c.DB.Transaction(func(tx *gorm.DB) error {
			file = models.File{
				OwnerUUID:  cf.OwnerUUID,
				ParentUUID: cf.ParentDirectory,
				Name:       cf.Filename,
			}
			err = tx.
				Create(&file).
				Error
			return err
		})
	}
	if err != nil {
		err = fmt.Errorf("failed to insert file: %w", err)
	}
	return file, err
}

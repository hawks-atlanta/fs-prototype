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

type DeleteFile struct {
	OwnerUUID uuid.UUID `json:"ownerUUID"`
	FileUUID  uuid.UUID `json:"fileUUID"`
}

func (df *DeleteFile) Check() (err error) {
	if df.OwnerUUID == uuid.Nil {
		err = fmt.Errorf("no owner UUID provided")
		return err
	}
	if df.FileUUID == uuid.Nil {
		err = fmt.Errorf("no file UUID provided")
	}
	return err
}

// TODO: List directory

type QueryFile struct {
	UserUUID uuid.UUID `json:"userUUID"`
	FileUUID uuid.UUID `json:"fileUUID"`
}

// TODO: Test this
func (qf *QueryFile) Check() (err error) {
	if qf.UserUUID == uuid.Nil {
		err = fmt.Errorf("no user UUID provided")
		return err
	}
	if qf.FileUUID == uuid.Nil {
		err = fmt.Errorf("no file UUID provided")
	}
	return err
}

// Intended to only be used by the Gateway
// The server checks if the user owns the file.
// If not the server tries to determine the access to the file by shared files with this account
func (c *Controller) QueryFile(qf *QueryFile) (archive models.Archive, err error) {
	err = qf.Check()
	if err != nil {
		err = fmt.Errorf("invalid query file request: %w", err)
		return archive, err
	}
	var crf = CanReadFile{
		UserUUID: qf.UserUUID,
		FileUUID: qf.FileUUID,
	}
	err = c.CanReadFile(&crf)
	if err != nil {
		return archive, err
	}
	err = c.DB.
		Raw(`
		SELECT archives.* 
		FROM archives, files 
		WHERE
			archives.uuid = files.archive_uuid 
			AND files.uuid = ?
		LIMIT 1`, qf.FileUUID).
		Scan(&archive).
		Error
	return archive, err
}

// Deletes file from the index
func (c *Controller) DeleteFile(df *DeleteFile) (err error) {
	err = df.Check()
	if err != nil {
		err = fmt.Errorf("invalid file deletion request: %w", err)
		return err
	}

	err = c.DB.
		Where("uuid = ? AND owner_uuid = ?", df.FileUUID, df.OwnerUUID).
		Delete(&models.File{}).
		Error
	if err != nil {
		err = fmt.Errorf("failed to delete file: %w", err)
	}
	return err
}

// TODO: Move file

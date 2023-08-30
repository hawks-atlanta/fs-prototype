package controller

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/models"
	"gorm.io/gorm"
)

type CanReadFile struct {
	UserUUID uuid.UUID `json:"userUUID"`
	FileUUID uuid.UUID `json:"fileUUID"`
}

func (crf *CanReadFile) Check() (err error) {
	if crf.UserUUID == uuid.Nil {
		err = fmt.Errorf("no user UUID provided")
		return err
	}
	if crf.FileUUID == uuid.Nil {
		err = fmt.Errorf("no file UUID provided")
	}
	return err
}

// Can read file is inteded to be used internally by other operations of the metadata
// Will check if user owns the file
// Or iif user has at least access by share directly or indirectly
func (c *Controller) CanReadFile(crf *CanReadFile) (err error) {
	err = c.DB.Transaction(func(tx *gorm.DB) error {
		var file models.File
		// Check if user is owner
		err := tx.
			Where("uuid = ? AND owner_uuid = ?", crf.FileUUID, crf.UserUUID).
			First(&file).
			Error
		if err == nil {
			return err
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("failed to query file information: %w", err)
			return err
		}
		var sf models.SharedFile
		// Check direct access
		err = c.DB.
			Where("file_uuid = ? AND user_uuid = ?", crf.FileUUID, crf.UserUUID).
			First(&sf).
			Error
		if err == nil {
			return err
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("failed to query shared files information: %w", err)
			return err
		}
		// Check nested
		var found struct {
			Found bool `gorm:"column:found"`
		}
		fmt.Println(crf.FileUUID, crf.UserUUID)
		err = c.DB.Raw(
			`WITH RECURSIVE file_hierarchy AS (
				-- Base case: start with the initial file UUID
				SELECT uuid, parent_uuid
				FROM files
				WHERE uuid = $1
			
				UNION ALL
			
				-- Recursive case: get the parent file of the current file
				SELECT f.uuid, f.parent_uuid
				FROM files f
				JOIN file_hierarchy fh ON f.uuid = fh.parent_uuid
			)
			
			-- Check if any of the files in the hierarchy are shared with the given user
			SELECT CASE 
					   WHEN EXISTS (
						   SELECT 1 
						   FROM shared_files sf
						   JOIN file_hierarchy fh ON sf.file_uuid = fh.uuid
						   WHERE sf.user_uuid = $2
					   ) THEN TRUE
					   ELSE FALSE
				   END AS found;
			`, crf.FileUUID, crf.UserUUID).
			Scan(&found).
			Error
		if err != nil {
			return err
		}
		if !found.Found {
			err = fmt.Errorf("permission denied")
		}
		return err
	})
	return err
}

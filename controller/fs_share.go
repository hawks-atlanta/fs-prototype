package controller

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hawks-atlanta/fs-prototype/models"
	"gorm.io/gorm"
)

type ShareWithMe struct {
	UserUUID uuid.UUID `json:"userUUID"`
}

func (swm *ShareWithMe) Check() (err error) {
	if swm.UserUUID == uuid.Nil {
		err = fmt.Errorf("user uuid not provided")
	}
	return err
}

// Used to list all the files shared with the current user
func (c *Controller) ShareWithMe(swm *ShareWithMe) (shared []models.SharedFile, err error) {
	err = swm.Check()
	if err != nil {
		err = fmt.Errorf("invalid shared with me request: %w", err)
		return shared, err
	}
	err = c.DB.
		Where("user_uuid = ?", swm.UserUUID).
		Find(&shared).
		Error
	if err != nil {
		err = fmt.Errorf("failed to obtain share files for user: %w", err)
	}
	return shared, err
}

type ShareWithWho struct {
	OwnerUUID uuid.UUID `json:"ownerUUID"`
	FileUUID  uuid.UUID `json:"fileUUID"`
}

func (sww *ShareWithWho) Check() (err error) {
	if sww.OwnerUUID == uuid.Nil {
		err = fmt.Errorf("no owner UUID provided")
		return err
	}
	if sww.FileUUID == uuid.Nil {
		err = fmt.Errorf("no file UUID provided")
	}
	return err
}

// Used to query users that have access to a file
func (c *Controller) ShareWithWho(sww *ShareWithWho) (shared []models.SharedFile, err error) {
	err = sww.Check()
	if err != nil {
		err = fmt.Errorf("invalid shared with who request: %w", err)
		return shared, err
	}
	err = c.DB.Transaction(func(tx *gorm.DB) error {
		var file models.File
		err := tx.
			Where("uuid = ? AND owner_uuid = ?", sww.FileUUID, sww.OwnerUUID).
			First(&file).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = fmt.Errorf("user doesn't have permissions over file: %w", err)
			} else {
				err = fmt.Errorf("failed to query file access: %w", err)
			}
			return err
		}
		err = tx.
			Where("file_uuid = ?", file.UUID).
			Find(&shared).
			Error
		if err != nil {
			err = fmt.Errorf("failed to query files: %w", err)
		}
		return err
	})
	return shared, err
}

type ShareRequest struct {
	OwnerUUID      uuid.UUID `json:"ownerUUID"`
	FileUUID       uuid.UUID `json:"fileUUID"`
	TargetUserUUID uuid.UUID `json:"targetUserUUID"`
}

func (sr *ShareRequest) Check() (err error) {
	if sr.OwnerUUID == uuid.Nil {
		err = fmt.Errorf("no owner UUID provided")
		return err
	}
	if sr.FileUUID == uuid.Nil {
		err = fmt.Errorf("no file UUID provided")
		return err
	}
	if sr.TargetUserUUID == uuid.Nil {
		err = fmt.Errorf("no target user UUID provided")
	}
	return err
}

// Use to share a file other users in the system
// Intended to be called after obtaining the UUID of the account thanks to the authentication service
func (c *Controller) ShareFile(sr *ShareRequest) (err error) {
	err = sr.Check()
	if err != nil {
		err = fmt.Errorf("invalid share request request: %w", err)
		return err
	}
	err = c.DB.Transaction(func(tx *gorm.DB) error {
		var file models.File
		err := tx.
			Where("uuid = ? AND owner_uuid = ?", sr.FileUUID, sr.OwnerUUID).
			First(&file).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = fmt.Errorf("permission denied: %w", err)
			} else {
				err = fmt.Errorf("failed to query file: %w", err)
			}
			return err
		}
		err = tx.
			Create(&models.SharedFile{
				FileUUID: file.UUID,
				UserUUID: sr.TargetUserUUID,
			}).
			Error
		if err != nil {
			err = fmt.Errorf("failed to create shared entry: %w", err)
		}
		return err
	})
	return err
}

// Work almost the same as the ShareFile but intended to remove files
func (c *Controller) UnshareFile(sr *ShareRequest) (err error) {
	err = sr.Check()
	if err != nil {
		err = fmt.Errorf("invalid share request request: %w", err)
		return err
	}
	err = c.DB.Transaction(func(tx *gorm.DB) error {
		var file models.File
		err := tx.
			Where("uuid = ? AND owner_uuid = ?", sr.FileUUID, sr.OwnerUUID).
			First(&file).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = fmt.Errorf("permission denied: %w", err)
			} else {
				err = fmt.Errorf("failed to query file: %w", err)
			}
			return err
		}
		err = tx.
			Where("file_uuid = ? AND user_uuid = ?", file.UUID, sr.TargetUserUUID).
			Delete(&models.SharedFile{}).
			Error
		if err != nil {
			err = fmt.Errorf("failed to create shared entry: %w", err)
		}
		return err
	})
	return err
}

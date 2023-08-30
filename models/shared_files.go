package models

import "github.com/google/uuid"

type SharedFile struct {
	Model
	UserUUID uuid.UUID `json:"userUUID" gorm:"uniqueIndex:idx_unique_shared_file;not null;"`
	File     *File     `json:"file,omitempty" gorm:"foreignKey:FileUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FileUUID uuid.UUID `json:"fileUUID,omitempty" gorm:"uniqueIndex:idx_unique_shared_file;not null;"`
}

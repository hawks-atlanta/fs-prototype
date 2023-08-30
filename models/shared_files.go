package models

import "github.com/google/uuid"

type SharedFile struct {
	Model
	OwnerUUID uuid.UUID `json:"ownerUUID" gorm:"uniqueIndex:idx_unique_shared_file;not null;"`
	File      *File     `json:"parent,omitempty" gorm:"foreignKey:FileUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FileUUID  uuid.UUID `json:"parentUUID,omitempty" gorm:"uniqueIndex:idx_unique_shared_file;not null;"`
}

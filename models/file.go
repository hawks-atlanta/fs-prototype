package models

import "github.com/google/uuid"

type File struct {
	Model
	OwnerUUID   uuid.UUID  `json:"ownerUUID" gorm:"uniqueIndex:idx_unique_file;not null;"`
	Parent      *File      `json:"parent,omitempty" gorm:"foreignKey:ParentUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentUUID  *uuid.UUID `json:"parentUUID,omitempty" gorm:"uniqueIndex:idx_unique_file;"`
	Archive     *Archive   `json:"archive,omitempty" gorm:"foreignKey:ArchiveUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ArchiveUUID *uuid.UUID `json:"archiveUUID,omitempty"`
	Name        string     `json:"name" gorm:"uniqueIndex:idx_unique_file;not null;"`
}

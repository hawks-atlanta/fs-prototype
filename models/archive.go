package models

type Archive struct {
	Model
	Hash    string `json:"hash" gorm:"uniqueIndex:idx_unique_archive;not null;"`
	Size    uint64 `json:"size" gorm:"uniqueIndex:idx_unique_archive;not null;"`
	IsReady bool   `json:"isReady" gorm:"not null;"`
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uint32         `gorm:"primary_key;auto_increment" json:"-"`
	UUID      uuid.UUID      `gorm:"uniqueIndex; type:char(36);" json:"id"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

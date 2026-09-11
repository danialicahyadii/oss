package model

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID          uint64         `gorm:"column:id;primaryKey"`
	UUID        string         `gorm:"column:uuid"`
	BussinessID string         `gorm:"column:business_id"`
	SofficeID   uint64         `gorm:"column:soffice_id"`
	VariantID   uint64         `gorm:"column:variant_id"`
	Notes       *string        `gorm:"column:notes"`
	ValueType   string         `gorm:"column:value_type"`
	CurrentRole string         `gorm:"column:current_role"`
	Status      string         `gorm:"column:status"`
	CreatedBy   uint64         `gorm:"column:created_by"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Submission) TableName() string {
	return "sm_submissions"
}

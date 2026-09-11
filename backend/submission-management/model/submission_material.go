package model

import (
	"time"

	"gorm.io/gorm"
)

type SubmissionMaterial struct {
	ID           uint64         `gorm:"column:id;primaryKey"`
	UUID         string         `gorm:"column:uuid"`
	SubmissionID uint64         `gorm:"column:submission_id"`
	MaterialID   uint64         `gorm:"column:material_id"`
	QtyJual      float64        `gorm:"column:qty_jual"`
	SalesUOM     string         `gorm:"column:sales_uom"`
	Value        float64        `gorm:"column:value"`
	CreatedBy    uint64         `gorm:"column:created_by"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (SubmissionMaterial) TableName() string {
	return "sm_submission_materials"
}

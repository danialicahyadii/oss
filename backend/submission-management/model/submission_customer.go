package model

import (
	"time"

	"gorm.io/gorm"
)

type SubmissionCustomer struct {
	ID           uint64         `gorm:"column:id;primaryKey"`
	UUID         string         `gorm:"column:uuid"`
	SubmissionID uint64         `gorm:"column:submission_id"`
	CustomerID   uint64         `gorm:"column:customer_id"`
	CreatedBy    uint64         `gorm:"column:created_by"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (SubmissionCustomer) TableName() string {
	return "sm_submission_customers"
}

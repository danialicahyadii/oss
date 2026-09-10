package model

import "time"

type Role struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	UUID        string     `gorm:"column:uuid"`
	Name        string     `gorm:"column:name"`
	Description *string    `gorm:"column:description"`
	IsActive    bool       `gorm:"column:is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (Role) TableName() string {
	return "um_roles"
}

package model

import "time"

type UserRole struct {
	ID        int        `gorm:"column:id;primaryKey"`
	UserID    uint64     `gorm:"column:user_id"`
	RoleID    uint64     `gorm:"column:role_id"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (UserRole) TableName() string {
	return "um_user_roles"
}

package model

import "time"

type RolePermission struct {
	ID           int       `gorm:"column:id;primaryKey"`
	RoleID       uint64    `gorm:"column:role_id"`
	PermissionID uint64    `gorm:"column:permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (RolePermission) TableName() string {
	return "um_role_permissions"
}

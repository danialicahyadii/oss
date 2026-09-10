package model

import "time"

type User struct {
	ID               uint64     `gorm:"column:id;primaryKey"`
	UUID             string     `gorm:"column:uuid"`
	Username         string     `gorm:"column:username"`
	Name             string     `gorm:"column:name"`
	LastName         string     `gorm:"column:last_name"`
	Email            string     `gorm:"column:email"`
	Password         string     `gorm:"column:password"`
	IsActive         bool       `gorm:"column:is_active"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "um_users"
}

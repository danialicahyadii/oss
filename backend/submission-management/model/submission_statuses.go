package model

import "time"

type SubmissionStatus struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	UUID         string    `gorm:"column:uuid"`
	SubmissionID uint64    `gorm:"column:submission_id"`
	FromStatus   string    `gorm:"column:from_status"`
	ToStatus     string    `gorm:"column:to_status"`
	ActorID      uint64    `gorm:"column:actor_id"`
	ActorRole    string    `gorm:"column:actor_role"`
	Notes        *string   `gorm:"column:notes"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (SubmissionStatus) TableName() string {
	return "sm_submission_statuses"
}

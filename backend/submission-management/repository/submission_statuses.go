package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/submission-management/model"
)

type SubmissionStatusRepository struct {
	db *gorm.DB
}

func NewSubmissionStatusRepository(db *gorm.DB) *SubmissionStatusRepository {
	return &SubmissionStatusRepository{
		db: db,
	}
}

func (r *SubmissionStatusRepository) Create(
	ctx context.Context,
	status *model.SubmissionStatus,
) error {
	return r.db.WithContext(ctx).Create(status).Error
}

func (r *SubmissionStatusRepository) FindBySubmissionID(
	ctx context.Context,
	submissionID uint64,
) ([]model.SubmissionStatus, error) {
	var statuses []model.SubmissionStatus

	err := r.db.WithContext(ctx).
		Where("submission_id = ?", submissionID).
		Order("created_at ASC").
		Find(&statuses).Error

	if err != nil {
		return nil, err
	}

	return statuses, nil
}

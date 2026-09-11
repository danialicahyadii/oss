package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/submission-management/model"
)

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}

func (r *SubmissionRepository) WithTx(tx *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{
		db: tx,
	}
}

func (r *SubmissionRepository) Create(
	ctx context.Context,
	submission *model.Submission,
) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *SubmissionRepository) FindByUUID(
	ctx context.Context,
	uuid string,
) (*model.Submission, error) {
	var submission model.Submission

	err := r.db.WithContext(ctx).
		Where("uuid = ?", uuid).
		First(&submission).Error

	if err != nil {
		return nil, err
	}

	return &submission, nil
}

func (r *SubmissionRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.Submission, error) {
	var submission model.Submission

	err := r.db.WithContext(ctx).
		First(&submission, id).Error

	if err != nil {
		return nil, err
	}

	return &submission, nil
}

func (r *SubmissionRepository) Update(
	ctx context.Context,
	submission *model.Submission,
) error {
	return r.db.WithContext(ctx).Save(submission).Error
}

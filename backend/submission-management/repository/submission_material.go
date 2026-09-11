package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/submission-management/model"
)

type SubmissionMaterialRepository struct {
	db *gorm.DB
}

func NewSubmissionMaterialRepository(db *gorm.DB) *SubmissionMaterialRepository {
	return &SubmissionMaterialRepository{
		db: db,
	}
}

func (r *SubmissionMaterialRepository) Create(
	ctx context.Context,
	material *model.SubmissionMaterial,
) error {
	return r.db.WithContext(ctx).Create(material).Error
}

func (r *SubmissionMaterialRepository) CreateBatch(
	ctx context.Context,
	materials []model.SubmissionMaterial,
) error {
	if len(materials) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Create(&materials).Error
}

func (r *SubmissionMaterialRepository) FindBySubmissionID(
	ctx context.Context,
	submissionID uint64,
) ([]model.SubmissionMaterial, error) {
	var materials []model.SubmissionMaterial

	err := r.db.WithContext(ctx).
		Where("submission_id = ?", submissionID).
		Find(&materials).Error

	if err != nil {
		return nil, err
	}

	return materials, nil
}

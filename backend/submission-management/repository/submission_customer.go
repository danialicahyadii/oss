package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/submission-management/model"
)

type SubmissionCustomerRepository struct {
	db *gorm.DB
}

func NewSubmissionCustomerRepository(db *gorm.DB) *SubmissionCustomerRepository {
	return &SubmissionCustomerRepository{
		db: db,
	}
}

func (r *SubmissionCustomerRepository) Create(
	ctx context.Context,
	customer *model.SubmissionCustomer,
) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *SubmissionCustomerRepository) CreateBatch(
	ctx context.Context,
	customers []model.SubmissionCustomer,
) error {
	if len(customers) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Create(&customers).Error
}

func (r *SubmissionCustomerRepository) FindBySubmissionID(
	ctx context.Context,
	submissionID uint64,
) ([]model.SubmissionCustomer, error) {
	var customers []model.SubmissionCustomer

	err := r.db.WithContext(ctx).
		Where("submission_id = ?", submissionID).
		Find(&customers).Error

	if err != nil {
		return nil, err
	}

	return customers, nil
}

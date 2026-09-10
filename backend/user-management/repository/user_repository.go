package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).
		Where("username = ? AND deleted_at IS NULL", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).
		Where(
			"id = ? AND deleted_at IS NULL",
			id,
		).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) UpdatePassword(
	ctx context.Context,
	userID uint64,
	password string,
) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("password", password).
		Error
}

package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{
		db: db,
	}
}

func (r *UserRoleRepository) FindByUserID(
	ctx context.Context,
	userID uint64,
) ([]model.UserRole, error) {

	var userRoles []model.UserRole

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&userRoles).Error

	if err != nil {
		return nil, err
	}

	return userRoles, nil
}

func (r *UserRoleRepository) Create(
	ctx context.Context,
	userRole *model.UserRole,
) error {

	return r.db.WithContext(ctx).
		Create(userRole).
		Error
}

func (r *UserRoleRepository) Delete(
	ctx context.Context,
	userID uint64,
	roleID uint64,
) error {

	result := r.db.WithContext(ctx).
		Where(
			"user_id = ? AND role_id = ?",
			userID,
			roleID,
		).
		Delete(&model.UserRole{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

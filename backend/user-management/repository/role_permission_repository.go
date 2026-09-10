package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{
		db: db,
	}
}

func (r *RolePermissionRepository) FindByRoleID(
	ctx context.Context,
	roleID uint64,
) ([]model.RolePermission, error) {

	var rolePermissions []model.RolePermission

	err := r.db.WithContext(ctx).
		Where(
			"role_id = ?",
			roleID,
		).
		Find(&rolePermissions).Error

	if err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

func (r *RolePermissionRepository) Create(
	ctx context.Context,
	rolePermission *model.RolePermission,
) error {

	return r.db.WithContext(ctx).
		Create(rolePermission).
		Error
}

func (r *RolePermissionRepository) Delete(
	ctx context.Context,
	roleID uint64,
	permissionID uint64,
) error {

	result := r.db.WithContext(ctx).
		Model(&model.RolePermission{}).
		Where(
			"role_id = ? AND permission_id = ?",
			roleID,
			permissionID,
		).Delete(&model.RolePermission{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

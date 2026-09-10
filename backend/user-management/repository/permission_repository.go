package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) FindAll(
	ctx context.Context,
	search string,
	offset int,
	limit int,
) ([]model.Permission, int64, error) {

	var permissions []model.Permission
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Where("deleted_at IS NULL")

	if search != "" {
		searchValue := "%" + search + "%"

		query = query.Where(
			"permission_code ILIKE ? OR permission_name ILIKE ? OR description ILIKE ?",
			searchValue,
			searchValue,
			searchValue,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *PermissionRepository) FindByUUID(
	ctx context.Context,
	permissionUUID string,
) (*model.Permission, error) {

	var permission model.Permission

	err := r.db.WithContext(ctx).
		Where(
			"uuid = ? AND deleted_at IS NULL",
			permissionUUID,
		).
		First(&permission).Error

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.Permission, error) {

	var permission model.Permission

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&permission).Error

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) FindByCode(
	ctx context.Context,
	permissionCode string,
) (*model.Permission, error) {

	var permission model.Permission

	err := r.db.WithContext(ctx).
		Where(
			"permission_code = ? AND deleted_at IS NULL",
			permissionCode,
		).
		First(&permission).Error

	if err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) Create(
	ctx context.Context,
	permission *model.Permission,
) error {

	return r.db.WithContext(ctx).
		Create(permission).
		Error
}

func (r *PermissionRepository) Update(
	ctx context.Context,
	permissionUUID string,
	updates map[string]interface{},
) error {

	result := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Where(
			"uuid = ? AND deleted_at IS NULL",
			permissionUUID,
		).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PermissionRepository) Delete(
	ctx context.Context,
	permissionUUID string,
) error {

	result := r.db.WithContext(ctx).
		Where(
			"uuid = ? AND deleted_at IS NULL",
			permissionUUID,
		).
		Delete(&model.Permission{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *PermissionRepository) HasPermission(
	ctx context.Context,
	userID uint64,
	permission string,
) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).
		Table("um_user_roles ur").
		Joins("JOIN um_role_permissions rp ON rp.role_id = ur.role_id").
		Joins("JOIN um_permissions p ON p.id = rp.permission_id").
		Where("ur.user_id = ?", userID).
		Where("p.name = ?", permission).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

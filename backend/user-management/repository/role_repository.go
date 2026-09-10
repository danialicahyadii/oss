package repository

import (
	"context"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) FindAll(
	ctx context.Context,
	search string,
	offset int,
	limit int,
) ([]model.Role, int64, error) {

	var roles []model.Role
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.Role{}).
		Where("deleted_at IS NULL")

	if search != "" {
		searchValue := "%" + search + "%"

		query = query.Where(
			"role_code ILIKE ? OR role_name ILIKE ? OR description ILIKE ?",
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
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleRepository) FindByUUID(
	ctx context.Context,
	roleUUID string,
) (*model.Role, error) {

	var role model.Role

	err := r.db.WithContext(ctx).
		Where(
			"uuid = ?",
			roleUUID,
		).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.Role, error) {

	var role model.Role

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) FindByCode(
	ctx context.Context,
	roleCode string,
) (*model.Role, error) {

	var role model.Role

	err := r.db.WithContext(ctx).
		Where(
			"role_code = ? AND deleted_at IS NULL",
			roleCode,
		).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) Create(
	ctx context.Context,
	role *model.Role,
) error {

	return r.db.WithContext(ctx).
		Create(role).
		Error
}

func (r *RoleRepository) Update(
	ctx context.Context,
	roleUUID string,
	updates map[string]interface{},
) error {

	result := r.db.WithContext(ctx).
		Model(&model.Role{}).
		Where(
			"uuid = ? AND deleted_at IS NULL",
			roleUUID,
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

func (r *RoleRepository) Delete(
	ctx context.Context,
	roleUUID string,
) error {

	result := r.db.WithContext(ctx).
		Where(
			"uuid = ? AND deleted_at IS NULL",
			roleUUID,
		).
		Delete(&model.Role{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

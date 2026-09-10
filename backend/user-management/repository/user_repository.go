package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ============================================================
// AUTH
// ============================================================

func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).
		Where(
			"username = ? AND deleted_at IS NULL",
			username,
		).
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
		Where(
			"email = ? AND deleted_at IS NULL",
			email,
		).
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

	return r.db.WithContext(ctx).
		Create(user).
		Error
}

func (r *UserRepository) UpdatePassword(
	ctx context.Context,
	userID uint64,
	password string,
) error {

	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where(
			"id = ? AND deleted_at IS NULL",
			userID,
		).
		Update("password", password).
		Error
}

// ============================================================
// USER CRUD
// ============================================================

func (r *UserRepository) FindAll(
	ctx context.Context,
	search string,
	isActive *bool,
	offset int,
	limit int,
) ([]model.User, int64, error) {

	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("deleted_at IS NULL")

	if search != "" {
		searchValue := "%" + search + "%"

		query = query.Where(
			"username ILIKE ? OR email ILIKE ? OR name ILIKE ? OR last_name ILIKE ?",
			searchValue,
			searchValue,
			searchValue,
			searchValue,
		)
	}

	if isActive != nil {
		query = query.Where(
			"is_active = ?",
			*isActive,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByUUID(
	ctx context.Context,
	userUUID string,
) (*model.User, error) {

	var user model.User

	err := r.db.WithContext(ctx).
		Where(
			"uuid = ?",
			userUUID,
		).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	userUUID string,
	updates map[string]interface{},
) error {

	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where(
			"uuid = ?",
			userUUID,
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

func (r *UserRepository) Delete(
	ctx context.Context,
	userUUID string,
) error {

	result := r.db.WithContext(ctx).
		Where(
			"uuid = ? AND deleted_at IS NULL",
			userUUID,
		).
		Delete(&model.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GenerateUUID hanya helper jika ingin digunakan repository.
func GenerateUUID() string {
	return uuid.New().String()
}

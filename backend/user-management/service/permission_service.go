package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/model"
	"oss.kftd.co.id/v2/user-management/repository"
)

type PermissionService struct {
	permissionRepo *repository.PermissionRepository
}

func NewPermissionService(
	permissionRepo *repository.PermissionRepository,
) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

func (s *PermissionService) GetPermissions(
	ctx context.Context,
	search string,
	page int,
	limit int,
) ([]model.Permission, int64, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	return s.permissionRepo.FindAll(
		ctx,
		strings.TrimSpace(search),
		offset,
		limit,
	)
}

func (s *PermissionService) GetPermission(
	ctx context.Context,
	permissionUUID string,
) (*model.Permission, error) {

	permissionUUID = strings.TrimSpace(permissionUUID)

	if permissionUUID == "" {
		return nil, errors.New("permission uuid is required")
	}

	permission, err := s.permissionRepo.FindByUUID(
		ctx,
		permissionUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}

		return nil, err
	}

	return permission, nil
}

func (s *PermissionService) CreatePermission(
	ctx context.Context,
	req request.CreatePermissionRequest,
) (*model.Permission, error) {

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if name == "" {
		return nil, errors.New("permission name is required")
	}

	var descriptionPtr *string
	if description != "" {
		descriptionPtr = &description
	}

	permission := &model.Permission{
		UUID:        uuid.New().String(),
		Name:        name,
		Description: descriptionPtr,
		IsActive:    true,
	}

	if err := s.permissionRepo.Create(
		ctx,
		permission,
	); err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *PermissionService) UpdatePermission(
	ctx context.Context,
	permissionUUID string,
	req request.UpdatePermissionRequest,
) error {

	permissionUUID = strings.TrimSpace(permissionUUID)

	if permissionUUID == "" {
		return errors.New("permission uuid is required")
	}

	_, err := s.permissionRepo.FindByUUID(
		ctx,
		permissionUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}

		return err
	}

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if name == "" {
		return errors.New("permission name is required")
	}

	var descriptionPtr *string

	if description != "" {
		descriptionPtr = &description
	}

	return s.permissionRepo.Update(
		ctx,
		permissionUUID,
		map[string]interface{}{
			"name":        name,
			"description": descriptionPtr,
			"is_active":   req.IsActive,
		},
	)
}

func (s *PermissionService) DeletePermission(
	ctx context.Context,
	permissionUUID string,
) error {

	permissionUUID = strings.TrimSpace(permissionUUID)

	if permissionUUID == "" {
		return errors.New("permission uuid is required")
	}

	if _, err := s.permissionRepo.FindByUUID(
		ctx,
		permissionUUID,
	); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission not found")
		}

		return err
	}

	if err := s.permissionRepo.Delete(
		ctx,
		permissionUUID,
	); err != nil {
		return err
	}

	return nil
}

func (s *PermissionService) HasPermission(
	ctx context.Context,
	userID uint64,
	permission string,
) (bool, error) {

	return s.permissionRepo.HasPermission(
		ctx,
		userID,
		permission,
	)
}

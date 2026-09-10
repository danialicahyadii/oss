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

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(
	roleRepo *repository.RoleRepository,
) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

func (s *RoleService) GetRoles(
	ctx context.Context,
	search string,
	page int,
	limit int,
) ([]model.Role, int64, error) {

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

	return s.roleRepo.FindAll(
		ctx,
		strings.TrimSpace(search),
		offset,
		limit,
	)
}

func (s *RoleService) GetRole(
	ctx context.Context,
	roleUUID string,
) (*model.Role, error) {

	roleUUID = strings.TrimSpace(roleUUID)

	if roleUUID == "" {
		return nil, errors.New("role uuid is required")
	}

	role, err := s.roleRepo.FindByUUID(
		ctx,
		roleUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}

		return nil, err
	}

	return role, nil
}

func (s *RoleService) CreateRole(
	ctx context.Context,
	req request.CreateRoleRequest,
) (*model.Role, error) {

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if name == "" {
		return nil, errors.New("role name is required")
	}

	var descriptionPtr *string
	if description != "" {
		descriptionPtr = &description
	}

	role := &model.Role{
		UUID:        uuid.New().String(),
		Name:        name,
		Description: descriptionPtr,
		IsActive:    true,
	}

	if err := s.roleRepo.Create(
		ctx,
		role,
	); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) UpdateRole(
	ctx context.Context,
	roleUUID string,
	req request.UpdateRoleRequest,
) error {

	roleUUID = strings.TrimSpace(roleUUID)

	if roleUUID == "" {
		return errors.New("role uuid is required")
	}

	_, err := s.roleRepo.FindByUUID(
		ctx,
		roleUUID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}

		return err
	}

	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)

	if name == "" {
		return errors.New("role name is required")
	}

	var descriptionPtr *string

	if description != "" {
		descriptionPtr = &description
	}

	return s.roleRepo.Update(
		ctx,
		roleUUID,
		map[string]interface{}{
			"name":        name,
			"description": descriptionPtr,
			"is_active":   req.IsActive,
		},
	)
}

func (s *RoleService) DeleteRole(
	ctx context.Context,
	roleUUID string,
) error {

	roleUUID = strings.TrimSpace(roleUUID)

	if roleUUID == "" {
		return errors.New("role uuid is required")
	}

	if _, err := s.roleRepo.FindByUUID(
		ctx,
		roleUUID,
	); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}

		return err
	}

	if err := s.roleRepo.Delete(
		ctx,
		roleUUID,
	); err != nil {
		return err
	}

	return nil
}

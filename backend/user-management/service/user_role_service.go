package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/dto/response"
	"oss.kftd.co.id/v2/user-management/model"
	"oss.kftd.co.id/v2/user-management/repository"
)

type UserRoleService struct {
	userRepo     *repository.UserRepository
	roleRepo     *repository.RoleRepository
	userRoleRepo *repository.UserRoleRepository
}

func NewUserRoleService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	userRoleRepo *repository.UserRoleRepository,
) *UserRoleService {
	return &UserRoleService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (s *UserRoleService) UpdateUserRoles(
	ctx context.Context,
	userUUID string,
	req request.UpdateUserRolesRequest,
) error {

	userUUID = strings.TrimSpace(userUUID)

	if userUUID == "" {
		return errors.New("user uuid is required")
	}

	// Find user
	user, err := s.userRepo.FindByUUID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}

		return err
	}

	// UUID → Role ID
	requestRoleIDs := make(map[uint64]bool)

	for _, roleUUID := range req.RoleUUIDs {

		roleUUID = strings.TrimSpace(roleUUID)

		if roleUUID == "" {
			continue
		}

		role, err := s.roleRepo.FindByUUID(ctx, roleUUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("role not found")
			}

			return err
		}

		requestRoleIDs[role.ID] = true
	}

	// Existing roles
	existingUserRoles, err := s.userRoleRepo.FindByUserID(
		ctx,
		user.ID,
	)
	if err != nil {
		return err
	}

	existingRoleIDs := make(map[uint64]bool)

	for _, userRole := range existingUserRoles {
		existingRoleIDs[userRole.RoleID] = true
	}

	// Add new roles
	for roleID := range requestRoleIDs {

		if existingRoleIDs[roleID] {
			continue
		}

		userRole := &model.UserRole{
			UserID: user.ID,
			RoleID: roleID,
		}

		if err := s.userRoleRepo.Create(
			ctx,
			userRole,
		); err != nil {
			return fmt.Errorf(
				"failed to assign role: role_id=%d user_id=%d: %w",
				user.ID,
				err,
			)
		}
	}

	// Remove unchecked roles
	for roleID := range existingRoleIDs {

		if requestRoleIDs[roleID] {
			continue
		}

		if err := s.userRoleRepo.Delete(
			ctx,
			user.ID,
			roleID,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *UserRoleService) GetUserRoles(
	ctx context.Context,
	userUUID string,
) ([]response.UserRoleResponse, error) {

	userUUID = strings.TrimSpace(userUUID)

	if userUUID == "" {
		return nil, errors.New("user uuid is required")
	}

	user, err := s.userRepo.FindByUUID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	userRoles, err := s.userRoleRepo.FindByUserID(
		ctx,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	result := make([]response.UserRoleResponse, 0)

	for _, userRole := range userRoles {

		role, err := s.roleRepo.FindByID(
			ctx,
			userRole.RoleID,
		)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			return nil, err
		}

		description := ""

		if role.Description != nil {
			description = *role.Description
		}

		result = append(result, response.UserRoleResponse{
			RoleUUID:    role.UUID,
			Name:        role.Name,
			Description: description,
		})
	}

	return result, nil
}

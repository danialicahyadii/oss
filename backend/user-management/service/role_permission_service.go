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

type RolePermissionService struct {
	roleRepo           *repository.RoleRepository
	permissionRepo     *repository.PermissionRepository
	rolePermissionRepo *repository.RolePermissionRepository
}

func NewRolePermissionService(
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
) *RolePermissionService {
	return &RolePermissionService{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
	}
}

func (s *RolePermissionService) UpdateRolePermissions(
	ctx context.Context,
	roleUUID string,
	req request.UpdateRolePermissionsRequest,
) error {

	roleUUID = strings.TrimSpace(roleUUID)

	if roleUUID == "" {
		return errors.New("role uuid is required")
	}

	// 1. Find role
	role, err := s.roleRepo.FindByUUID(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}

		return err
	}

	// 2. Ambil permission UUID dari request
	permissionUUIDs := make(map[string]struct{})

	for _, permissionUUID := range req.PermissionUUIDs {
		permissionUUID = strings.TrimSpace(permissionUUID)

		if permissionUUID == "" {
			continue
		}

		permissionUUIDs[permissionUUID] = struct{}{}
	}

	// 3. Ambil permission yang sekarang terpasang
	existingRolePermissions, err := s.rolePermissionRepo.FindByRoleID(
		ctx,
		role.ID,
	)
	if err != nil {
		return err
	}

	// 4. Mapping existing permission berdasarkan PermissionID
	existingPermissionIDs := make(map[uint64]bool)

	for _, rolePermission := range existingRolePermissions {
		existingPermissionIDs[rolePermission.PermissionID] = true
	}

	// 5. Resolve semua permission UUID dari request
	requestPermissionIDs := make(map[uint64]bool)

	for permissionUUID := range permissionUUIDs {

		permission, err := s.permissionRepo.FindByUUID(
			ctx,
			permissionUUID,
		)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("permission not found")
			}

			return err
		}

		requestPermissionIDs[permission.ID] = true
	}

	// 6. Tambahkan permission baru
	for permissionID := range requestPermissionIDs {

		if existingPermissionIDs[permissionID] {
			continue
		}

		rolePermission := &model.RolePermission{
			RoleID:       role.ID,
			PermissionID: permissionID,
		}

		if err := s.rolePermissionRepo.Create(
			ctx,
			rolePermission,
		); err != nil {
			return fmt.Errorf(
				"failed to assign permission: role_id=%d permission_id=%d: %w",
				role.ID,
				permissionID,
				err,
			)
		}
	}

	// 7. Hapus permission yang sudah tidak dicentang
	for permissionID := range existingPermissionIDs {

		if requestPermissionIDs[permissionID] {
			continue
		}

		if err := s.rolePermissionRepo.Delete(
			ctx,
			role.ID,
			permissionID,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *RolePermissionService) GetRolePermissions(
	ctx context.Context,
	roleUUID string,
) ([]response.RolePermissionResponse, error) {

	roleUUID = strings.TrimSpace(roleUUID)

	if roleUUID == "" {
		return nil, errors.New("role uuid is required")
	}

	// Find role
	role, err := s.roleRepo.FindByUUID(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}

		return nil, err
	}

	// Find role permissions
	rolePermissions, err := s.rolePermissionRepo.FindByRoleID(
		ctx,
		role.ID,
	)
	if err != nil {
		return nil, err
	}

	result := make([]response.RolePermissionResponse, 0)

	for _, rolePermission := range rolePermissions {

		permission, err := s.permissionRepo.FindByID(
			ctx,
			rolePermission.PermissionID,
		)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			return nil, err
		}

		result = append(result, response.RolePermissionResponse{
			PermissionUUID: permission.UUID,
			Name:           permission.Name,
			Description:    permission.Description,
		})
	}

	return result, nil
}

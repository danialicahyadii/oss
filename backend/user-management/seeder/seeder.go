package seeder

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"oss.kftd.co.id/v2/user-management/model"
)

func SeedUserManagement(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// =========================================================
		// 1. PERMISSIONS
		// =========================================================

		permissions := []struct {
			Name        string
			Description string
		}{
			{
				Name:        "user.view",
				Description: "View users",
			},
			{
				Name:        "user.create",
				Description: "Create user",
			},
			{
				Name:        "user.update",
				Description: "Update user",
			},
			{
				Name:        "user.delete",
				Description: "Delete user",
			},
			{
				Name:        "role.view",
				Description: "View roles",
			},
			{
				Name:        "role.create",
				Description: "Create role",
			},
			{
				Name:        "role.update",
				Description: "Update role",
			},
			{
				Name:        "role.delete",
				Description: "Delete role",
			},
			{
				Name:        "role.permission.view",
				Description: "View role permissions",
			},
			{
				Name:        "role.permission.update",
				Description: "Update role permissions",
			},
			{
				Name:        "user.role.view",
				Description: "View user roles",
			},
			{
				Name:        "user.role.update",
				Description: "Update user roles",
			},
		}

		permissionMap := make(map[string]*model.Permission)

		for _, item := range permissions {
			description := item.Description

			var permission model.Permission

			err := tx.
				Where("name = ?", item.Name).
				First(&permission).Error

			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}

				permission = model.Permission{
					UUID:        uuid.New().String(),
					Name:        item.Name,
					Description: &description,
					IsActive:    true,
				}

				if err := tx.Create(&permission).Error; err != nil {
					return err
				}
			}

			permissionMap[item.Name] = &permission
		}

		// =========================================================
		// 2. ROLES
		// =========================================================

		roleDefinitions := []struct {
			Name        string
			Description string
		}{
			{
				Name:        "Super Admin",
				Description: "Full access to the application",
			},
			{
				Name:        "Admin",
				Description: "Administrative access",
			},
		}

		roleMap := make(map[string]*model.Role)

		for _, item := range roleDefinitions {
			description := item.Description

			var role model.Role

			err := tx.
				Where("name = ?", item.Name).
				First(&role).Error

			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}

				role = model.Role{
					UUID:        uuid.New().String(),
					Name:        item.Name,
					Description: &description,
					IsActive:    true,
				}

				if err := tx.Create(&role).Error; err != nil {
					return err
				}
			}

			roleMap[item.Name] = &role
		}

		// =========================================================
		// 3. ROLE -> PERMISSIONS
		// =========================================================

		// Super Admin mendapatkan semua permission.
		superAdmin := roleMap["Super Admin"]

		for _, permission := range permissionMap {
			var count int64

			err := tx.
				Table("um_role_permissions").
				Where(
					"role_id = ? AND permission_id = ?",
					superAdmin.ID,
					permission.ID,
				).
				Count(&count).Error

			if err != nil {
				return err
			}

			if count == 0 {
				err := tx.Table("um_role_permissions").Create(map[string]interface{}{
					"role_id":       superAdmin.ID,
					"permission_id": permission.ID,
				}).Error

				if err != nil {
					return err
				}
			}
		}

		// Admin permissions.
		admin := roleMap["Admin"]

		adminPermissions := []string{
			"user.view",
			"user.create",
			"user.update",
			"user.delete",

			"role.view",
			"role.create",
			"role.update",
			"role.delete",

			"role.permission.view",
			"role.permission.update",

			"user.role.view",
			"user.role.update",
		}

		for _, permissionName := range adminPermissions {
			permission := permissionMap[permissionName]

			var count int64

			err := tx.
				Table("um_role_permissions").
				Where(
					"role_id = ? AND permission_id = ?",
					admin.ID,
					permission.ID,
				).
				Count(&count).Error

			if err != nil {
				return err
			}

			if count == 0 {
				err := tx.Table("um_role_permissions").Create(map[string]interface{}{
					"role_id":       admin.ID,
					"permission_id": permission.ID,
				}).Error

				if err != nil {
					return err
				}
			}
		}

		// =========================================================
		// 4. USERS
		// =========================================================

		password, err := bcrypt.GenerateFromPassword(
			[]byte("ITkftd@2026"),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}

		users := []struct {
			Username string
			Name     string
			LastName string
			Email    string
			RoleName string
		}{
			{
				Username: "superadmin",
				Name:     "Super",
				LastName: "Admin",
				Email:    "superadmin@example.com",
				RoleName: "Super Admin",
			},
			{
				Username: "admin",
				Name:     "Admin",
				LastName: "User",
				Email:    "admin@example.com",
				RoleName: "Admin",
			},
		}

		for _, item := range users {
			var user model.User

			err := tx.
				Where("username = ?", item.Username).
				First(&user).Error

			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}

				user = model.User{
					UUID:     uuid.New().String(),
					Username: item.Username,
					Name:     item.Name,
					LastName: item.LastName,
					Email:    item.Email,
					Password: string(password),
					IsActive: true,
				}

				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}

			// =====================================================
			// 5. USER -> ROLE
			// =====================================================

			role := roleMap[item.RoleName]

			var count int64

			err = tx.
				Table("um_user_roles").
				Where(
					"user_id = ? AND role_id = ?",
					user.ID,
					role.ID,
				).
				Count(&count).Error

			if err != nil {
				return err
			}

			if count == 0 {
				err := tx.Table("um_user_roles").Create(map[string]interface{}{
					"user_id": user.ID,
					"role_id": role.ID,
				}).Error

				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

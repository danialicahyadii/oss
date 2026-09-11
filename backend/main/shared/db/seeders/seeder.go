package seeders

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	rbacModel "oss.kftd.co.id/v2/user-management/model"
	userModel "oss.kftd.co.id/v2/user-management/model"
)

func Seed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// =========================================================
		// 1. PERMISSIONS
		// =========================================================

		permissions := []rbacModel.Permission{
			{
				UUID:        uuid.NewString(),
				Name:        "role.delete",
				Description: stringPtr("permission delete pada fitur role"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "role.update",
				Description: stringPtr("permission update pada fitur role"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "role.read",
				Description: stringPtr("permission read pada fitur role"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "role.create",
				Description: stringPtr("permission create pada fitur role"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "user.delete",
				Description: stringPtr("permission delete pada fitur user"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "user.update",
				Description: stringPtr("permission update pada fitur user"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "user.create",
				Description: stringPtr("permission create pada fitur user"),
				IsActive:    true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "user.read",
				Description: stringPtr("permission read pada fitur user"),
				IsActive:    true,
			},
		}

		for _, permission := range permissions {
			var existing rbacModel.Permission

			err := tx.
				Where("name = ?", permission.Name).
				First(&existing).
				Error

			if err == nil {
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := tx.Create(&permission).Error; err != nil {
				return err
			}
		}

		// =========================================================
		// 2. ROLES
		// =========================================================

		// adminDescription := "Admin adalah role yang memiliki semua permission"
		// userDescription := "User adalah role yang memiliki beberapa permission tertentu"

		roles := []rbacModel.Role{
			{
				UUID:        uuid.NewString(),
				Name:        "Admin",
				Description: stringPtr(
					"Admin adalah role yang memiliki semua permission",
				),
				IsActive: true,
			},
			{
				UUID:        uuid.NewString(),
				Name:        "User",
				Description: stringPtr(
					"User adalah role yang memiliki beberapa permission tertentu",
				),
				IsActive: true,
			},
		}

		for _, role := range roles {
			var existing rbacModel.Role

			err := tx.
				Where("name = ?", role.Name).
				First(&existing).
				Error

			if err == nil {
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := tx.Create(&role).Error; err != nil {
				return err
			}
		}

		// Ambil role yang sudah ada
		var adminRole rbacModel.Role
		if err := tx.
			Where("name = ?", "Admin").
			First(&adminRole).Error; err != nil {
			return err
		}

		var userRole rbacModel.Role
		if err := tx.
			Where("name = ?", "User").
			First(&userRole).Error; err != nil {
			return err
		}

		// =========================================================
		// 3. ROLE PERMISSIONS
		// =========================================================

		var allPermissions []rbacModel.Permission

		if err := tx.
			Where("deleted_at IS NULL").
			Find(&allPermissions).Error; err != nil {
			return err
		}

		// ADMIN = SEMUA PERMISSION
		for _, permission := range allPermissions {
			if err := seedRolePermission(
				tx,
				adminRole.ID,
				permission.ID,
			); err != nil {
				return err
			}
		}

		// USER = user.read
		var userReadPermission rbacModel.Permission

		if err := tx.
			Where("name = ?", "user.read").
			First(&userReadPermission).Error; err != nil {
			return err
		}

		if err := seedRolePermission(
			tx,
			userRole.ID,
			userReadPermission.ID,
		); err != nil {
			return err
		}

		// =========================================================
		// 4. USERS
		// =========================================================

		passwordHash, err := bcrypt.GenerateFromPassword(
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
			RoleID   uint64
		}{
			{
				Username: "admin_oss",
				Name:     "Admin",
				LastName: "Oss",
				Email:    "admin_oss@kftd.co.id",
				RoleID:   adminRole.ID,
			},
			{
				Username: "user_oss",
				Name:     "User",
				LastName: "Oss",
				Email:    "user_oss@kftd.co.id",
				RoleID:   userRole.ID,
			},
		}

		for _, data := range users {
			var user userModel.User

			err := tx.
				Where("email = ?", data.Email).
				First(&user).
				Error

			if err == nil {
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			user = userModel.User{
				UUID:      uuid.NewString(),
				Username:  data.Username,
				Name:      data.Name,
				LastName:  data.LastName,
				Email:     data.Email,
				Password:  string(passwordHash),
				IsActive:  true,
			}

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			// User Role
			if err := seedUserRole(
				tx,
				user.ID,
				data.RoleID,
			); err != nil {
				return err
			}
		}

		log.Println("Database seeder completed successfully")

		return nil
	})
}

func seedRolePermission(
	tx *gorm.DB,
	roleID uint64,
	permissionID uint64,
) error {

	// Sesuaikan nama tabel dengan tabel pivot kamu
	var count int64

	err := tx.Table("um_role_permissions").
		Where(
			"role_id = ? AND permission_id = ?",
			roleID,
			permissionID,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return tx.Table("um_role_permissions").Create(map[string]interface{}{
		"role_id":       roleID,
		"permission_id": permissionID,
	}).Error
}

func seedUserRole(
	tx *gorm.DB,
	userID uint64,
	roleID uint64,
) error {

	var count int64

	err := tx.Table("um_user_roles").
		Where(
			"user_id = ? AND role_id = ?",
			userID,
			roleID,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return tx.Table("um_user_roles").Create(map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
	}).Error
}

func stringPtr(value string) *string {
	return &value
}

func init() {
	fmt.Print("")
}
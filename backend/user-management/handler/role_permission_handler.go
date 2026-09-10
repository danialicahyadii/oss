package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/service"

	"gorm.io/gorm"
)

type RolePermissionHandler struct {
	service *service.RolePermissionService
}

func NewRolePermissionHandler(
	service *service.RolePermissionService,
) *RolePermissionHandler {
	return &RolePermissionHandler{
		service: service,
	}
}

func (h *RolePermissionHandler) UpdateRolePermissions(c *gin.Context) {

	roleUUID := c.Param("uuid")

	var req request.UpdateRolePermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := h.service.UpdateRolePermissions(
		c.Request.Context(),
		roleUUID,
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"message": "role or permission not found",
			})

		case err.Error() == "role uuid is required":
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})

		case err.Error() == "role not found":
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})

		case err.Error() == "permission not found":
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role permissions updated successfully",
	})
}

func (h *RolePermissionHandler) GetRolePermissions(c *gin.Context) {

	roleUUID := c.Param("uuid")

	data, err := h.service.GetRolePermissions(
		c.Request.Context(),
		roleUUID,
	)

	if err != nil {
		switch err.Error() {
		case "role uuid is required":
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})

		case "role not found":
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": data,
	})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/service"
)

type PermissionHandler struct {
	service *service.PermissionService
}

func NewPermissionHandler(
	service *service.PermissionService,
) *PermissionHandler {
	return &PermissionHandler{
		service: service,
	}
}

func (h *PermissionHandler) GetPermissions(c *gin.Context) {

	search := c.Query("search")

	page, err := strconv.Atoi(
		c.DefaultQuery("page", "1"),
	)

	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(
		c.DefaultQuery("limit", "10"),
	)

	if err != nil {
		limit = 10
	}

	permissions, total, err := h.service.GetPermissions(
		c.Request.Context(),
		search,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get permissions",
		})
		return
	}

	data := make([]gin.H, 0, len(permissions))

	for _, permission := range permissions {
		data = append(data, gin.H{
			"uuid":        permission.UUID,
			"name":        permission.Name,
			"description": permission.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data":    data,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *PermissionHandler) GetPermission(c *gin.Context) {

	permissionUUID := c.Param("uuid")

	permission, err := h.service.GetPermission(
		c.Request.Context(),
		permissionUUID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data": gin.H{
			"uuid":        permission.UUID,
			"name":        permission.Name,
			"description": permission.Description,
		},
	})
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {

	var req request.CreatePermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	permission, err := h.service.CreatePermission(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "permission created successfully",
		"data": gin.H{
			"uuid":        permission.UUID,
			"name":        permission.Name,
			"description": permission.Description,
		},
	})
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {

	permissionUUID := c.Param("uuid")

	var req request.UpdatePermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdatePermission(
		c.Request.Context(),
		permissionUUID,
		req,
	); err != nil {

		status := http.StatusBadRequest

		if err.Error() == "permission not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "permission updated successfully",
	})
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {

	permissionUUID := c.Param("uuid")

	if err := h.service.DeletePermission(
		c.Request.Context(),
		permissionUUID,
	); err != nil {

		status := http.StatusBadRequest

		if err.Error() == "permission not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "permission deleted successfully",
	})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/service"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(
	service *service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

func (h *RoleHandler) GetRoles(c *gin.Context) {

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

	roles, total, err := h.service.GetRoles(
		c.Request.Context(),
		search,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get roles",
		})
		return
	}

	data := make([]gin.H, 0, len(roles))

	for _, role := range roles {
		data = append(data, gin.H{
			"uuid":        role.UUID,
			"name":        role.Name,
			"description": role.Description,
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

func (h *RoleHandler) GetRole(c *gin.Context) {

	roleUUID := c.Param("uuid")

	role, err := h.service.GetRole(
		c.Request.Context(),
		roleUUID,
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
			"uuid":        role.UUID,
			"name":        role.Name,
			"description": role.Description,
		},
	})
}

func (h *RoleHandler) CreateRole(c *gin.Context) {

	var req request.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	role, err := h.service.CreateRole(
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
		"message": "role created successfully",
		"data": gin.H{
			"uuid":        role.UUID,
			"name":        role.Name,
			"description": role.Description,
		},
	})
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {

	roleUUID := c.Param("uuid")

	var req request.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdateRole(
		c.Request.Context(),
		roleUUID,
		req,
	); err != nil {

		status := http.StatusBadRequest

		if err.Error() == "role not found" {
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
		"message": "role updated successfully",
	})
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {

	roleUUID := c.Param("uuid")

	if err := h.service.DeleteRole(
		c.Request.Context(),
		roleUUID,
	); err != nil {

		status := http.StatusBadRequest

		if err.Error() == "role not found" {
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
		"message": "role deleted successfully",
	})
}

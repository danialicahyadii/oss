package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/service"
)

type UserRoleHandler struct {
	service *service.UserRoleService
}

func NewUserRoleHandler(
	service *service.UserRoleService,
) *UserRoleHandler {
	return &UserRoleHandler{
		service: service,
	}
}

func (h *UserRoleHandler) UpdateUserRoles(c *gin.Context) {

	userUUID := c.Param("uuid")

	var req request.UpdateUserRolesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := h.service.UpdateUserRoles(
		c.Request.Context(),
		userUUID,
		req,
	)

	if err != nil {
		switch err.Error() {
		case "user uuid is required":
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})

		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{
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
		"message": "user roles updated successfully",
	})
}

func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {

	userUUID := c.Param("uuid")

	data, err := h.service.GetUserRoles(
		c.Request.Context(),
		userUUID,
	)

	if err != nil {
		switch err.Error() {
		case "user uuid is required":
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})

		case "user not found":
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
		"roles": data,
	})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(
	service *service.UserService,
) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := request.UserListRequest{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
	}

	if value := c.Query("is_active"); value != "" {
		isActive := value == "true"
		req.IsActive = &isActive
	}

	result, err := h.service.GetUsers(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {

	userUUID := c.Param("uuid")

	result, err := h.service.GetUser(
		c.Request.Context(),
		userUUID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := h.service.CreateUser(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"data":    result,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {

	userUUID := c.Param("uuid")

	var req request.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result, err := h.service.UpdateUser(
		c.Request.Context(),
		userUUID,
		req,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"data":    result,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

	userUUID := c.Param("uuid")

	if err := h.service.DeleteUser(
		c.Request.Context(),
		userUUID,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}

package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"oss.kftd.co.id/v2/user-management/dto/request"
	"oss.kftd.co.id/v2/user-management/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(
	service *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	token, err := h.service.Login(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data": gin.H{
			"token": token,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	user, err := h.service.Register(
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
		"message": "registration successful",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user identity",
		})
		return
	}

	user, err := h.service.Me(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"name":       user.Name,
			"last_name":  user.LastName,
			"email":      user.Email,
			"is_active":  user.IsActive,
		},
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user identity",
		})
		return
	}

	var req request.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.ChangePassword(
		c.Request.Context(),
		userID,
		req,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "password changed successfully",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "authorization header is required",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid authorization header",
		})
		return
	}

	token, err := h.service.RefreshToken(parts[1])

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "token refreshed",
		"data": gin.H{
			"token": token,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logout successful",
	})
}
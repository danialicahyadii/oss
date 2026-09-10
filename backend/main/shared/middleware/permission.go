package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"oss.kftd.co.id/v2/user-management/service"
)

func RequirePermission(
	permissionService *service.PermissionService,
	permission string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		userIDValue, exists := c.Get("user_id")

		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		userID, ok := userIDValue.(uint64)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		hasPermission, err := permissionService.HasPermission(
			c.Request.Context(),
			userID,
			permission,
		)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "permission denied",
			})
			return
		}

		c.Next()
	}
}

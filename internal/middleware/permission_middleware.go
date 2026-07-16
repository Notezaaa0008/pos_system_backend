package middleware

import (
	"errors"
	"log"
	"net/http"
	"pos-system-backend/internal/models"
	"pos-system-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MiddlewarePermissionService interface {
	ValidatePermissionService(userIDStr string, storeIDStr string) (*models.UserStore, error)
}

func PermissionMiddleware(authService MiddlewarePermissionService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, userExists := c.Get("userID")
		systemRole, roleExists := c.Get("systemRole")

		if !userExists || !roleExists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized"})
			c.Abort()
			return
		}

		userIDStr := userID.(string)
		systemRoleStr := systemRole.(string)
		if systemRoleStr == "SYSTEM_ADMIN" {
			storeID := c.GetHeader("X-Store-ID")
    		if storeID != "" {
        		c.Set("storeID", storeID)
    		}
			c.Set("storeRole", "SYSTEM_ADMIN")
			c.Next()
			return
		}

		storeIDStr := c.GetHeader("X-Store-ID")
		if storeIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing X-Store-ID header"})
			c.Abort()
			return
		}

		userStore, err := authService.ValidatePermissionService(userIDStr, storeIDStr)

		if err != nil {
			var appErr *utils.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, gin.H{
					"status":  "error",
					"message": appErr.Message,
				})
				c.Abort()
				return
			}

			// หาข้อมูลสิทธิ์ในตารางกลางไม่เจอ (แปลว่าไม่มีสิทธิ์ในร้านนี้)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[RBAC WARN] Access Denied - UserID: %s has no permission for StoreID: %s", userIDStr, storeIDStr)
				c.JSON(http.StatusForbidden, gin.H{
					"status":  "error",
					"message": "forbidden: you do not have permission for this store",
				})
				c.Abort()
				return
			}

			// 🔍 เคสที่ 3: Database พัง หรือเกิดระบบขัดข้องที่คาดไม่ถึง
			log.Printf("[RBAC ERROR] Unexpected error during permission validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "internal server error: please try again later",
			})
			c.Abort()
			return
		}

		currentRole := userStore.Role.RoleName
		isAllowed := false
		for _, role := range allowedRoles {
			if currentRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			log.Printf("[RBAC WARN] UserID %v with Role %s tried to access restriction %v", userIDStr, currentRole, allowedRoles)
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "forbidden: insufficient permissions for this store",
			})
			c.Abort()
			return
		}

		c.Set("storeID", userStore.StoreID)
		c.Set("storeRole", currentRole)
		c.Set("storeRoleID", userStore.RoleID)
		c.Next()
	}

}
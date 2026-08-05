package routes

import (
	"pos-system-backend/internal/middleware"
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/users"

	"github.com/gin-gonic/gin"
)

func initUserRoutes(routesGroup *gin.RouterGroup, userCtrl *users.UsersController, authService *auth.AuthService, authMiddleware gin.HandlerFunc) {

	user := routesGroup.Group("/users")
	{
		protectedUser := user.Group("/")
		protectedUser.Use(authMiddleware)
		{
			protectedUser.GET("/profile", middleware.PermissionMiddleware(authService, "OWNER", "MANAGER", "STAFF", "CASHIER"), userCtrl.GetProfileController)
		}
	}
}
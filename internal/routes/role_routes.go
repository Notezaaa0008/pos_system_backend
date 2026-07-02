package routes

import (
	"pos-system-backend/internal/middleware"
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/roles"

	"github.com/gin-gonic/gin"
)

func initRoleRoutes(routesGroup *gin.RouterGroup, roleCtrl *roles.RolesController, authService *auth.AuthService, authMiddleware gin.HandlerFunc) {
	role := routesGroup.Group("/role")
	{
		protectedRole := role.Group("/")
		protectedRole.Use(authMiddleware)
		{
			protectedRole.GET("/all-role", middleware.PermissionMiddleware(authService, "SYSTEM_ADMIN"), roleCtrl.GetAllRolesController)

			protectedRole.POST("/create", middleware.PermissionMiddleware(authService, "SYSTEM_ADMIN"), roleCtrl.CreateRoleController)
		
			protectedRole.PUT("/update/:id", middleware.PermissionMiddleware(authService, "SYSTEM_ADMIN"), roleCtrl.UpdateRoleController)
		} 
		
	}
}
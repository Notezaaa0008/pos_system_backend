package routes

import (
	"pos-system-backend/internal/module/auth"

	"github.com/gin-gonic/gin"
)

func initAuthRoutes(routesGroup *gin.RouterGroup, authCtrl *auth.AuthController, authService *auth.AuthService, middleware gin.HandlerFunc) {

	
	
	auth := routesGroup.Group("/auth")
	{
		auth.POST("/signup-system-admin", authCtrl.RegisterSystemAdminController)
		auth.POST("/login", authCtrl.LoginController)
		auth.POST("/forgot-password", authCtrl.ForgotPasswordController)
        auth.POST("/reset-password", authCtrl.ResetPasswordController)

		protectedAuth := auth.Group("/")
		protectedAuth.Use(middleware) 
		{
			// middleware.permissionMiddleware(authService, "MANAGER", "CASHIER")
			protectedAuth.POST("/logout", authCtrl.LogoutController)
		}
		
	}
}
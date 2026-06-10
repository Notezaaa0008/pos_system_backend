package routes

import (
	"gin-quickstart/internal/module/auth"

	"github.com/gin-gonic/gin"
)

func initAuthRoutes(routesGroup *gin.RouterGroup, authCtrl *auth.AuthController, requireAuth gin.HandlerFunc) {

	
	
	auth := routesGroup.Group("/auth")
	{
		auth.POST("/signup-system-admin", authCtrl.RegisterSystemAdminController)
		auth.POST("/login", authCtrl.LoginController)
	}
}
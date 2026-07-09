package routes

import (
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/master"

	"github.com/gin-gonic/gin"
)

// gin.HandlerFunc
func initMasterRoutes(routesGroup *gin.RouterGroup, masterCtrl *master.MasterController, authService *auth.AuthService, authMiddleware gin.HandlerFunc) {

	
	
	master := routesGroup.Group("/master")
	{
		protectedPrefix := master.Group("/")
		protectedPrefix.Use(authMiddleware)
		master.GET("/all-prefix", masterCtrl.GetAllPrefixController)
		
	}
}
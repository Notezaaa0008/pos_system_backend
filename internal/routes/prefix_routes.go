package routes

import (
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/prefix"

	"github.com/gin-gonic/gin"
)

// gin.HandlerFunc
func initPrefixRoutes(routesGroup *gin.RouterGroup, perfixCtrl *prefix.PrefixController, authService *auth.AuthService, authMiddleware gin.HandlerFunc) {

	
	
	prefix := routesGroup.Group("/prefix")
	{
		protectedPrefix := prefix.Group("/")
		protectedPrefix.Use(authMiddleware)
		prefix.GET("/all-prefix", perfixCtrl.GetAllPrefixController)
		
	}
}
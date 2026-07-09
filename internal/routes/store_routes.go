package routes

import (
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/store"

	"github.com/gin-gonic/gin"
)

func initStoreRoutes(routesGroup *gin.RouterGroup, StoreCtrl *store.StoreController, authService *auth.AuthService, authMiddleware gin.HandlerFunc) {

	store := routesGroup.Group("/prefix")
	{
		protectedPrefix := store.Group("/")
		protectedPrefix.Use(authMiddleware)
		store.GET("/store", StoreCtrl.GetUserStoreController)

	}
}
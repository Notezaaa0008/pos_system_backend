package routes

import (
	"pos-system-backend/internal/middleware"
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/store"

	"github.com/gin-gonic/gin"
)

func initStoreRoutes(routesGroup *gin.RouterGroup, storeCtrl *store.StoreController, authService *auth.AuthService, authMiddleware gin.HandlerFunc) {

	store := routesGroup.Group("/store")
	{
		protectedStore := store.Group("/")
		protectedStore.Use(authMiddleware)
		{
			protectedStore.GET("/all-store", storeCtrl.GetStoreListController)
			protectedStore.POST("/create-store", middleware.PermissionMiddleware(authService, "OWNER"), storeCtrl.CreateStoreController)
			protectedStore.PUT("/update-store", middleware.PermissionMiddleware(authService, "OWNER"), storeCtrl.UpdateStoreController)
			protectedStore.PATCH("/update-status", middleware.PermissionMiddleware(authService, "OWNER"), storeCtrl.UpdateStoreStatusController)
			protectedStore.DELETE("delete-store", middleware.PermissionMiddleware(authService, "OWNER"), storeCtrl.DeleteStoreController)
		}
		

	}
}
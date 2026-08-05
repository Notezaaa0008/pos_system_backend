package routes

import (
	"pos-system-backend/internal/middleware"
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/master"
	"pos-system-backend/internal/module/store"
	"pos-system-backend/internal/module/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func InitRouter(server *gin.Engine, db *gorm.DB) {
	v1 := server.Group("/api/v1")

	// init module
	authModule := auth.InitModule(db)
	masterModule := master.InitModule(db)
	storeModule := store.InitModule(db)
	userModule := users.InitModule(db)

	authMiddleware := middleware.AuthMiddleware(authModule.Service)
	
	initAuthRoutes(v1, authModule.Controller, authModule.Service, authMiddleware)
	initMasterRoutes(v1, masterModule.Controller, authModule.Service, authMiddleware)
	initStoreRoutes(v1, storeModule.Controller, authModule.Service, authMiddleware)
	initUserRoutes(v1, userModule.Controller, authModule.Service, authMiddleware)
	
}
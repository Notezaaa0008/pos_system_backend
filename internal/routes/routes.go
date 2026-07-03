package routes

import (
	"pos-system-backend/internal/middleware"
	"pos-system-backend/internal/module/auth"
	"pos-system-backend/internal/module/prefix"
	"pos-system-backend/internal/module/roles"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func InitRouter(server *gin.Engine, db *gorm.DB) {
	v1 := server.Group("/api/v1")

	// init module
	roleModule := roles.InitModule(db)
	authModule := auth.InitModule(db, roleModule)
	prefixModule := prefix.InitModule(db)

	authMiddleware := middleware.AuthMiddleware(authModule.Service)
	
	initAuthRoutes(v1, authModule.Controller, authModule.Service, authMiddleware)
	initRoleRoutes(v1, roleModule.Controller, authModule.Service, authMiddleware)
	initPrefixRoutes(v1, prefixModule.Controller, authModule.Service, authMiddleware)
	initUserRoutes(v1, db)
	
}
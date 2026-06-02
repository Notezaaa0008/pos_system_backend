package routes

import (
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/module/auth"
	"gin-quickstart/internal/module/roles"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func InitRouter(server *gin.Engine, db *gorm.DB) {
	v1 := server.Group("/api/v1")

	// init module
	roleModule := roles.InitModule(db)
	authModule := auth.InitModule(db, roleModule)

	authMiddleware := middleware.AuthWithRefreshMiddleware(authModule.Service)
	
	initAuthRoutes(v1, authModule.Controller, authMiddleware)
	initUserRoutes(v1, db)
	
}
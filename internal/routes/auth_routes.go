package routes

import (
	"gin-quickstart/internal/module/auth"
	"gin-quickstart/internal/module/roles"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initAuthRoutes(routesGroup *gin.RouterGroup, db *gorm.DB) {

	roleMod := roles.InitModule(db)
	authMod := auth.InitModule(db, roleMod)
	
	auth := routesGroup.Group("/auth")
	{
		auth.POST("/signup-super-admin", authMod.Controller.SignupSuperAdminController)
	}
}
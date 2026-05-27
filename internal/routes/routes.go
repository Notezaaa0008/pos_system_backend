package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func InitRouter(server *gin.Engine, db *gorm.DB) {
	v1 := server.Group("/api/v1")
	
	initAuthRoutes(v1, db)
	initUserRoutes(v1, db)
	
}
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initUserRoutes(routesGroup *gin.RouterGroup, db *gorm.DB) {

	// userMod := users.InitModule(db)

	user := routesGroup.Group("/users")
	{
		user.GET("/")
	}
}
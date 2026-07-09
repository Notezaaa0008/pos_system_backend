package users

import "github.com/gin-gonic/gin"

type UsersController struct {
	service *UsersService
}

func NewUsersController (service *UsersService) *UsersController{
	return &UsersController{service: service}
}

func (userClrt *UsersController) GetProfile(c *gin.Context) {

}




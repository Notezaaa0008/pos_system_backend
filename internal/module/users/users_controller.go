package users

type UsersController struct {
	service *UsersService
}

func NewUsersController (service *UsersService) *UsersController{
	return &UsersController{service: service}
}






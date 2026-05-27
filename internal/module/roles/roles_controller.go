package roles


type RolesController struct {
	service *RolesService
}

func NewRoleController (service *RolesService) *RolesController{
	return &RolesController{service: service}
}
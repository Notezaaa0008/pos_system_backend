package roles

import "gorm.io/gorm"

type Module struct {
    Repo    	*RolesRepository
    Service 	*RolesService
	Controller 	*RolesController
}

func InitModule(db *gorm.DB) *Module {
    repo := NewRolesRepository(db)
    service := NewRolesService(repo)
	controller := NewRoleController(service)
    
    return &Module{
        Repo:    	repo,
        Service: 	service,
		Controller: controller,
    }
}
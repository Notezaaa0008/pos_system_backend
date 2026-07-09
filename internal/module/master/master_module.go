package master

import "gorm.io/gorm"

type Module struct{
	Repo       *MasterRepository
    Service    *MasterService
    Controller *MasterController
}

func InitModule(db *gorm.DB) *Module {
    repo := NewMasterRepository(db)
    service := NewMasterService(repo) 
    controller := NewMasterController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}

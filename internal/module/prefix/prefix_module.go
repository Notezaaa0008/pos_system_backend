package prefix

import "gorm.io/gorm"

type Module struct{
	Repo       *PrefixRepository
    Service    *PrefixService
    Controller *PrefixController
}

func InitModule(db *gorm.DB) *Module {
    repo := NewPrefixRepository(db)
    
    // ดึงเฉพาะตัว Service ของ Role ออกมาจากกระเป๋าเพื่อเอามาต่อสายไฟ 🔌
    service := NewPrefixService(repo) 
    controller := NewPrefixController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}

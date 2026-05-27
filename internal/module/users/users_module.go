package users

import (
	"gorm.io/gorm"
)

type Module struct {
    Repo       *UsersRepository
    Service    *UsersService
    Controller *UsersController
}

func InitModule(db *gorm.DB) *Module {
    repo := NewUserRepository(db)
    
    // ดึงเฉพาะตัว Service ของ Role ออกมาจากกระเป๋าเพื่อเอามาต่อสายไฟ 🔌
    service := NewUsersService(repo) 
    controller := NewUsersController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}
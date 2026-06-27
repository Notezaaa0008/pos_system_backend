package auth

import (
	"pos-system-backend/internal/module/roles"

	"gorm.io/gorm"
)

type Module struct {
    Repo       *AuthRepository
    Service    *AuthService
    Controller *AuthController
}


func InitModule(db *gorm.DB, roleModule *roles.Module) *Module {
    repo := NewAuthRepository(db)
    
    // ดึงเฉพาะตัว Service ของ Role ออกมาจากกระเป๋าเพื่อเอามาต่อสายไฟ 🔌
    service := NewAuthService(repo, roleModule.Service) 
    controller := NewAuthController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}

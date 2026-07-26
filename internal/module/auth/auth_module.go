package auth

import (
	"gorm.io/gorm"
)

type Module struct {
    Repo       *AuthRepository
    Service    *AuthService
    Controller *AuthController
}


func InitModule(db *gorm.DB) *Module {
    repo := NewAuthRepository(db)
    service := NewAuthService(repo) 
    controller := NewAuthController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}

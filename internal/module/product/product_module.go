package product

import "gorm.io/gorm"

type Module struct{
	Repo       *ProductRepository
    Service    *ProductService
    Controller *ProductController
}

func InitModule(db *gorm.DB) *Module {
    repo := NewProductRepository(db)
    service := NewProductService(repo) 
    controller := NewProductController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}
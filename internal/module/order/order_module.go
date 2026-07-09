package order

import "gorm.io/gorm"

type Module struct{
	Repo       *OrderRepository
    Service    *OrderService
    Controller *OrderController
}

func InitModule(db *gorm.DB) *Module {
    repo := NewOrderRepository(db)
    service := NewOrderService(repo) 
    controller := NewOrderController(service)
    
    return &Module{
        Repo:       repo,
        Service:    service,
        Controller: controller,
    }
}
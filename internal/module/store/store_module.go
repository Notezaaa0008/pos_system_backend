package store

import "gorm.io/gorm"

type Module struct {
	Repo       *StoreRepository
	Service    *StoreService
	Controller *StoreController
}

func InitModule(db *gorm.DB) *Module {
	repo := NewStoreRepository(db)
	service := NewStoreService(repo)
	controller := NewStoreController(service)

	return &Module{
		Repo:       repo,
		Service:    service,
		Controller: controller,
	}
}
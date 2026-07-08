package store

type storeServiceInterface interface {
}

type StoreController struct {
	service storeServiceInterface
}

func NewStoreController(service storeServiceInterface) *StoreController {
	return &StoreController{service: service}
}
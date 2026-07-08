package store

type storeServiceInterface interface {
}

type StoreController struct {
	service storeServiceInterface
}

func NewPrefixController(service storeServiceInterface) *StoreController {
	return &StoreController{service: service}
}
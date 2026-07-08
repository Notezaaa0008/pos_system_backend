package store

type storeRepositoryInterface interface {
	
}

type StoreService struct {
	repo storeRepositoryInterface
}

func NewStoreService(repo storeRepositoryInterface) *StoreService {
	return &StoreService{repo: repo}
}
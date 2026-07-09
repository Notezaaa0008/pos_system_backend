package order

type orderRepositoryInterface interface {
	
}

type OrderService struct {
	repo orderRepositoryInterface
}

func NewOrderService(repo orderRepositoryInterface) *OrderService {
	return &OrderService{repo: repo}
}
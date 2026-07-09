package order


type orderServiceInterface interface {
	
}

type OrderController struct {
	service orderServiceInterface
}

func NewOrderController(service orderServiceInterface) *OrderController {
	return &OrderController{service: service}
}
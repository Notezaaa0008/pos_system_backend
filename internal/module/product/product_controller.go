package product


type productServiceInterface interface {
	
}

type ProductController struct {
	service productServiceInterface
}

func NewProductController(service productServiceInterface) *ProductController {
	return &ProductController{service: service}
}
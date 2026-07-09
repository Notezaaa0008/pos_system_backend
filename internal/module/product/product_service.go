package product


type productRepositoryInterface interface {
	
}

type ProductService struct {
	repo productRepositoryInterface
}

func NewProductService(repo productRepositoryInterface) *ProductService {
	return &ProductService{repo: repo}
}
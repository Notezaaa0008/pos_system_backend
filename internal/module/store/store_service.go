package store

import (
	"pos-system-backend/internal/models"
	storeDto "pos-system-backend/internal/module/store/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type storeRepositoryInterface interface {
	GetUserStores(userID uuid.UUID, data *storeDto.GetUserStoreRequest) ([]models.UserStore, int64, error)
}

type StoreService struct {
	repo storeRepositoryInterface
}

func NewStoreService(repo storeRepositoryInterface) *StoreService {
	return &StoreService{repo: repo}
}

func (sevice *StoreService) GetUserStoreService(userID uuid.UUID, req *storeDto.GetUserStoreRequest) ([]gin.H, int64, error){
	userStores, total, err := sevice.repo.GetUserStores(userID, req)
	if err != nil {
		return nil, 0, err
	}

	// []gin.H คือ map[string]interface
	var result []gin.H
	for _, us := range userStores {
		result = append(result, gin.H{
			"store_id":   us.StoreID,
			"store_code": us.Store.StoreCode,
			"store_name": us.Store.StoreName,
			"role_name":  us.Role.RoleName,
		})
	}

	return result, total, nil
}
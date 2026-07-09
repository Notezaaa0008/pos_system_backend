package store

import (
	"fmt"
	"pos-system-backend/internal/models"
	storeDto "pos-system-backend/internal/module/store/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (repo *StoreRepository) GetUserStores(userID uuid.UUID, data *storeDto.GetUserStoreRequest) ([]models.UserStore, int64, error) {
	var userStores []models.UserStore
	var total int64

	if data.Page <= 0 {
        data.Page = 1 
    }
    if data.Limit <= 0 {
        data.Limit = 10 
    } else if data.Limit > 50 {
        data.Limit = 50 
    }

	query := repo.db.Model(&models.UserStore{}).
		Preload("Store").
		Preload("Role").
		Joins("JOIN stores ON stores.id = user_stores.store_id").
		Where("user_stores.user_id = ? AND user_stores.deleted_at IS NULL", userID)

	if data.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", data.Search)
		query = query.Where("stores.store_code LIKE ? OR stores.store_name LIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&total).Error
	if  err != nil {
		return nil, 0, err
	}

	offset := (data.Page - 1) * data.Limit
	err = query.Limit(data.Limit).Offset(offset).Find(&userStores).Error
	if err != nil {
		return nil, 0, err
	}

	return userStores, total, nil
}
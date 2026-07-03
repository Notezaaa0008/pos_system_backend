package prefix

import (
	"pos-system-backend/internal/models"

	"gorm.io/gorm"
)

type PrefixRepository struct {
	db  *gorm.DB
}

func NewPrefixRepository(db *gorm.DB) *PrefixRepository {
	return &PrefixRepository{db: db}
}

func (repo *PrefixRepository) GetAllPrefix() ([]models.Prefix, error) {
	var prefix []models.Prefix

	err := repo.db.Where("is_active = ? AND deleted_at IS NULL", true).Find(&prefix).Error

	if err != nil {
		return nil, err
	}

	return prefix, nil
}
package master

import (
	"pos-system-backend/internal/models"

	"gorm.io/gorm"
)

type MasterRepository struct {
	db  *gorm.DB
}

func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{db: db}
}

func (repo *MasterRepository) GetAllPrefix() ([]models.Prefix, error) {
	var prefix []models.Prefix

	err := repo.db.Where("is_active = ? AND deleted_at IS NULL", true).Find(&prefix).Error

	if err != nil {
		return nil, err
	}

	return prefix, nil
}
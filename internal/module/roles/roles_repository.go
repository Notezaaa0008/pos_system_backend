package roles

import (
	"errors"
	"gin-quickstart/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolesRepository struct {
	db *gorm.DB
}

func NewRolesRepository (db *gorm.DB) *RolesRepository {
	return &RolesRepository{db: db}
}

func (repo *RolesRepository) GetRoleIdByRoleName(roleName string) (uuid.UUID, error) {
	var role models.Role

	err := repo.db.Select("id").Where("role_name = ? AND is_active = ?", roleName, true).First(&role).Error

	if err != nil {
		// ไม่เจอ Role นี้เลยในตาราง ให้มองว่ายังไม่มี super_admin และไม่ถือว่าระบบพัง
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return uuid.Nil, errors.New("system error: super_admin role not found in the database, please configure the system first.") 
        }
		// เกิดข้อผิดพลาดอื่น เช่น DB ล่ม
		return uuid.Nil, err
	}

	return role.ID, nil
}



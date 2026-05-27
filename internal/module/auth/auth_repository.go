package auth

import (
	"errors"
	"gin-quickstart/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db  *gorm.DB
}

func NewAuthRepository (db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) CheckSuperAdminExists(rolesID uuid.UUID) (bool, error) {
	var user models.User
	

	err := repo.db.Select("id").Where("role_id = ?", rolesID).First(&user).Error

	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
		// เกิดข้อผิดพลาดอื่น เช่น DB ล่ม
		return false, err
	}

	// เจอทั้ง Role และเจอทั้ง User ที่มีสิทธิ์นี้
    return true, nil
}

func (repo *AuthRepository) FineUserByUserName(userName string) (*models.User, error){
	var user models.User

	err := repo.db.Preload("Role").Where("user_name = ? AND is_active = ?", userName, true).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil , errors.New("system error: invalid username or password.")
		}

		return nil, err
	}

	return &user, nil
}

func (repo *AuthRepository) CreateUser(user *models.User) error {
	err := repo.db.Create(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateRefreshTokenRecord(refreshToken *models.RefreshToken) error {
	err := repo.db.Create(refreshToken).Error

	if err != nil {
		return err
	}

	return nil
}
package users

import (
	"log"
	"pos-system-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUserRepository (db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (repo *UsersRepository) GetAllUserRepository() {
	
}

func (repo *UsersRepository) GetProfileRepository(userID uuid.UUID) (*models.User, error){
	var user models.User

	err := repo.db.Model(&models.User{}).Where("id = ? AND is_active = ?", userID, true).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (repo *UsersRepository) UpdateProfileRepository(userID uuid.UUID, userData *models.User) error{
	err := repo.db.Model(&models.User{}).
	Where("id = ?", userID).
	Updates(userData).Error

	if err != nil {
		log.Printf("[Repository UpdateProfile ERROR] : %v", err)
        return err
	}

	return nil
}

func (repo *UsersRepository) DeleteUserRepository(userStoreData *models.UserStore, userData *models.User, deleteUserID uuid.UUID) error{
	tx := repo.db.Begin()
    if tx.Error != nil {
        return tx.Error
    }

	defer tx.Rollback()

	// Hard Delete
	err := tx.Unscoped().Where("user_id = ?", deleteUserID).Delete(&models.RefreshToken{}).Error;
	if err != nil {
		return err
	}
	err = tx.Unscoped().Where("user_id = ?", deleteUserID).Delete(&models.ResetPassword{}).Error;
	if err != nil {
		return err
	}

	err = tx.Model(&models.UserStore{}).
    Where("user_id = ?", deleteUserID).
    Select("IsActive", "DeletedBy").
    Updates(userStoreData).Error
	if err != nil {
        return err
    }

	err = tx.Where("user_id = ?", deleteUserID).Delete(&models.UserStore{}).Error
	if err != nil {
        return err
    }

	err = tx.Model(&models.User{}).
    Where("id = ?", deleteUserID).
    Select("IsActive", "DeletedBy").
    Updates(userData).Error
	if err != nil {
        return err
    }

	err = tx.Where("id = ?", deleteUserID).Delete(&models.User{}).Error
	if err != nil {
        return err
    }

	err = tx.Commit().Error
    if err != nil {
        log.Printf("[Repository DeleteUser DATABASE ERROR] Failed to commit transaction : %v", err)
        return err
    }

	return nil
}
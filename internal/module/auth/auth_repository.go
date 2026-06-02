package auth

import (
	"errors"
	"gin-quickstart/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db  *gorm.DB
}

func NewAuthRepository (db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// Search
func (repo *AuthRepository) CheckSuperAdminExists(rolesID uuid.UUID) (bool, error) {
	if rolesID == uuid.Nil {
		return false, errors.New("roles id is required.")
	}

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

func (repo *AuthRepository) CheckRefreshTokenValid(hashedRefreshToken string) (bool, error) {
	if hashedRefreshToken == "" {
		return false, errors.New("hash refresh token is required.")
	}

	var count int64
	now := time.Now()
	
	err := repo.db.Model(&models.RefreshToken{}).
		Where("token_hash = ? AND is_revoked = ? AND expires_at > ?", hashedRefreshToken, false, now).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	// คืนค่ากลับไป: ถ้า count > 0 แปลว่าโทเคนนี้ "ยังมีอายุและใช้งานได้จริง" (ส่ง true)
	// ถ้าหาไม่เจอหรือเงื่อนไขไม่ผ่าน count จะเป็น 0 (ส่ง false)
	return count > 0, nil
}

func (repo *AuthRepository) FineUserByUserName(userName string) (*models.User, error){
	if userName == "" {
		return nil, errors.New("user name is required.")
	}

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

func (repo *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required.")
	}

	var user models.User

	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *AuthRepository) FindValidResetToken(token string) (*models.ResetPassword, error) {
	if token == "" {
		return nil, errors.New("token is required.")
	}

	var reset models.ResetPassword

	err := repo.db.Where("token = ? AND is_used = ? AND expired_at > ?", token, false, time.Now()).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// Create
func (repo *AuthRepository) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("data user is required.")
	}

	err := repo.db.Create(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateRefreshTokenRecord(refreshToken *models.RefreshToken) error {
	if refreshToken == nil {
		return errors.New("data refresh token is required.")
	}

	err := repo.db.Create(refreshToken).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateResetPassword(reset *models.ResetPassword) error{
	if reset == nil {
		return errors.New("data reset password is required.")
	}

	err := repo.db.Create(reset).Error

	if err != nil {
		return err
	}

	return nil
}

// Update
func (repo *AuthRepository) RevokeRefreshToken(userID uuid.UUID, hashedToken string) error {

	query := repo.db.Model(&models.RefreshToken{}).Where("is_revoked = ?", false)

	if hashedToken != "" {
		query = query.Where("token_hash = ?", hashedToken)
	} else if userID != uuid.Nil {
		query = query.Where("user_id = ? AND expires_at > ?", userID, time.Now())
	} else {
		return errors.New("either hashed token or user id is required to revoke sessions")
	}

	return query.Update("is_revoked", true).Error
}

func (repo *AuthRepository) UpdatePasswordAndRevokeToken(userID uuid.UUID, hashedPwd string, resetID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("user id is required.")
	}

	if hashedPwd == "" {
		return errors.New("hashed password is required.")
	}

	if resetID == uuid.Nil {
		return errors.New("reset id is required.")
	}
	
	tx := repo.db.Begin()

	// 1. อัปเดตรหัสผ่านใหม่ให้ User
	if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPwd).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. อัปเดตตั๋วใบนี้ว่า "ใช้แล้ว" (is_used = true)
	if err := tx.Model(&models.ResetPassword{}).Where("id = ?", resetID).Update("is_used", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


package auth

import (
	"errors"
	"log"
	"pos-system-backend/internal/models"
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

// Get
func (repo *AuthRepository) CheckSystemAdminExists(systemRole string) (bool, error) {
	var user models.User
	
	err := repo.db.Select("id").Where("system_role = ?", systemRole).First(&user).Error

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

func (repo *AuthRepository) CheckPermission(userID uuid.UUID, storeID uuid.UUID) (*models.UserStore, error) {
	var userStore models.UserStore

	err := repo.db.Preload("Role").Where("user_id = ? AND store_id = ?", userID, storeID).First(&userStore).Error

	if err != nil {
        return nil, err // ส่ง error กลับไป (เช่น gorm.ErrRecordNotFound)
    }

    return &userStore, nil
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


func (repo *AuthRepository) FindUserByEmail(email string, findType string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required.")
	}

	var user models.User

	if(findType == "LOGIN"){
		// ถ้าหาไม่เจอจะคืนเป็น error ถ้าใช้ First
		err := repo.db.Preload("UserStores.Store").Preload("UserStores.Role").Where("email = ? AND is_active = ?", email, true).First(&user).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := repo.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error
		if err != nil {
			return nil, err
		}
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
func (repo *AuthRepository) CreateUserSystemAdmin(userData *models.User) error {
	if userData == nil {
		return errors.New("data user is required.")
	}

	err := repo.db.Create(userData).Error

	if err != nil {
		log.Printf("[Repository CreateUserSystemAdmin DATABASE ERROR] Failed to insert system admin (%s): %v", userData.Email, err)
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateUser(userData *models.User, userStoreData *models.UserStore) error {
	if userData == nil || userStoreData == nil {
        return errors.New("user data and store permission data are required")
    }

	tx := repo.db.Begin()
    if tx.Error != nil {
		log.Printf("[Repository CreateUser DATABASE ERROR] Failed to start transaction for user %s: %v", userData.Email, tx.Error)
        return tx.Error
    }

	err := tx.Create(userData).Error
	if  err != nil {
		log.Printf("[Repository CreateUser DATABASE ERROR] Step 1 failed -> Inserting user record (%s): %v. Rolling back...", userData.Email, err)
        tx.Rollback() // 🚨 พังตรงนี้ให้ยกเลิกทั้งหมดทันที
        return err
    }

	userStoreData.UserID = userData.ID

	err = tx.Create(userStoreData).Error
	if err != nil {
		log.Printf("[Repository CreateUser DATABASE ERROR] Step 2 failed -> Inserting user_store permission for user %s (StoreID: %s, RoleID: %s): %v. Rolling back...", userData.Email, userStoreData.StoreID.String(), userStoreData.RoleID.String(), err)
        tx.Rollback() // 🚨 หากผูกสิทธิ์พัง ให้กดยกเลิกการสร้าง User ก่อนหน้าไปด้วยเพื่อความปลอดภัย
        return err
    }

	err = tx.Commit().Error
    if err != nil {
        log.Printf("[Repository CreateUser DATABASE ERROR] Transaction commit failed for user %s: %v", userData.Email, err)
        return err
    }

	return nil
}

func (repo *AuthRepository) CreateRefreshTokenRecord(refreshTokenData *models.RefreshToken) error {
	if refreshTokenData == nil {
		return errors.New("data refresh token is required.")
	}

	err := repo.db.Create(refreshTokenData).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateResetPassword(reset *models.ResetPassword) error {
	if reset == nil {
		return errors.New("data reset password is required.")
	}

	err := repo.db.Create(reset).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateLogEmail(logEmail *models.LogSendEmail) error {
	if logEmail == nil {
        return errors.New("log email data is required")
    }

	err := repo.db.Create(logEmail).Error
    if err != nil {
        log.Printf("[Repository CreateLogEmail DATABASE ERROR] Failed to insert email log for %s: %v", logEmail.Recipient, err)
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

	updates := map[string]interface{}{
        "is_revoked": true,
        "updated_at": time.Now(),
		"updated_by": userID,
    }

	err := query.Updates(updates).Error

	if err != nil {
        log.Printf("[Repository RevokeRefreshToken DATABASE ERROR] Failed to revoke session (User: %s, HasToken: %t): %v", 
            userID.String(), hashedToken != "", err)
        return err
    }

    return nil
}

func (repo *AuthRepository) UpdatePasswordAndRevokeToken(userID uuid.UUID, hashedPwd string, resetID uuid.UUID) error {
	if userID == uuid.Nil || hashedPwd == "" || resetID == uuid.Nil {
        return errors.New("missing required arguments for updating password and revoking token")
    }
    
    tx := repo.db.Begin()
    if tx.Error != nil {
        log.Printf("[Repository ResetPassword DATABASE ERROR] Failed to start tx: %v", tx.Error)
        return tx.Error
    }

    // ขั้นตอนที่ 1: อัปเดตรหัสผ่านใหม่ให้กับตาราง User หลัก
    resultUser := tx.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPwd)
    if resultUser.Error != nil {
        log.Printf("[Repository ResetPassword DATABASE ERROR] Step 1 failed -> Updating user password: %v. Rolling back...", resultUser.Error)
        tx.Rollback()
        return resultUser.Error
    }
    if resultUser.RowsAffected == 0 {
        log.Printf("[Repository ResetPassword WARN] Step 1 failed -> User %s not found. Rolling back...", userID.String())
        tx.Rollback()
        return errors.New("user not found, password update failed")
    }

    // ขั้นตอนที่ 2: อัปเดตสถานะตั๋วใบนี้ว่าโดนใช้งานไปเรียบร้อยแล้ว (Prevent Replay Attack)
    resultReset := tx.Model(&models.ResetPassword{}).Where("id = ?", resetID).Update("is_used", true)
    if resultReset.Error != nil {
        log.Printf("[Repository ResetPassword DATABASE ERROR] Step 2 failed -> Updating reset token state: %v. Rolling back...", resultReset.Error)
        tx.Rollback()
        return resultReset.Error
    }
    if resultReset.RowsAffected == 0 {
        log.Printf("[Repository ResetPassword WARN] Step 2 failed -> Reset token %s is already invalid. Rolling back...", resetID.String())
        tx.Rollback()
        return errors.New("reset token is invalid or has already been used")
    }

    // ขั้นตอนที่ 3: สั่งกวาดล้างและยกเลิก Refresh Token (Session) ทั้งหมดของยูสเซอร์คนนี้ที่เคยเปิดทิ้งไว้เครื่องอื่น
    // 💡 รันผ่านตัวแปร tx เพื่อความปลอดภัยของธุรกรรม
    err := tx.Model(&models.RefreshToken{}).
        Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now()).
        Update("is_revoked", true).Error
        
    if err != nil {
        log.Printf("[Repository ResetPassword DATABASE ERROR] Step 3 failed -> Revoking active refresh tokens: %v. Rolling back...", err)
        tx.Rollback()
        return err
    }

    // 🎉 ทุกด่านสมบูรณ์แบบ ทำการเซฟลงดิสก์ถาวร
    err = tx.Commit().Error
    if err != nil {
        log.Printf("[Repository ResetPassword DATABASE ERROR] Transaction commit crashed: %v", err)
        return err
    }

    return nil
}

func (repo *AuthRepository) UpdateLogEmailStatus(logID uuid.UUID, status string, errMsg *string, userID uuid.UUID) error {
	now := time.Now()

	updates := map[string]interface{}{
        "status":     status,
        "updated_at": &now,
		"updated_by": userID,
    }

	if errMsg != nil {
        updates["error_message"] = *errMsg
    }

	err := repo.db.Model(&models.LogSendEmail{}).Where("id = ?", logID).Updates(updates).Error
    if err != nil {
        log.Printf("[Repository UpdateLogEmailStatus DATABASE ERROR] Failed to update log ID %s to %s: %v", logID.String(), status, err)
        return err
    }

    return nil
}


package store

import (
	"errors"
	"fmt"
	"log"
	"pos-system-backend/internal/models"
	storeDto "pos-system-backend/internal/module/store/dto"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}


func (repo *StoreRepository) CheckIsOwner(userID uuid.UUID) (bool, error) {
    var count int64
    err := repo.db.Table("user_stores").
        Joins("JOIN roles ON roles.id = user_stores.role_id").
        Where("user_stores.user_id = ? AND roles.role_name = ? AND user_stores.is_active = ? AND user_stores.deleted_at IS NULL", userID, "OWNER", true).
        Count(&count).Error
        
    return count > 0, err
}

func (repo *StoreRepository) GetUserStoresList(userID uuid.UUID, systemRole string, isOwner bool, data *storeDto.GetStoreRequest) (interface{}, int64, error) {
	if data.Page <= 0 {
        data.Page = 1 
    }
    if data.Limit <= 0 {
        data.Limit = 10 
    } else if data.Limit > 50 {
        data.Limit = 50 
    }
    offset := (data.Page - 1) * data.Limit

	// เคสที่ 1: เป็น SYSTEM_ADMIN หรือ USER ที่มีสิทธิ์ OWNER (เห็นทุกสาขาในระบบ)
    if systemRole == "SYSTEM_ADMIN" || isOwner {
        var stores []models.Store
        var total int64

        query := repo.db.Model(&models.Store{}).Where("deleted_at IS NULL")

        if data.Search != "" {
            searchPattern := fmt.Sprintf("%%%s%%", data.Search)
            query = query.Where("store_code LIKE ? OR store_name LIKE ?", searchPattern, searchPattern)
        }

        if err := query.Count(&total).Error; err != nil {
            return nil, 0, err
        }

        if err := query.Limit(data.Limit).Offset(offset).Find(&stores).Error; err != nil {
            return nil, 0, err
        }

        return stores, total, nil
    }

	// เคสที่ 2: เป็น USER ทั่วไป (MANAGER, STAFF) -> ดึงเฉพาะสาขาที่ตัวเองมีสิทธิ์ผูกอยู่
    var userStores []models.UserStore
    var total int64

    query := repo.db.Model(&models.UserStore{}).
        Preload("Store").
        Preload("Role").
        Joins("JOIN stores ON stores.id = user_stores.store_id").    
    	Where("user_stores.user_id = ? " +
          "AND user_stores.deleted_at IS NULL " +
          "AND stores.deleted_at IS NULL " +
          "AND user_stores.is_active = ? " + // เช็คว่าสิทธิ์พนักงานยังใช้งานได้อยู่ไหม
          "AND stores.is_active = ?",        // เช็คว่าสาขานี้ยังเปิดให้บริการอยู่ไหม
          userID, true, true)

    if data.Search != "" {
        searchPattern := fmt.Sprintf("%%%s%%", data.Search)
        query = query.Where("stores.store_code LIKE ? OR stores.store_name LIKE ?", searchPattern, searchPattern)
    }

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.Limit(data.Limit).Offset(offset).Find(&userStores).Error; err != nil {
        return nil, 0, err
    }

    return userStores, total, nil
}

func (repo *StoreRepository) GetLastStoreCode(tx *gorm.DB) (string, error) {
    var lastCode string
    err := tx.Model(&models.Store{}).
        Unscoped().  //เอาค่าทั้งหมดไม่สนว่ามี deleteAt ไหม
        Clauses(clause.Locking{Strength: "UPDATE"}). // lock แถวล่าสุดไว้เพื่อป้องกันปัญหา Race Condition แย่งเลขกัน
        Select("store_code").
        Order("store_code DESC"). 
        Limit(1).
        Pluck("store_code", &lastCode).Error //การแกะเอาเฉพาะค่าในคอลัมน์ (Column) ที่เราเลือก ออกมาใส่ไว้ในตัวแปรเดี่ยว ๆ

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return "", nil
    }
    if err != nil {
        return "", err
    }

    return lastCode, nil
}

func (repo *StoreRepository) CreateStore(store *models.Store, isBindOwner bool, userID uuid.UUID, ownerRoleID uuid.UUID) error {
	tx := repo.db.Begin()
    if tx.Error != nil {
		log.Printf("[Repository CreateStore DATABASE ERROR] Failed to start transaction : %v", tx.Error)
        return tx.Error
    }

    defer tx.Rollback()

    lastCode, err := repo.GetLastStoreCode(tx)
    if err != nil {
        return err 
    }

    nextNumber := 1
    if lastCode != "" {
        parts := strings.Split(lastCode, "-")
        if len(parts) == 2 {
            currentNum, _ := strconv.Atoi(parts[1])
            nextNumber = currentNum + 1
        }
    }
    store.StoreCode = fmt.Sprintf("ST-%04d", nextNumber)

    err = tx.Create(store).Error;
    if  err != nil {
        log.Printf("[Repo ERROR][CreateStore] Failed to create store and its associations: %v", err)
        return err
    }

    if isBindOwner {
        userStore := models.UserStore{
            UserID:  userID,
            StoreID: store.ID, // GORM เติม ID ร้านค้าใหม่ลงมาในตัวแปรนี้ให้เรียบร้อยแล้วหลัง tx.Create ข้างบน
            RoleID:  ownerRoleID,
        }

        err := tx.Create(&userStore).Error;
        if  err != nil {
            log.Printf("[Repo ERROR][CreateStore] Failed to bind user_store role: %v", err)
            return err
        }
    }

    err = tx.Commit().Error
    if err != nil {
        log.Printf("[Repository CreateStore DATABASE ERROR] Transaction commit failed : %v", err)
        return err
    }

    return  nil
}

func (repo *StoreRepository) UpdateStore(store *models.Store, storeAddress *models.StoreAddress, storeID uuid.UUID) error {
    tx := repo.db.Begin()
    if tx.Error != nil {
		log.Printf("[Repository UpdateStore DATABASE ERROR] Failed to start transaction : %v", tx.Error)
        return tx.Error
    }

    defer tx.Rollback()

    err := tx.Model(&models.Store{}).Where("id = ?", storeID).Updates(store).Error
    if err != nil {
        log.Printf("[Repository UpdateStore DATABASE ERROR] Failed to update store data : %v", err)
        return err
    }

    err = tx.Model(&models.StoreAddress{}).Where("store_id = ?", storeID).Updates(storeAddress).Error
    if err != nil {
        log.Printf("[Repository UpdateStore DATABASE ERROR] Failed to update store address : %v", err)
        return err
    }

    err = tx.Commit().Error
    if err != nil {
        log.Printf("[Repository UpdateStore DATABASE ERROR] Failed to commit transaction : %v", err)
        return err
    }

    return nil
}

func (repo *StoreRepository) UpdateStoreStatus(storeID uuid.UUID, store *models.Store) error {
    err := repo.db.Model(&models.Store{}).
        Where("id = ?", storeID).
        Updates(store).Error

    if err != nil {
        log.Printf("[Repository UpdateStoreFields ERROR] : %v", err)
        return err
    }
    return nil
}
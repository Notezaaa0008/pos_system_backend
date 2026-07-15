package store

import (
	"log"
	"pos-system-backend/internal/models"
	storeDto "pos-system-backend/internal/module/store/dto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type storeRepositoryInterface interface {
	CheckIsOwner(userID uuid.UUID) (bool, error)
	GetUserStoresList(userID uuid.UUID, systemRole string, isOwner bool, data *storeDto.GetStoreRequest) (interface{}, int64, error)
    CreateStore(storeData *models.Store, isBindOwner bool, userID uuid.UUID, ownerRoleID uuid.UUID) error
    UpdateStore(store *models.Store, storeAddress *models.StoreAddress, storeID uuid.UUID) error
    UpdateStoreStatus(storeID uuid.UUID, store *models.Store) error
}

type StoreService struct {
	repo storeRepositoryInterface
}

func NewStoreService(repo storeRepositoryInterface) *StoreService {
	return &StoreService{repo: repo}
}

func (service *StoreService) GetStoreListService(userID uuid.UUID, systemRole string, req *storeDto.GetStoreRequest) ([]gin.H, int64, error){
	//เช็คสิทธิ์ก่อนว่าผู้ใช้คนนี้เป็น OWNER ไหม
    isOwner := false
    if systemRole == "USER" {
        var err error
        isOwner, err = service.repo.CheckIsOwner(userID)
        if err != nil {
            return nil, 0, err
        }
    }
	
	//เรียก Repository เพื่อดึงข้อมูลตามเงื่อนไขสิทธิ์
	rawData, total, err := service.repo.GetUserStoresList(userID, systemRole, isOwner, req)
    if err != nil {
        return nil, 0, err
    }

	// []gin.H คือ map[string]interface
	var result []gin.H

	// 3. จัดการแปลงข้อมูล (Type Assertion) เพื่อพ่น JSON รูปแบบเดียวกันออกไป
    if systemRole == "SYSTEM_ADMIN" || isOwner {
        // แตกข้อมูลจากกรณีดึงตาราง Store ตรงๆ
        stores := rawData.([]models.Store)
        for _, s := range stores {
            roleName := "OWNER"
            if systemRole == "SYSTEM_ADMIN" {
                roleName = "SYSTEM_ADMIN"
            }
            
            result = append(result, gin.H{
                "store_id":   s.ID,
                "store_code": s.StoreCode,
                "store_name": s.StoreName,
                "role_name":  roleName, // บังคับสิทธิ์ให้แสดงเป็น OWNER หรือ ADMIN บนหน้าจอไปเลย
            })
        }
    } else {
        // แตกข้อมูลจากกรณีพนักงานทั่วไป (ตาราง UserStore)
        userStores := rawData.([]models.UserStore)
        for _, us := range userStores {
            result = append(result, gin.H{
                "store_id":   us.StoreID,
                "store_code": us.Store.StoreCode,
                "store_name": us.Store.StoreName,
                "role_name":  us.Role.RoleName,
            })
        }
    }

    return result, total, nil
}

func (service *StoreService) CreateStoreService(userID uuid.UUID, rolrID uuid.UUID, systemRole string, req *storeDto.CreateStoreRequest) (error) {
    isBindOwner := false
    var ownerRoleID uuid.UUID
    if systemRole == "USER" {
        isBindOwner = true
        ownerRoleID = rolrID
    }

    // ปั้นข้อมูลที่อยู่ (models.StoreAddress)
    address := models.StoreAddress{
        ProvinceID:    req.ProvinceID,
        DistrictID:    req.DistrictID,
        SubdistrictID: req.SubdistrictID,
        PostCodeID:    req.PostCodeID,
        IsActive:      true,
        CreatedBy:     userID,
    }

    // มัดรวมทุกอย่างเข้าก้อน Store หลัก
    storeData := models.Store{
        StoreName:      req.StoreName,
        BranchName:     req.BranchName,
        Description:    req.Description,
        PrimaryPhone:   req.PrimaryPhone,
        SecondaryPhone: req.SecondaryPhone,
        LineID:         req.LineID, 
        IsActive:       true,
        CreatedBy:      userID,
        StoreAddress:   &address, // 🔥 ใส่ความสัมพันธ์ลูกลงไปตรง ๆ
    }

    err := service.repo.CreateStore(&storeData, isBindOwner, userID, ownerRoleID)
    if err != nil {
        return err
    }

    return nil
}

func (service *StoreService) UpdateStoreService(userID uuid.UUID, storeID uuid.UUID, req *storeDto.UpdateStoreRequest) error {
    now := time.Now()
    updateStore := models.Store{
        StoreName:      req.StoreName,
        BranchName:     req.BranchName,
        Description:    req.Description,
        PrimaryPhone:   req.PrimaryPhone,
        SecondaryPhone: req.SecondaryPhone,
        LineID:         req.LineID,
        UpdatedAt:      &now,
        UpdatedBy:      &userID,
    }

    updateStoreAddress := models.StoreAddress{
        ProvinceID:     req.ProvinceID,
        DistrictID:     req.DistrictID,
        SubdistrictID:  req.SubdistrictID,
        PostCodeID:     req.PostCodeID,
        UpdatedAt:      &now,
        UpdatedBy:      &userID,
    }

    err := service.repo.UpdateStore(&updateStore, &updateStoreAddress, storeID)
    if err != nil {
        log.Printf("[Service UpdateStoreService ERROR] Failed to update Store Error: %v", err)
        return err
    }

    return nil
}

func (service *StoreService) UpdateStoreStatusService(storeID uuid.UUID, isActive bool, userID uuid.UUID) error {
    now := time.Now()
    updateStore := models.Store{
        IsActive:  isActive,
        UpdatedAt: &now,
        UpdatedBy: &userID,
    }

    err := service.repo.UpdateStoreStatus(storeID, &updateStore)
    if err != nil {
        log.Printf("[Service UpdateStoreStatusService ERROR] Failed to update Store Status Error: %v", err)
        return err
    }
    
    return nil
}
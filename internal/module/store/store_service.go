package store

import (
	"log"
	"pos-system-backend/internal/models"
	storeDto "pos-system-backend/internal/module/store/dto"
	"time"

	"github.com/google/uuid"
)

type storeRepositoryInterface interface {
	CheckIsOwner(userID uuid.UUID) (bool, error)
	GetUserStoresList(userID uuid.UUID, systemRole string, isOwner bool, data *storeDto.GetStoreRequest) (interface{}, int64, error)
    CreateStore(storeData *models.Store, isBindOwner bool, userID uuid.UUID, ownerRoleID uuid.UUID) error
    UpdateStore(store *models.Store, storeAddress *models.StoreAddress, storeID uuid.UUID) error
    UpdateStoreStatus(storeID uuid.UUID, store *models.Store) error
    DeleteStore(storeID uuid.UUID, storeData *models.Store, userStoreData *models.UserStore) error
}

type StoreService struct {
	repo storeRepositoryInterface
}

func NewStoreService(repo storeRepositoryInterface) *StoreService {
	return &StoreService{repo: repo}
}

func (service *StoreService) GetStoreListService(userID uuid.UUID, systemRole string, req *storeDto.GetStoreRequest) ([]storeDto.GetStoreListResponse, int64, error){
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
	result := []storeDto.GetStoreListResponse{}

	// จัดการแปลงข้อมูล (Type Assertion) เพื่อพ่น JSON รูปแบบเดียวกันออกไป
    if systemRole == "SYSTEM_ADMIN" || isOwner {
        // แตกข้อมูลจากกรณีดึงตาราง Store ตรงๆ
        stores := rawData.([]models.Store)
        for _, s := range stores {
            roleName := "OWNER"
            if systemRole == "SYSTEM_ADMIN" {
                roleName = "SYSTEM_ADMIN"
            }
            
            formatted := storeDto.GetStoreListResponse{
                ID: s.ID,
                StoreCode: s.StoreCode,
                StoreName: s.StoreName,
                RoleName: roleName,
            }

            result = append(result, formatted)
        }
    } else {
        // แตกข้อมูลจากกรณีพนักงานทั่วไป (ตาราง UserStore)
        userStores := rawData.([]models.UserStore)
        for _, us := range userStores {
            formatted := storeDto.GetStoreListResponse{
                ID: us.StoreID,
                StoreCode: us.Store.StoreCode,
                StoreName: us.Store.StoreName,
                RoleName: us.Role.RoleName,
            }

            result = append(result, formatted)
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
        Address:       req.Address,  
        ProvinceID:    req.ProvinceID,
        DistrictID:    req.DistrictID,
        SubdistrictID: req.SubdistrictID,
        PostcodeID:    req.PostcodeID,
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
        Address:        req.Address,
        ProvinceID:     req.ProvinceID,
        DistrictID:     req.DistrictID,
        SubdistrictID:  req.SubdistrictID,
        PostcodeID:     req.PostcodeID,
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

func (service *StoreService) DeleteStoreService(storeID uuid.UUID, userID uuid.UUID) error {
    // 🎯 1. ปั้นข้อมูลสำหรับอัปเดตตาราง Store
    storeData := models.Store{
        IsActive:  false,
        DeletedBy: &userID,
    }

    // 🎯 2. ปั้นข้อมูลสำหรับอัปเดตตาราง UserStore
    userStoreData := models.UserStore{
        IsActive:  false,
        DeletedBy: &userID,
    }

    err := service.repo.DeleteStore(storeID, &storeData, &userStoreData)
    if err != nil {
        log.Printf("[Service DeleteStoreService ERROR] Failed to delete store: %v", err)
        return err
    }

    return nil
}
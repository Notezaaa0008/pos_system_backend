package store

import (
	"log"
	"net/http"
	storeDto "pos-system-backend/internal/module/store/dto"
	"pos-system-backend/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type storeServiceInterface interface {
	GetStoreListService(userID uuid.UUID, systemRole string, req *storeDto.GetStoreRequest) ([]gin.H, int64, error)
    CreateStoreService(userID uuid.UUID, rolrID uuid.UUID, systemRole string, req *storeDto.CreateStoreRequest) (error)
    UpdateStoreService(userID uuid.UUID, storeID uuid.UUID, req *storeDto.UpdateStoreRequest) error
    UpdateStoreStatusService(storeID uuid.UUID, isActive bool, userID uuid.UUID) error
}

type StoreController struct {
	service storeServiceInterface
}

func NewStoreController(service storeServiceInterface) *StoreController {
	return &StoreController{service: service}
}

func (storeCtrl *StoreController) GetStoreListController(c *gin.Context) {
	userID, err := utils.GetFromCtx(c, "userID")
	if err != nil {
        log.Printf("[STORE][GET_USER_STORES][CTX_ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: user identity not found",
        })
        return
    }

	systemRole := c.GetString("systemRole")
	if systemRole == "" {
        log.Printf("[STORE][GET_USER_STORES][CTX_ERROR] systemRole missing or empty in context")
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: system role not found",
        })
        return
    }

	var req storeDto.GetStoreRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_USER_STORES][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

	stores, total, err := storeCtrl.service.GetStoreListService(userID, systemRole, &req)
	if err != nil {
		log.Printf("[STORE][GET_USER_STORES][SERVICE_ERROR] userID=%v error=%v", userID, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Internal server error: Failed to retrieve user stores",
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    stores,
		"total":   total,       
		"page":    req.Page,  
		"limit":   req.Limit, 
	})

}

func (storeCtrl *StoreController) CreateStoreController(c *gin.Context) {
	userID, err := utils.GetFromCtx(c, "userID")
	if err != nil {
        log.Printf("[STORE][CREATE_STORES][CTX_ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: user identity not found",
        })
        return
    }

	systemRole := c.GetString("systemRole")
	if systemRole == "" {
        log.Printf("[STORE][GET_USER_STORES][CTX_ERROR] systemRole missing or empty in context")
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: system role not found",
        })
        return
    }

    roleID := uuid.Nil 
    if systemRole == "USER" {
        storeRoleID, err := utils.GetFromCtx(c, "storeRoleID")
        if err != nil {
            log.Printf("[STORE][GET_USER_STORES][CTX_ERROR] Failed to get storeRoleID from context: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{
                "status":  "error",
                "message": "Unauthorized: store role identity not found",
            })
            return
        }

        roleID = storeRoleID
    }  

    var req storeDto.CreateStoreRequest
    err = c.ShouldBind(&req) // สำหรับ Form-data
    if err != nil {
		log.Printf("[CreateStore WARN] Validation failed for user input form-data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid request body format",
            "error":   err.Error(),
        })
        return
    }

    storeName, isBlank := utils.IsBlank(req.StoreName)
    if isBlank {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid store name",
        })
        return
    }
    req.StoreName = storeName

    branchName, isBlank := utils.IsBlank(req.BranchName)
    if isBlank {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid branch name",
        })
        return
    }
    req.BranchName = branchName

    req.PrimaryPhone = utils.CleanInputPhoneNumber(req.PrimaryPhone)
    if req.PrimaryPhone == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "primary_phone cannot be blank"})
        return
    }

    if req.SecondaryPhone != nil {
    cleaned := utils.CleanInputPhoneNumber(*req.SecondaryPhone)
        if cleaned == "" {
            req.SecondaryPhone = nil
        } else {
            req.SecondaryPhone = &cleaned
        }
    }

    if req.LineID != nil {
        trimmedLineID := strings.TrimSpace(*req.LineID)
        if trimmedLineID == "" {
            req.LineID = nil
        } else {
            req.LineID = &trimmedLineID
        }
    }

    err = storeCtrl.service.CreateStoreService(userID, roleID, systemRole, &req)
    if err != nil {
        log.Printf("[CreateStore Controller ERROR] Service failed to execute create process: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to create new store profile",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "status":  "success",
        "message": "The new store and address have been successfully registered",
    })
}

func (storeCtrl *StoreController) UpdateStoreController(c *gin.Context) {
    userID, err := utils.GetFromCtx(c, "userID")
	if err != nil {
        log.Printf("[STORE][UPDATE_STORES][CTX_ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: user identity not found",
        })
        return
    }

    storeID, err := utils.GetFromCtx(c, "storeID")
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
		log.Printf("[STORE][UPDATE_STORES][CTX_ERROR] Failed to get storeID from context: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing store identity for this operation"})
        return
    }

    var req storeDto.UpdateStoreRequest
    err = c.ShouldBind(&req) // สำหรับ Form-data
    if err != nil {
		log.Printf("[UpdateStore WARN] Validation failed for user input form-data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid request body format",
            "error":   err.Error(),
        })
        return
    }

    err = storeCtrl.service.UpdateStoreService(userID, storeID, &req)
    if  err != nil {
        log.Printf("[UpdateStore Controller ERROR] Service failed to execute update store process: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update store details"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Store and address details have been successfully updated",
    })
}

func (storeCtrl *StoreController) UpdateStoreStatusController(c *gin.Context) {
    userID, err := utils.GetFromCtx(c, "userID")
	if err != nil {
        log.Printf("[STORE][UPDATE_STORES_STATUS][CTX_ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: user identity not found",
        })
        return
    }

    storeID, err := utils.GetFromCtx(c, "storeID")
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
		log.Printf("[STORE][UPDATE_STORES_STATUS][CTX_ERROR] Failed to get storeID from context: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing store identity for this operation"})
        return
    }

    var req storeDto.UpdateStorestatus
    err = c.ShouldBindJSON(&req) 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid request body"})
        return
    }

    err = storeCtrl.service.UpdateStoreStatusService(storeID, *req.IsActive, userID)
    if err != nil {
        log.Printf("[UpdateStoreStatus Controller ERROR] Service failed to execute update store status process: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to update status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Store status updated successfully",
    })
}
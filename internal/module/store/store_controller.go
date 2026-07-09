package store

import (
	"log"
	"net/http"
	storeDto "pos-system-backend/internal/module/store/dto"
	"pos-system-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type storeServiceInterface interface {
	GetUserStoreService(userID uuid.UUID, req *storeDto.GetUserStoreRequest) ([]gin.H, int64, error)
}

type StoreController struct {
	service storeServiceInterface
}

func NewStoreController(service storeServiceInterface) *StoreController {
	return &StoreController{service: service}
}

func (storeCtrl *StoreController) GetUserStoreController(c *gin.Context) {
	userID, err := utils.GetFromCtx(c, "userID")
	
	if err != nil {
        log.Printf("[STORE][GET_USER_STORES][CTX_ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: user identity not found",
        })
        return
    }

	var req storeDto.GetUserStoreRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_USER_STORES][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

	stores, total, err := storeCtrl.service.GetUserStoreService(userID, &req)
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
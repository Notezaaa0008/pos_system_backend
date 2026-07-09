package master

import (
	"net/http"
	masterDto "pos-system-backend/internal/module/master/dto"

	"github.com/gin-gonic/gin"
)

type masterServiceInterface interface {
	GetAllPrefixService() ([]masterDto.GetAllPrefixResponse, error)
}

type MasterController struct {
	service masterServiceInterface
}

func NewMasterController(service masterServiceInterface) *MasterController {
	return &MasterController{service: service}
}

func (MasterCtrl *MasterController) GetAllPrefixController(c *gin.Context) {
	perfixs, err := MasterCtrl.service.GetAllPrefixService()

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch roles data",
            "error":   err.Error(), // ถ้าขึ้น Production อาจจะซ่อนไว้เพื่อความปลอดภัย
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active roles",
        "data":  perfixs, // หน้าบ้านจะแกะข้อมูลจาก Field "data" อันนี้ไปใช้ต่อ
    })
}
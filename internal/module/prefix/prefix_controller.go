package prefix

import (
	"net/http"
	prefixDto "pos-system-backend/internal/module/prefix/dto"

	"github.com/gin-gonic/gin"
)

type prefixServiceInterface interface {
	GetAllPrefixService() ([]prefixDto.GetAllPrefixResponse, error)
}

type PrefixController struct {
	service prefixServiceInterface
}

func NewPrefixController(service prefixServiceInterface) *PrefixController {
	return &PrefixController{service: service}
}

func (prefixCtrl *PrefixController) GetAllPrefixController(c *gin.Context) {
	perfixs, err := prefixCtrl.service.GetAllPrefixService()

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
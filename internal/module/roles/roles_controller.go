package roles

import (
	"errors"
	roledto "gin-quickstart/internal/module/roles/dto"
	"gin-quickstart/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type rolesServiceInterface interface {
	GetAllRolesService() ([]roledto.GetAllRoleResponse, error)
    CreateRoleService(req *roledto.CreateRoleRequest, userId uuid.UUID) error
	UpadateRoleService(req *roledto.UpdateRoleRequest, userId uuid.UUID) error
}

type RolesController struct {
	service rolesServiceInterface
}

func NewRoleController (service rolesServiceInterface) *RolesController{
	return &RolesController{service: service}
}

func (roleClrt *RolesController) GetAllRolesController(c *gin.Context) {
	roles, err := roleClrt.service.GetAllRolesService()

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
        "data":    roles, // หน้าบ้านจะแกะข้อมูลจาก Field "data" อันนี้ไปใช้ต่อ
    })

}

func (roleClrt *RolesController) CreateRoleController(c *gin.Context) {
	userId, err := utils.GetUserIDFromCtx(c)
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

    var req roledto.CreateRoleRequest

    err = c.ShouldBindJSON(&req)

    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

    err = roleClrt.service.CreateRoleService(&req, userId)

    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Role created successfully.",
	})
}

func (roleClrt *RolesController) UpdateRoleController(c *gin.Context) {
	userId, err := utils.GetUserIDFromCtx(c)
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	var req roledto.UpdateRoleRequest

	err = c.ShouldBindJSON(&req)

    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = roleClrt.service.UpadateRoleService(&req, userId)

	if err != nil {
		var appErr *utils.AppError

		if errors.As(err, &appErr) {
            c.JSON(appErr.StatusCode, gin.H{
                "status":  "error",
                "message": appErr.Message,
            })
            return
        }
		
		c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Internal server error. Something went wrong.",
        })
		return
	}

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Role updated successfully",
    })
}
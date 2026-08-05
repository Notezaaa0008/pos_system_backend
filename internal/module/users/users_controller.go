package users

import (
	"log"
	"net/http"
	usersdto "pos-system-backend/internal/module/users/dto"
	"pos-system-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userServiceInterface interface {
	GetProfileService(userID uuid.UUID) (*usersdto.GetProfileResponse, error)
	UpdateProfileService(userID uuid.UUID, req *usersdto.UpdateProfileRequest) error
	DeleteUserService(deleteUserID uuid.UUID, userID uuid.UUID) error
}

type UsersController struct {
	service userServiceInterface
}

func NewUsersController (service userServiceInterface) *UsersController{
	return &UsersController{service: service}
}

func (userClrt *UsersController) GetAllUserController(c *gin.Context) {
	
}

func (userClrt *UsersController) GetProfileController(c *gin.Context) {
	userID, err := utils.GetFromCtx(c, "userID")
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
		log.Printf("[GET PROFILE ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	profile, err:= userClrt.service.GetProfileService(userID)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch profile data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully get profile.",
        "data":  profile,
    })
}

func (userClrt *UsersController) UpdateProfileController(c *gin.Context) {
	userID, err := utils.GetFromCtx(c, "userID")
	
	if err != nil {
		log.Printf("[UPDATE PROFILE ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	var req usersdto.UpdateProfileRequest
	err = c.ShouldBind(&req) // สำหรับ Form-data
    if err != nil {
		log.Printf("[Update profile WARN] Validation failed for user input form-data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid request body format",
        })
        return
    }

	err = userClrt.service.UpdateProfileService(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to update profile data",
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully update profile.",
    })
}

func (userClrt *UsersController) DeleteUserController(c *gin.Context) {
	deleteUserIDStr := c.Param("id")

	deleteUserID, err := uuid.Parse(deleteUserIDStr)
	if err != nil {
		log.Printf("[USERS][DELETE_USER][INVALID_UUID] param: %s, err: %v", deleteUserIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload or missing required parameters",
		})
		return
	}

	userID, err := utils.GetFromCtx(c, "userID")
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
		log.Printf("[DELETE USER ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	err = userClrt.service.DeleteUserService(deleteUserID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to delete user.",
        })
        return
	}

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully delete user.",
    })

}




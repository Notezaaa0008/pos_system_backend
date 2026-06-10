package auth

import (
	"errors"
	"fmt"
	"gin-quickstart/internal/models"
	authdto "gin-quickstart/internal/module/auth/dto"
	"gin-quickstart/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authServiceInterface interface {
	RegisterSystemAdminService(req *authdto.RegisterSystemAdminRequest) error
	RegisterUserService(req *authdto.RegisterUserRequest, userId uuid.UUID) error
	LoginService(req *authdto.LoginRequest, userAgent string) (string, string, *models.User, error)
	LogoutService(userId uuid.UUID, rawRefreshToken string, allDevices bool) error
	ForgotPasswordService(req *authdto.ForgotPasswordRequest) error
	ResetPasswordService(req *authdto.ResetPasswordRequest) error
}

type AuthController struct {
	service authServiceInterface
}

func NewAuthController (service authServiceInterface) *AuthController {
	return &AuthController{service: service}
}

func (authCtrl *AuthController) clearAuthCookies(c *gin.Context, isProduction bool) {
    // การใส่ MaxAge = -1 คือการตะโกนบอกเบราว์เซอร์ว่า "ลบคุกกี้ใบนี้ทิ้งซะเดี๋ยวนี้!"
    c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
    c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
}

func (authCtrl *AuthController) RegisterSystemAdminController(c *gin.Context) {
	var req authdto.RegisterSystemAdminRequest

	// Bind JSON data เข้ากับ Struct และ Validate เบื้องต้น
	err := c.ShouldBindJSON(&req)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = authCtrl.service.RegisterSystemAdminService(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "System Admin account created successfully.",
	})
}

func (authCtrl *AuthController) RegisterUserController(c *gin.Context) {
	userId, err := utils.GetUserIDFromCtx(c)
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	var req authdto.RegisterUserRequest

	err = c.ShouldBindJSON(&req)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = authCtrl.service.RegisterUserService(&req, userId)

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

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Account created successfully.",
	})

}

func (authCtrl *AuthController) LoginController(c *gin.Context) {
	var req authdto.LoginRequest

	// Bind JSON data เข้ากับ Struct และ Validate เบื้องต้น
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Printf(
        	"[AUTH][LOGIN][REQUEST_INVALID] path=%s error=%v",
        	c.Request.URL.Path,
        	err,
    	)

    	c.JSON(http.StatusBadRequest, gin.H{
        	"status":  "error",
        	"message": "Invalid request payload",
    	})
    	return
	}

	// ดึงข้อมูลอุปกรณ์ (DeviceInfo) อัตโนมัติจาก HTTP Header
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown Device"
	}

	accessToken, refreshToken, user, err := authCtrl.service.LoginService(&req, userAgent)

	if err != nil {
		var appErr *utils.AppError

    	if errors.As(err, &appErr) {
        	c.JSON(appErr.StatusCode, gin.H{
            	"status":  "error",
            	"message": appErr.Message,
        	})
        	return
    	}

		log.Printf("[AUTH][LOGIN][UNEXPECTED_ERROR] %v", err)

    	c.JSON(http.StatusInternalServerError, gin.H{
        	"status":  "error",
        	"message": "Internal server error",
    	})
    	return
	}

	accMaxAgeStr := os.Getenv("COOKIE_ACC_MAX_AGE")
    refMaxAgeStr := os.Getenv("COOKIE_REF_MAX_AGE")

    // ตั้งค่า Default สำรองไว้เผื่อลืมเขียนใน .env (กันระบบพัง)
    accMaxAge := 60 * 15 // 15 นาที
    refMaxAge := 3600 * 24 * 7 // 7 วัน

    if accMaxAgeStr != "" {
        if val, err := strconv.Atoi(accMaxAgeStr); err == nil {
            accMaxAge = val
        }
    }
    if refMaxAgeStr != "" {
        if val, err := strconv.Atoi(refMaxAgeStr); err == nil {
            refMaxAge = val
        }
    }

	isProduction := os.Getenv("APP_ENV") == "production"

	c.SetCookie(
		"access_token",   
    	accessToken,      
    	accMaxAge, // 15 นาที (หน่วยเป็นวินาที)
    	"/", "", isProduction, 
    	true,  // 🔒 HttpOnly: true (กัน XSS)
	)

	c.SetCookie(
 	   "refresh_token",   
    	refreshToken,      
    	refMaxAge, // 7 วัน
    	"/", "", isProduction, 
    	true,  // 🔒 HttpOnly: true (กัน XSS)
	)

    displayName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
    
	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Login successfully",
        "user": gin.H{
            "id":       	user.ID,         // หน้าบ้านอาจต้องใช้ผูกอ้างอิง
            "display_name": displayName,   // เอาไว้โชว์มุมขวาบนของเว็บ: "สวัสดีคุณ..."
            "role":     	user.Role.RoleName,   // เอาไว้ให้หน้าบ้านเช็กเพื่อ ซ่อน/แสดง ปุ่มเมนู
			"image":		user.ImageUrl,
        },
    })
}

func (authCtrl *AuthController) LogoutController(c *gin.Context) {
	// 🔄 ดึงค่า ?all=true ถ้าเป็นคำอื่นหรือไม่ได้ส่งมา จะได้ค่าเป็น false ทันที
	allQuery := strings.TrimSpace(c.Query("all"))
	allDevices, err := strconv.ParseBool(allQuery)
	if err != nil {
    	c.JSON(http.StatusBadRequest, gin.H{
        	"status":  "error",
        	"message": "Invalid query parameter 'all'. Must be a boolean value.",
    	})
    	return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
    	c.JSON(http.StatusUnauthorized, gin.H{
        	"status":  "error",
        	"message": "Refresh token is missing. Please log in again.",
    	})
    	return
	}

	userId, err := utils.GetUserIDFromCtx(c)
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	err = authCtrl.service.LogoutService(userId, refreshToken, allDevices)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
        return
    }

    // 5. ล้างคุกกี้หน้าบ้านทิ้งตามปกติ
    isProduction := os.Getenv("APP_ENV") == "production"
	authCtrl.clearAuthCookies(c, isProduction)
    // c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
    // c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "logged out successfully"})

}

func (authCtrl *AuthController) ForgotPasswordController(c *gin.Context) {
	var req authdto.ForgotPasswordRequest

	// Bind JSON data เข้ากับ Struct และ Validate เบื้องต้น
	err := c.ShouldBindJSON(&req)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = authCtrl.service.ForgotPasswordService(&req)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "If the email exists, a reset link has been sent.",
	})
}

func (authCtrl *AuthController) ResetPasswordController(c *gin.Context) {
	var req authdto.ResetPasswordRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = authCtrl.service.ResetPasswordService(&req)
	if err != nil {
		// แยกแยะเคส: ถ้า Error เกิดจากตั๋วหมดอายุ/ปลอม (ที่เราเขียนดักไว้ใน Service)
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}
		
		// เคสอื่นๆ เช่น DB พัง สั่งพ่น 500 กลับไป
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to reset password"})
		return
	}

	isProduction := os.Getenv("APP_ENV") == "production"
	authCtrl.clearAuthCookies(c, isProduction)
	// c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
	// c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "password has been reset successfully. please login again.",
	})
}
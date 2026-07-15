package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"pos-system-backend/internal/models"
	authDto "pos-system-backend/internal/module/auth/dto"
	"pos-system-backend/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authServiceInterface interface {
	RegisterSystemAdminService(req *authDto.RegisterSystemAdminRequest) error
	RegisterUserService(req *authDto.RegisterUserRequest, userId uuid.UUID, storeID uuid.UUID) error
	LoginService(req *authDto.LoginRequest, userAgent string) (string, string, *models.User, int64, error)
	LogoutService(userId uuid.UUID, rawRefreshToken string, allDevices bool) error
	ForgotPasswordService(req *authDto.ForgotPasswordRequest) error
	ResetPasswordService(req *authDto.ResetPasswordRequest) error
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
	var req authDto.RegisterSystemAdminRequest

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
	userID, err := utils.GetFromCtx(c, "userID")
	
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
		log.Printf("[RegisterUser ERROR] Failed to get userID from context: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	storeID, err := utils.GetFromCtx(c, "storeID")
	if err != nil {
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
		log.Printf("[RegisterUser ERROR] Failed to get storeID from context: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing store identity for this operation"})
        return
    }
	
	var req authDto.RegisterUserRequest
	err = c.ShouldBind(&req) // สำหรับ Form-data
    if err != nil {
		log.Printf("[RegisterUser WARN] Validation failed for user input form-data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "invalid request body format",
            "error":   err.Error(),
        })
        return
    }

	err = authCtrl.service.RegisterUserService(&req, userID, storeID)

	if err != nil {
		var appErr *utils.AppError

		if errors.As(err, &appErr) {
			log.Printf("[RegisterUser WARN] Service validation error for user %s: %s (Status: %d)", req.Email, appErr.Message, appErr.StatusCode)
            c.JSON(appErr.StatusCode, gin.H{
                "status":  "error",
                "message": appErr.Message,
            })
            return
        }
		
		log.Printf("[RegisterUser CRITICAL ERROR] Internal crash while registering user %s: %v", req.Email, err)
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
	var req authDto.LoginRequest

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

	accessToken, refreshToken, user, storeNumber, err := authCtrl.service.LoginService(&req, userAgent)

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
            "user_id":       user.ID,
            "display_name":  displayName,
            "system_role":   user.SystemRole, // 👑 "SUPER_ADMIN" หรือ "USER"
            "image":         user.ImageUrl,
            "store_number": storeNumber,     
        },
    })
}

func (authCtrl *AuthController) LogoutController(c *gin.Context) {
	// 🔄 ดึงค่า ?all=true ถ้าเป็นคำอื่นหรือไม่ได้ส่งมา จะได้ค่าเป็น false ทันที
	allQuery := strings.TrimSpace(c.Query("all"))
	var allDevices bool
	var err error

	if allQuery == "" {
		allDevices = false 
	} else {
		allDevices, err = strconv.ParseBool(allQuery)
		if err != nil {
			log.Printf("[AUTH][LOGOUT][BAD_REQUEST] invalid 'all' query value='%s' error=%v", allQuery, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid query parameter 'all'. Must be a boolean value.",
			})
			return
		}
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("[Logout Warning] Refresh token missing from cookies.")
    	c.JSON(http.StatusUnauthorized, gin.H{
        	"status":  "error",
        	"message": "Refresh token is missing. Please log in again.",
    	})
    	return
	}

	userId, err := utils.GetFromCtx(c, "userID")
	if err != nil {
		log.Printf("[Logout ERROR] Failed to retrieve userID from context: %v. Ensure auth middleware is applied.", err)
        // ถ้าแอดมินลืมใส่ Middleware หรือแปลงไทป์พลาด มันจะดีดออกตรงนี้เลย
        c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized access"})
        return
    }

	err = authCtrl.service.LogoutService(userId, refreshToken, allDevices)
    if err != nil {
		log.Printf("[Logout ERROR] Service layer failed for UserID: %s, allDevices: %t, error: %v", userId, allDevices, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
        return
    }

    // ล้างคุกกี้หน้าบ้านทิ้ง
    isProduction := os.Getenv("APP_ENV") == "production"
	authCtrl.clearAuthCookies(c, isProduction)

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "logged out successfully"})

}

func (authCtrl *AuthController) ForgotPasswordController(c *gin.Context) {
	var req authDto.ForgotPasswordRequest

	// Bind JSON data เข้ากับ Struct และ Validate เบื้องต้น
	err := c.ShouldBindJSON(&req)
	
	if err != nil {
		log.Printf("[ForgotPassword WARN] Invalid request body or email format format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = authCtrl.service.ForgotPasswordService(&req)
	
	if err != nil {
        var appErr *utils.AppError
        
        // 🔍 เคสโดนดักด้วย Business Logic ที่เราตั้งใจพ่นออกมา (เช่น คอนฟิกใน .env หาย)
        if errors.As(err, &appErr) {
            log.Printf("[ForgotPassword WARN] Service process stopped for %s: %s (Status: %d)", req.Email, appErr.Message, appErr.StatusCode)
            // ถ้าเป็นแอปทั่วไปจะส่ง 200 หลอกหน้าบ้านไปเลยเพื่อความปลอดภัย แต่ถ้าอยากให้หน้าบ้านรู้ภายในก็ส่งตาม AppError ครับ
            c.JSON(appErr.StatusCode, gin.H{"status": "error", "message": appErr.Message})
            return
        }

        // 🚨 เคสที่ระบบแครช/DB ล่ม/เน็ตดับ (Unexpected Error)
        log.Printf("[ForgotPassword CRITICAL ERROR] Crash while processing reset password for %s: %v", req.Email, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Internal server error. Something went wrong.",
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "If the email exists, a reset link has been sent.",
	})
}

func (authCtrl *AuthController) ResetPasswordController(c *gin.Context) {
	var req authDto.ResetPasswordRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Printf("[ResetPassword WARN] Invalid JSON body submitted: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body format",
			"error":   err.Error(),
		})
		return
	}

	err = authCtrl.service.ResetPasswordService(&req)
	if err != nil {
		var appErr *utils.AppError
        
        if errors.As(err, &appErr) {
            log.Printf("[ResetPassword WARN] Request rejected by service: %s (Status: %d)", appErr.Message, appErr.StatusCode)
            c.JSON(appErr.StatusCode, gin.H{"status": "error", "message": appErr.Message})
            return
        }
        
        log.Printf("[ResetPassword CRITICAL ERROR] Internal engine crash during reset processing: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to reset password"})
        return
	}

	isProduction := os.Getenv("APP_ENV") == "production"
	authCtrl.clearAuthCookies(c, isProduction)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "password has been reset successfully. please login again.",
	})
}
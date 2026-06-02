package auth

import (
	"fmt"
	"gin-quickstart/internal/models"
	authdto "gin-quickstart/internal/module/auth/dto"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthServiceInterface interface {
	RegisterSuperAdminService(req *authdto.RegisterSuperAdminRequest) error
	LoginService(req *authdto.LoginRequest, userAgent string) (string, string, *models.User, error)
	LogoutService(userIDStr string, rawRefreshToken string, allDevices bool) error
	ForgotPasswordService(req *authdto.ForgotPassword) error
	ResetPasswordService(req *authdto.ResetPassword) error
}

type AuthController struct {
	service AuthServiceInterface
}

func NewAuthController (service AuthServiceInterface) *AuthController {
	return &AuthController{service: service}
}

func (authCtrl *AuthController) clearAuthCookies(c *gin.Context, isProduction bool) {
    // การใส่ MaxAge = -1 คือการตะโกนบอกเบราว์เซอร์ว่า "ลบคุกกี้ใบนี้ทิ้งซะเดี๋ยวนี้!"
    c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
    c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
}

func (authCtrl *AuthController) SignupSuperAdminController(c *gin.Context) {
	var req authdto.RegisterSuperAdminRequest

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

	err = authCtrl.service.RegisterSuperAdminService(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(), // จะพ่นข้อความภาษาไทยที่เราตั้งไว้ใน Service ออกไปหาหน้าบ้านทันที
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Super Admin account created successfully.",
	})
}

func (authCtrl *AuthController) LoginController(c *gin.Context) {
	var req authdto.LoginRequest

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

	// ดึงข้อมูลอุปกรณ์ (DeviceInfo) อัตโนมัติจาก HTTP Header
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown Device"
	}

	accessToken, refreshToken, user, err := authCtrl.service.LoginService(&req, userAgent)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "invalid username or password.",
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

	var displayName string

	// ถ้ามีทั้งชื่อและนามสกุล
	if user.FirstName != "" && user.LastName != "" {
    	displayName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	} else if user.FirstName != "" {
    	displayName = user.FirstName // มีแต่ชื่อ เอาแค่ชื่อ
	} else {
    	displayName = user.UserName // ไม่มีเลย เอา Username ไปกินก่อน
	}

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Login successfully",
        "user": gin.H{
            "id":       	user.ID,         // หน้าบ้านอาจต้องใช้ผูกอ้างอิง
            "display_name": displayName,   // เอาไว้โชว์มุมขวาบนของเว็บ: "สวัสดีคุณ..."
            "role":     	user.Role.RoleName,   // เอาไว้ให้หน้าบ้านเช็กเพื่อ ซ่อน/แสดง ปุ่มเมนู
        },
    })
}

func (authCtrl *AuthController) LogoutController(c *gin.Context) {
	// 🔄 ดึงค่า ?all=true ถ้าเป็นคำอื่นหรือไม่ได้ส่งมา จะได้ค่าเป็น false ทันที
	allDevices := c.Query("all") == "true"
	refreshToken, _ := c.Cookie("refresh_token")
	userIDStr, _ := c.Get("userID")

	err := authCtrl.service.LogoutService(userIDStr.(string), refreshToken, allDevices)
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
	var req authdto.ForgotPassword

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
	var req authdto.ResetPassword

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
package auth

import (
	"fmt"
	authdto "gin-quickstart/internal/module/auth/dto"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *AuthService
}

func NewAuthController (service *AuthService) *AuthController {
	return &AuthController{service: service}
}

func (authCtrl *AuthController) SignupSuperAdminController (c *gin.Context) {
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

func (userCtrl *AuthController) LoginController (c *gin.Context) {
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

	accessToken, refreshToken, user, err := userCtrl.service.LoginService(&req, userAgent)

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

	c.JSON(200, gin.H{
        "status":  "success",
        "message": "Login successfully",
        "user": gin.H{
            "id":       	user.ID,         // หน้าบ้านอาจต้องใช้ผูกอ้างอิง
            "display_name": displayName,   // เอาไว้โชว์มุมขวาบนของเว็บ: "สวัสดีคุณ..."
            "role":     	user.Role.RoleName,   // เอาไว้ให้หน้าบ้านเช็กเพื่อ ซ่อน/แสดง ปุ่มเมนู
        },
    })
}
package middleware

import (
	"gin-quickstart/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MiddlewareAuthService interface {
	ValidateRefreshTokenService(hashedRefreshToken string) (bool, error)
}

func AuthWithRefreshMiddleware(authService MiddlewareAuthService) gin.HandlerFunc {

	return func(c *gin.Context) {
		isProduction := os.Getenv("APP_ENV") == "production"

		// เอา Access Token มาตรวจก่อน
		accessToken, err := c.Cookie("access_token")
		if err == nil {
			//มี Access Token ลองเอาไปแกะดูว่ายังไม่หมดอายุใช่ไหม
			claims, err := utils.ParseAndValidateToken(accessToken)
			if err == nil {
				// ✅ ตั๋วผ่านปกติ! เอาสิทธิ์ยัดใส่คอนเท็กซ์ แล้วให้เดินผ่านไปเลย
				c.Set("userID", claims.UserID)
				c.Set("roleName", claims.RoleID)
				c.Next()
				return
			}
		}

		// ถ้า Access Token ไม่มี หรือ หมดอายุแล้ว
		// เราจะหันมาพึ่งพา Refresh Token
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			log.Printf("[Middleware Warning] Both tokens missing or expired.")
			// ถ้าไม่มีทั้งคู่แปลว่าไม่ได้ล็อกอินมา หรือลบคุกกี้ทิ้ง ส่ง 401 บล็อกไว้เลย
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized: session expired"})
			c.Abort() // 🚨 สั่ง Gin ว่า: "เฮ้ย หยุดจ่ายงาน! ล็อกประตูเดี๋ยวนี้!"
			return
		}

		// ตรวจสอบไส้ในของ Refresh Token
		refreshClaims, err := utils.ParseAndValidateToken(refreshToken)
		if err != nil {
			log.Printf("[Middleware Warning] Refresh token validation failed. Error: %v", err)
			// 🚨 เคสที่ตั๋วหมดอายุจริง หรือ ตั๋วปลอม: สั่งทำลายคุกกี้เน่าบนเบราว์เซอร์ทิ้งทันที!
			c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
			c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
			// ถ้า Refresh Token ก็ปลอมหรือหมดอายุด้วย ยูสเซอร์ต้องไปล็อกอินใหม่จริงๆ
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized: please login again"})
			c.Abort() // 🚨 สั่ง Gin ว่า: "เฮ้ย หยุดจ่ายงาน! ล็อกประตูเดี๋ยวนี้!"
			return
		}

		isValid, err := authService.ValidateRefreshTokenService(refreshToken)

		if err != nil {
			log.Printf("[Middleware ERROR] Database/Service error during refresh token validation: %v", err)
            // ❌ ไม่ลบคุกกี้! แค่แจ้งว่าระบบหลังบ้านมีปัญหาชั่วคราว
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error", 
                "message": "internal server error: please try again later",
            })
            c.Abort()
            return
		}

		if !isValid {
            log.Printf("[Middleware Warning] Refresh Token has been revoked or invalid for UserID: %s", refreshClaims.UserID)
            // 🔴 ล้างคุกกี้เน่าทิ้งทันที เพราะตั๋วใบนี้ใช้ไม่ได้อีกต่อไปแล้ว
            c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
            c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
            c.JSON(http.StatusUnauthorized, gin.H{
                "status": "error", 
                "message": "session expired or revoked: please login again",
            })
            c.Abort()
            return
        }

		// เมื่อ Refresh Token ยังไม่หมดอายุหรือถูกต้อง แปลว่าปลอดภัย
		// ออก Access Token ใหม่ให้
		userIDStr := refreshClaims.UserID
		roleIDStr := refreshClaims.RoleID

		// 💡 เรียกใช้ฟังก์ชันเจน JWT 
		timeAccessTokenStr := os.Getenv("TIME_ACC_TOKEN")

		if timeAccessTokenStr == "" {
			log.Println("[Middleware ERROR] Security configuration missing: TIME_ACC_TOKEN is not set in .env")
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error: security configuration missing"})
            c.Abort()
            return
		}

		timeAccessToken, err := strconv.Atoi(timeAccessTokenStr)
	
		if err != nil {
			log.Printf("[Middleware ERROR] Invalid format for TIME_ACC_TOKEN in .env: %v (must be an integer)", err)
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error: invalid security configuration format"})
            c.Abort()
            return
		}

		durationTimeAccessToken := time.Minute * time.Duration(timeAccessToken)

		newAccessToken, err := utils.GenerateJWT(userIDStr, roleIDStr, durationTimeAccessToken) 

		if err != nil {
			log.Printf("[Middleware ERROR] Failed to generate new access token for UserID %s: %v", userIDStr, err)
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to renew session"})
            c.Abort()
            return
		}

		// ส่ง Access Token ใหม่ แปะกลับไปในคุกกี้ของเบราว์เซอร์ User ทันที
		c.SetCookie(
			"access_token",   
			newAccessToken,      
			timeAccessToken, // 15 นาที (หน่วยวินาที)
			"/", "", isProduction, true,  
		)

		c.SetCookie(
			"user_role",   
			roleIDStr,      
			timeAccessToken, // 15 นาที (หน่วยวินาที)
			"/", "", isProduction, 
			true,  // 🔒 HttpOnly: true (กัน XSS)
		)

		// ฝังค่าลง Context ให้เรียบร้อยเพื่อส่งไม้ต่อให้ Controller
		c.Set("userID", userIDStr)
		c.Set("roleName", roleIDStr)

		// 🚀 เดินต่อเข้าไปทำงานได้เลย
		c.Next()

	}
}
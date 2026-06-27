package middleware

import (
	"log"
	"net/http"
	"os"
	"pos-system-backend/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MiddlewareAuthService interface {
	ValidateRefreshTokenService(hashedRefreshToken string) (bool, error)
}

func AuthMiddleware(authService MiddlewareAuthService) gin.HandlerFunc {

	return func(c *gin.Context) {
		isProduction := os.Getenv("APP_ENV") == "production"

		// เอา Access Token มาตรวจก่อน
		accessToken, err := c.Cookie("access_token")
		if err == nil {
			// มี Access Token -> ลองแกะเช็กวันหมดอายุและความถูกต้อง
			claims, err := utils.ParseAndValidateToken(accessToken)
			if err == nil {
				// ✅ ตั๋วผ่านปกติ! เซฟสิทธิ์หลักลง Context แล้วให้เดินผ่านไปด่านถัดไป
				c.Set("userID", claims.UserID)
				c.Set("systemRole", claims.SystemRole)
				c.Next()
				return
			}
		}

		// ถ้า Access Token ไม่มี หรือ หมดอายุแล้ว
		// เราจะหันมาพึ่งพา Refresh Token
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			log.Printf("[Middleware WARN] Both tokens missing or expired. Access denied.")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized: session expired"})
			c.Abort()
			return
		}

		// ตรวจสอบไส้ในของ Refresh Token
		refreshClaims, err := utils.ParseAndValidateToken(refreshToken)
		if err != nil {
			log.Printf("[Middleware WARN] Refresh token validation failed. Error: %v", err)
			// 🚨 ตั๋วปลอมหรือหมดอายุจริง: สั่งล้างทำลายคุกกี้เน่าบนเบราว์เซอร์ทิ้งทันที
			c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
			c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized: please login again"})
			c.Abort()
			return
		}

		isValid, err := authService.ValidateRefreshTokenService(refreshToken)
		if err != nil {
			log.Printf("[Middleware ERROR] Database/Service error during refresh token validation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "internal server error: please try again later",
			})
			c.Abort()
			return
		}

		if !isValid {
			log.Printf("[Middleware WARN] Refresh Token has been revoked or invalid for UserID: %s", refreshClaims.UserID)
			// 🔴 โดนเตะออกจากระบบชัวร์ ๆ ล้างคุกกี้ทิ้งซะ
			c.SetCookie("access_token", "", -1, "/", "", isProduction, true)
			c.SetCookie("refresh_token", "", -1, "/", "", isProduction, true)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "session expired or revoked: please login again",
			})
			c.Abort()
			return
		}

		// เมื่อ Refresh Token ยังไม่หมดอายุหรือถูกต้อง แปลว่าปลอดภัย
		// ออก Access Token ใหม่ให้
		userIDStr := refreshClaims.UserID
		systemRoleStr := refreshClaims.SystemRole

		timeAccessTokenStr := os.Getenv("TIME_ACC_TOKEN")
		if timeAccessTokenStr == "" {
			log.Println("[Middleware ERROR] Security configuration missing: TIME_ACC_TOKEN is not set in .env")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error: security configuration missing"})
			c.Abort()
			return
		}

		timeAccessToken, err := strconv.Atoi(timeAccessTokenStr)
		if err != nil {
			log.Printf("[Middleware ERROR] Invalid format for TIME_ACC_TOKEN in .env: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error: invalid security configuration"})
			c.Abort()
			return
		}

		durationTimeAccessToken := time.Minute * time.Duration(timeAccessToken)

		newAccessToken, err := utils.GenerateJWT(userIDStr, systemRoleStr, durationTimeAccessToken) 
		if err != nil {
			log.Printf("[Middleware ERROR] Failed to generate new access token for UserID %s: %v", userIDStr, err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to renew session"})
			c.Abort()
			return
		}

		cookieMaxAge := timeAccessToken * 60

		// ส่ง Access Token ใหม่ แปะกลับไปในคุกกี้ของเบราว์เซอร์ User ทันที
		c.SetCookie(
			"access_token",   
			newAccessToken,      
			cookieMaxAge, // 15 นาที (หน่วยวินาที)
			"/", "", isProduction, true,  
		)

		// ฝังค่าลง Context ให้เรียบร้อยเพื่อส่งไม้ต่อให้ Controller
		c.Set("userID", userIDStr)
		c.Set("systemRole", systemRoleStr)

		// 🚀 เดินต่อเข้าไปทำงานได้เลย
		c.Next()

	}
}
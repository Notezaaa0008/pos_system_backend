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

func AuthWithRefreshMiddleware() gin.HandlerFunc {

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
			// ถ้าไม่มีทั้งคู่แปลว่าไม่ได้ล็อกอินมา หรือลบคุกกี้ทิ้ง ส่ง 401 บล็อกไว้เลย
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized: session expired"})
			c.Abort() // 🚨 สั่ง Gin ว่า: "เฮ้ย หยุดจ่ายงาน! ล็อกประตูเดี๋ยวนี้!"
			return
		}

		// ตรวจสอบไส้ในของ Refresh Token
		refreshClaims, err := utils.ParseAndValidateToken(refreshToken)
		if err != nil {
			// ถ้า Refresh Token ก็ปลอมหรือหมดอายุด้วย จบเกมครับ ยูสเซอร์ต้องไปล็อกอินใหม่จริงๆ
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized: please login again"})
			c.Abort() // 🚨 สั่ง Gin ว่า: "เฮ้ย หยุดจ่ายงาน! ล็อกประตูเดี๋ยวนี้!"
			return
		}

		// เมื่อ Refresh Token ยังไม่หมดอายุหรือถูกต้อง แปลว่าปลอดภัย
		// ออก Access Token ใหม่ให้
		userIDStr := refreshClaims.UserID
		roleIDStr := refreshClaims.RoleID

		// 💡 เรียกใช้ฟังก์ชันเจน JWT 
		timeAccessTokenStr := os.Getenv("TIME_ACC_TOKEN")

		if timeAccessTokenStr == "" {
			log.Println("Error: TIME_ACC_TOKEN is missing in .env")
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error: security configuration missing"})
            c.Abort()
            return
		}

		timeAccessToken, err := strconv.Atoi(timeAccessTokenStr)
	
		if err != nil {
			log.Println("Admin Warning: TIME_ACC_TOKEN in .env must be a number:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error: invalid security configuration format"})
            c.Abort()
            return
		}

		durationTimeAccessToken := time.Minute * time.Duration(timeAccessToken)

		newAccessToken, err := utils.GenerateJWT(userIDStr, roleIDStr, durationTimeAccessToken) 

		if err != nil {
			log.Println("Error generating new access token:", err)
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

		// ฝังค่าลง Context ให้เรียบร้อยเพื่อส่งไม้ต่อให้ Controller
		c.Set("userID", userIDStr)
		c.Set("roleName", roleIDStr)

		// 🚀 เดินต่อเข้าไปทำงานได้เลย
		c.Next()

	}
}
package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	jwt.RegisteredClaims // ฟิลด์มาตรฐาน เช่น exp (หมดอายุ), iat (เวลาที่สร้าง)
}

func GenerateJWT(userID string, roleID string, duration time.Duration) (string, error){
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if secretKey == "" {
		log.Println("Error: JWT_SECRET_KEY is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return "", errors.New("JWT secret key is missing in environment configuration")
	}

	// จัดเตรียมข้อมูลสิทธิ์ที่ต้องการฝังลงใน Token
	claims := MyCustomClaims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			// กำหนดเวลาหมดอายุ (ยึดตาม duration ที่ส่งเข้ามา)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			// กำหนดเวลาที่สร้าง Token ชิ้นนี้
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// เลือกรหัสอัลกอริทึมในการเข้ารหัส (นิยมใช้ HS256) พร้อมแนบ Claims เข้าไป
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. เซ็นรับรอง Token ด้วยคีย์ลับของเรา เพื่อเปลี่ยนมันให้เป็น String ยาวๆ ส่งออกไป
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Printf("[JWT][SIGN_ERROR] %v", err)
		return "", err
	}

	return tokenString, nil
}

func ParseAndValidateToken(tokenStr string) (*MyCustomClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if secretKey == "" {
		log.Println("[JWT ERROR] JWT_SECRET_KEY is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return nil, errors.New("JWT secret key is missing in environment configuration")
	}

	jwtSecret := []byte(secretKey)

	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        // แนะนำให้ดักจับ Signing Method เพิ่มเติมเพื่อความปลอดภัย
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })

	if err != nil || !token.Valid {
		log.Printf("[JWT Warning] Token parsing failed: %v", err)
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*MyCustomClaims);

	if  ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("cannot parse claims")
}
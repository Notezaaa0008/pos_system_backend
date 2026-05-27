package main

import (
	"fmt"
	"gin-quickstart/internal/database"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/routes"
	"gin-quickstart/pkg/validator"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// 🌟 บังคับแอปพลิเคชัน Go ทั้งตัวให้ล็อกเวลาเป็นโซนไทย (Asia/Bangkok)
    loc, err := time.LoadLocation("Asia/Bangkok")
    if err == nil {
        time.Local = loc // เปลี่ยนมาตรฐานคำสั่ง time.Now() ทั้งแอปให้เป็นเวลาไทย
    }
	
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	port := os.Getenv("PORT")

	if port == "" {
        port = "8080" // default port
    }
	
	

	db := database.ConnectDatabase()

	err = models.Migrate(db)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	server := gin.Default()

	// ลงทะเบียน Custom Validator ชื่อ "strong_password"
	validator.CustomValidatorPassword()

	// ตั้งค่า CORS Middleware
    // server.Use(cors.New(cors.Config{
    //     AllowOrigins:     []string{"http://localhost:3000"}, // ใส่ URL ของ Frontend คุณ
    //     AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    //     AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    //     AllowCredentials: true,
    // }))

    // หรือถ้าจะเอาแบบง่าย (อนุญาตหมดทุกที่ - ไม่ค่อยปลอดภัยแต่สะดวกตอนพัฒนา)
    server.Use(cors.Default())

	routes.InitRouter(server, db)

	
	server.Run(fmt.Sprintf(":%s", port)) // localhost:8080

}

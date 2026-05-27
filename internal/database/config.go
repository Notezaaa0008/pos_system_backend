package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDatabase () *gorm.DB{
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("database connection string is missing in environment configuration")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// ตรวจสอบการเชื่อมต่อจริงผ่านตัวแปร sql.DB (Optional)
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("Database is unreachable:", err)
    }

	// ตั้งค่าจำนวนการเชื่อมต่อ
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
    sqlDB.SetMaxIdleConns(10)

    // SetMaxOpenConns sets the maximum number of open connections to the database.
    sqlDB.SetMaxOpenConns(100)

    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
    sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established")

	return db
}
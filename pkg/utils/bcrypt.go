package utils

import (
	"errors"
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) ([]byte ,error) {

	costStr := os.Getenv("HASH_COST")
	if costStr == "" {
		log.Println("Error: HASH_COST is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return nil, errors.New("internal server error: security configuration is missing")
	}

	// แปลงจาก string เป็น int 
	cost, err := strconv.Atoi(costStr)

	if err != nil {
		// ถ้ากรอกมาไม่ใช่ตัวเลข (เช่น HASH_COST=twelve) ให้ดีดข้อความฟ้องทันที ไม่ปล่อยผ่านแล้ว
		log.Println("Admin Warning: HASH_COST in .env must be a number:", err)
		return nil, errors.New("internal server error: invalid security configuration format")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Printf("Bcrypt Error: failed to hash password with cost %d. Details: %v", cost, err)
		return nil, errors.New("internal server error: failed to process user registration")
	}

	return hashedPassword, nil
}

func ComparePassword(hashPassword string, password string) (error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	if err != nil {
		// ตรวจสอบว่าเป็นเรื่องรหัสผ่านไม่ตรงกันจริง ๆ หรือไม่
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return errors.New("invalid username or password")
        }
        // เผื่อกรณีเกิด error อื่น ๆ ที่ไม่ได้เกี่ยวกับรหัสผ่านผิด (เช่น ค่า hash ปลอมจน bcrypt ตรวจสอบไม่ได้)
        return err
	}
	
	return nil
}
package models

import "gorm.io/gorm"


func Migrate(db *gorm.DB) error {
	// สร้างตารางที่ไม่มีให้เอง แต่การ update ชื่อ column หรือ ลบ column จะไม่สามารถทำให้ได้
    return db.AutoMigrate(
        &User{},
        &Role{},
        &RefreshToken{},
        &ResetPassword{},
        // &Post{},
        // &Product{}, // เพิ่มตารางใหม่ๆ ที่นี่ที่เดียว
    )
}
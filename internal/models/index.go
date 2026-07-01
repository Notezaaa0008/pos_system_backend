package models

import "gorm.io/gorm"


func Migrate(db *gorm.DB) error {
	// สร้างตารางที่ไม่มีให้เอง แต่การ update ชื่อ column หรือ ลบ column จะไม่สามารถทำให้ได้
    // ต้องเรียงตามลำดับตารางแม่ตารางลูกด้วย
    return db.AutoMigrate(
        &LogSendEmail{},

        &Prefix{},
        &Role{},
        &Status{},
        &Unit{},
        &Store{},
        &StoreUnit{},
        &StorePhone{},

        &Province{},
        &District{},
        &Subdistrict{},
        &PostalCode{},
        &PostalCodeArea{},
        &StoreAddress{},

        &ProductCategory{},
        &ProductSubcategory{},
        &Product{},
        &ProductPicture{},
        &ProductOption{},
        &ProductOptionItem{},

        &StockTransaction{},

        &Order{},
        &OrderItem{},

        &Discount{},
        &DiscountProduct{},

        &User{},
        &RefreshToken{},
        &ResetPassword{},
        &UserStore{},
        // &Promotion{},
    )
}
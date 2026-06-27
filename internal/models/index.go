package models

import "gorm.io/gorm"


func Migrate(db *gorm.DB) error {
	// สร้างตารางที่ไม่มีให้เอง แต่การ update ชื่อ column หรือ ลบ column จะไม่สามารถทำให้ได้
    return db.AutoMigrate(
        &DiscountProduct{},
        &Discount{},
        &District{},
        &LogSendEmail{},
        &OrderItem{},
        &Order{},
        &PostalCodeArea{},
        &PostalCode{},
        &Prefix{},
        &ProductCategory{},
        &ProductOptionItem{},
        &ProductOption{},
        &ProductPicture{},
        &ProductSubcategory{},
        &Product{},
        &Promotion{},
        &Province{},
        &RefreshToken{},
        &ResetPassword{},
        &Role{},
        &Status{},
        &StockTransaction{},
        &StoreAddress{},
        &StorePhone{},
        &StoreUnit{},
        &Store{},
        &Subdistrict{},
        &Unit{},
        &UserStore{},
        &User{},
    )
}
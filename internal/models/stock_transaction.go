package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type StockTransaction struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Quantity                decimal.Decimal     `gorm:"type:numeric(12,2);not null"`
	BeforeQuantity          decimal.Decimal     `gorm:"type:numeric(12,2);not null"` 
    AfterQuantity           decimal.Decimal     `gorm:"type:numeric(12,2);not null"`
	TransactionType         string              `gorm:"type:varchar(50);not null;index"`// IN (รับเข้า), OUT (จ่ายออก)
    Reason                  string              `gorm:"type:varchar(50);not null;index"` // SALE, RECEIVE, ADJUST_ADD, ADJUST_SUB, WASTE, RETURN
	ReferenceID             *uuid.UUID          `gorm:"type:uuid;index"`                 // ID อ้างอิง เช่น OrderID (ยอมให้ null ได้)
    Note                    string              `gorm:"type:varchar(255)"`               // หมายเหตุเพิ่มเติม เช่น "นับสต็อกประจำปีแล้วสินค้าหาย"
	CreatedAt 				time.Time 			`gorm:"not null"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid; not null"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	StoreID					uuid.UUID			`gorm:"type:uuid;not null;index"`
	Store					Store				`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ProductID				uuid.UUID			`gorm:"type:uuid;not null;index"`
	Product					Product				`gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
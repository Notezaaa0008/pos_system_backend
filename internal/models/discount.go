package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Discount struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DiscountType			string				`gorm:"type:varchar(50);not null"`
	DiscountCode			*string				`gorm:"type:varchar(50);index"`
	DiscountName			string				`gorm:"type:varchar(255);not null"`

	// เอาไว้แยกขอบเขตส่วนลด (เช่น GLOBAL = ทั้งบิล, SPECIFIC = เฉพาะสินค้าที่เลือก)
    ApplyScope              string              `gorm:"type:varchar(20);not null;default:'GLOBAL'"`
	ValueType               string              `gorm:"type:varchar(20);not null;default:'FIXED'"`
	Value 					decimal.Decimal		`gorm:"type:numeric(12,2);not null;default:0"`
	MinOrderAmount          decimal.Decimal     `gorm:"type:numeric(12,2);not null;default:0"`
	MaxDiscountAmount       *decimal.Decimal    `gorm:"type:numeric(12,2);default:null"`
	IsActive  				bool				`gorm:"not null;default:true;"`

	// type: date จะไม่ใส่เวลาลงไป
	StartDate				time.Time			`gorm:"type:date;not null;index"`
	EndDate					time.Time			`gorm:"type:date;not null;index"`
	CreatedAt 				time.Time 			`gorm:"not null"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	StoreID					uuid.UUID			`gorm:"type:uuid;not null;index"`
	Store					Store				`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	DiscountProducts		[]DiscountProduct	`gorm:"goreignKey:DiscountID"`
}
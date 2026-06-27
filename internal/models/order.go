package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderCode				string				`gorm:"type:varchar(50);not null;uniqueIndex:idx_store_order_code"`
	TotalAmount				decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	DiscountAmount			decimal.Decimal		`gorm:"type:numeric(12,2);default:0"`
	NetAmount				decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	Vat 					decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	CreatedAt 				time.Time 			`gorm:"not null;index"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null;index"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	OrderStatusID			uuid.UUID			`gorm:"type:uuid;not null;index"`
	OrderStatus				Status				`gorm:"foreignKey:OrderStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	OrderPaymentMethodID	uuid.UUID			`gorm:"type:uuid;not null;index"`
	OrderPaymentMethod		Status				`gorm:"foreignKey:OrderPaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	StoreID					uuid.UUID			`gorm:"type:uuid;not null;uniqueIndex:idx_store_order_code"`
	Store					Store				`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	OrderItems				[]OrderItem			`gorm:"foreignKey:OrderID"`
}
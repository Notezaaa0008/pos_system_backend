package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// ใส่ไว้เพราะถ้าอนาคตชื่อ product เปลี่ยน จะได้รู้ว่า ณ เวลานั้นเป็นอะไร
	ProductName				string				`gorm:"type:varchar(255);not null"`
	Quantity				decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	// ราคา ณ วันที่ขาย
	UnitPrice				decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	ItemDiscount			decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	CreatedAt 				time.Time 			`gorm:"not null"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	OrderID					uuid.UUID			`gorm:"type:uuid;not null;index"`
	Order					Order				`gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ProductID				uuid.UUID			`gorm:"type:uuid;not null;index"`
	Product					Product				`gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
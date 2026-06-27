package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductOptionItem struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItemName				string				`gorm:"type:varchar(100);not null"`
	ExtraPrice    			decimal.Decimal 	`gorm:"type:numeric(12,2);not null;default:0"`
	CreatedAt 				time.Time 			`gorm:"not null;index"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null;index"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	ProductOptionID			uuid.UUID			`gorm:"type:uuid;not null;index"`
	ProductOption			ProductOption		`gorm:"foreignKey:ProductOptionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountProduct struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt 				time.Time 			`gorm:"not null"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	DiscountID				uuid.UUID			`gorm:"type:uuid;not null;index"`
	Discount				Discount			`gorm:"foreignKey:DiscountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ProductID				uuid.UUID			`gorm:"type:uuid;not null;index"`
	Product					Product				`gorm:"foreignKey:DiscountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
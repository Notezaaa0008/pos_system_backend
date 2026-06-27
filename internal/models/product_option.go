package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductOption struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OptionName				string				`gorm:"type:varchar(100);not null"`
	IsRequired  			bool           		`gorm:"not null;default:false"`
	CreatedAt 				time.Time 			`gorm:"not null;index"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null;index"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	ProductID				uuid.UUID			`gorm:"type:uuid;not null;index"`
	Product					Product				`gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ProductOptionItems		[]ProductOptionItem	`gorm:"foreignKey:ProductOptionID"`
}
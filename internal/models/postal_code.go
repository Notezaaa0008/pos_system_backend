package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostalCode struct {
	ID        			uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Code 				string				`gorm:"type:varchar(10);not null;uniqueIndex:idx_postal_code"`	
	IsActive  			bool				`gorm:"not null;default:true"`
	CreatedAt 			time.Time 			`gorm:"not null"`
	CreatedBy			uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 			*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID			`gorm:"type:uuid;"`
	DeletedAt 			*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy			*uuid.UUID			`gorm:"type:uuid;"`

	PostalCodeAreas		[]PostalCodeArea	`gorm:"foreignKey:PostalCodeID"`
	StoreAddresses		[]StoreAddress		`gorm:"foreignKey:PostalCodeID"`
}
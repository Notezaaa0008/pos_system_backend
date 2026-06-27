package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductPicture struct {
	ID        					uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductPictureName			string			`gorm:"type:varchar(100);not null"`
	ProductPictureOriginalName	string			`gorm:"type:varchar(100);not null"`
	ProductPictureUrl			string			`gorm:"type:varchar(255);not null"`
	IsActive  					bool			`gorm:"not null;default:true"`
	CreatedAt 					time.Time 		`gorm:"not null"`
	CreatedBy					uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt 					*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy					*uuid.UUID		`gorm:"type:uuid"`
	DeletedAt 					*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy					*uuid.UUID		`gorm:"type:uuid"`	
	
	ProductID 					uuid.UUID		`gorm:"type:uuid;not null;index"`	
	Product						Product			`gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorePhone struct {
	ID        			uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// ทั้งระบบ ห้ามมีเบอร์นี้ซ้ำกันเลยเด็ดขาด
	PhoneNumber			string			`gorm:"type:varchar(50);not null;uniqueIndex:idx_store_phone_number"`
	IsActive  			bool			`gorm:"not null;default:true"`
	CreatedAt 			time.Time 		`gorm:"not null"`
	CreatedBy			uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt 			*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID		`gorm:"type:uuid"`
	DeletedAt 			*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy			*uuid.UUID		`gorm:"type:uuid"`

	StoreID				uuid.UUID		`gorm:"type:uuid;not null;index"`
	Store				Store			`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
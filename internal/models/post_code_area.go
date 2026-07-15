package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostCodeArea struct {
	ID        			uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IsActive  			bool			`gorm:"not null;default:true"`
	CreatedAt 			time.Time 		`gorm:"not null"`
	CreatedBy			uuid.UUID		`gorm:"type:uuid; not null"`
	UpdatedAt 			*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID		`gorm:"type:uuid;"`
	DeletedAt 			*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy			*uuid.UUID		`gorm:"type:uuid;"`

	// uniqueIndex:idx_subdistrict_post_code ไม่ให้ ตำบล ซ้ำซ้อนใน รหัสไปรษณีย์ เดียวกัน และค้นหาตำบลพ่วง รหัสไปรษณีย์ ไวขึ้น
	SubdistrictID		uuid.UUID		`gorm:"type:uuid;not null;uniqueIndex:idx_subdistrict_post_code"`
	Subdistrict			Subdistrict		`gorm:"foreignKey:SubdistrictID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// uniqueIndex:idx_subdistrict_post_code ไม่ให้ ตำบล ซ้ำซ้อนใน รหัสไปรษณีย์ เดียวกัน และค้นหาตำบลพ่วง รหัสไปรษณีย์ ไวขึ้น
	PostCodeID			uuid.UUID		`gorm:"type:uuid;not null;uniqueIndex:idx_subdistrict_post_code"`
	PostCode			PostCode		`gorm:"foreignKey:PostCodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	
}
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Province struct {
	ID        				uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// uniqueIndex:idx_province_name บล็อกข้อมูลซ้ำ และได้ความเร็วเวลาเราเขียน Query ค้นหาจังหวัด
	ProvinceName			string			`gorm:"type:varchar(100);not null;uniqueIndex:idx_province_name"`
	IsActive  				bool			`gorm:"not null;default:true"`
	CreatedAt 				time.Time 		`gorm:"not null"`
	CreatedBy				uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID		`gorm:"type:uuid;"`
	DeletedAt 				*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy				*uuid.UUID		`gorm:"type:uuid;"`

	Districts				[]District		`gorm:"foreignKey:ProvinceID"`
	StoreAddresses			[]StoreAddress	`gorm:"foreignKey:ProvinceID"`
}
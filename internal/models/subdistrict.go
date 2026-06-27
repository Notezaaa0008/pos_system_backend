package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subdistrict struct {
	ID              	uuid.UUID       	`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// uniqueIndex:idx_subdistrict_name_district ไม่ให้ ตำบล ซ้ำซ้อนในอำเภอเดียวกัน และค้นหาตำบลพ่วงจำเภอไวขึ้น
	SubDistrictName 	string          	`gorm:"type:varchar(100);not null;uniqueIndex:idx_subdistrict_name_district"`
	IsActive        	bool            	`gorm:"not null;default:true;"`
	CreatedAt       	time.Time       	`gorm:"not null"`
	CreatedBy       	uuid.UUID       	`gorm:"type:uuid;not null"`
	UpdatedAt       	*time.Time      	`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy       	*uuid.UUID      	`gorm:"type:uuid"`
	DeletedAt       	*gorm.DeletedAt 	`gorm:"index"`
	DeletedBy       	*uuid.UUID      	`gorm:"type:uuid"`

	// uniqueIndex:idx_subdistrict_name_district ไม่ให้ ตำบล ซ้ำซ้อนในอำเภอเดียวกัน และค้นหาตำบลพ่วงอำเภอไวขึ้น
	DistrictID			uuid.UUID			`gorm:"type:uuid;not null;uniqueIndex:idx_subdistrict_name_district"`
	District			District			`gorm:"foreignKey:DistrictID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	PostalCodeAreas		[]PostalCodeArea	`gorm:"foreignKey:SubdistrictID"`
	StoreAddresses		[]StoreAddress		`gorm:"foreignKey:SubdistrictID"`
}
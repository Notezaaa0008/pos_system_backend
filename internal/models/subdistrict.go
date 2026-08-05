package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subdistrict struct {
	ID              	uuid.UUID       	`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// uniqueIndex:idx_subdistrict_name_district ไม่ให้ ตำบล ซ้ำซ้อนในอำเภอเดียวกัน และค้นหาตำบลพ่วงจำเภอไวขึ้น
	SubdistrictNameTh 	string          	`gorm:"type:varchar(100);not null;uniqueIndex:idx_subdistrict_name_district_th,priority:1"`
	SubdistrictNameEn 	string          	`gorm:"type:varchar(100);not null;uniqueIndex:idx_subdistrict_name_district_en,priority:1"`
	IsActive        	bool            	`gorm:"not null;default:true;"`
	CreatedAt       	time.Time       	`gorm:"not null"`
	CreatedBy       	uuid.UUID       	`gorm:"type:uuid;not null"`
	UpdatedAt       	*time.Time      	`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy       	*uuid.UUID      	`gorm:"type:uuid"`
	DeletedAt       	gorm.DeletedAt 		`gorm:"uniqueIndex:idx_subdistrict_name_district_th,priority:3;uniqueIndex:idx_subdistrict_name_district_en,priority:3"`
	DeletedBy       	*uuid.UUID      	`gorm:"type:uuid"`

	// uniqueIndex:idx_subdistrict_name_district ไม่ให้ ตำบล ซ้ำซ้อนในอำเภอเดียวกัน และค้นหาตำบลพ่วงอำเภอไวขึ้น
	DistrictID			uuid.UUID			`gorm:"type:uuid;not null;uniqueIndex:idx_subdistrict_name_district_th,priority:2;uniqueIndex:idx_subdistrict_name_district_en,priority:2"`
	District			District			`gorm:"foreignKey:DistrictID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	PostcodeAreas		[]PostcodeArea		`gorm:"foreignKey:SubdistrictID"`
	StoreAddresses		[]StoreAddress		`gorm:"foreignKey:SubdistrictID"`
}
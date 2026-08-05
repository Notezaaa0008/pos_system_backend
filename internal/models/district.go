package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type District struct {
	ID        				uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// uniqueIndex:idx_district_name_province ไม่ให้ อำเภอเมือง ซ้ำซ้อนในจังหวัดเดียวกัน และค้นหาอำเภอพ่วงจังหวัดไวขึ้น
	DistrictNameTh			string			`gorm:"type:varchar(100);not null;uniqueIndex:idx_district_name_province_th,priority:1"`
	DistrictNameEn			string			`gorm:"type:varchar(100);not null;uniqueIndex:idx_district_name_province_en,priority:1"`
	IsActive  				bool			`gorm:"not null;default:true;"`
	CreatedAt 				time.Time 		`gorm:"not null"`
	CreatedBy				uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID		`gorm:"type:uuid"`
	DeletedAt 				gorm.DeletedAt	`gorm:"uniqueIndex:idx_district_name_province_th,priority:3;uniqueIndex:idx_district_name_province_en,priority:3"`	
	DeletedBy				*uuid.UUID		`gorm:"type:uuid"`

	// uniqueIndex:idx_district_name_province ไม่ให้ อำเภอเมือง ซ้ำซ้อนในจังหวัดเดียวกัน และค้นหาอำเภอพ่วงจังหวัดไวขึ้น
	ProvinceID				uuid.UUID		`gorm:"type:uuid;not null;uniqueIndex:idx_district_name_province_th,priority:2;uniqueIndex:idx_district_name_province_en,priority:2"`
	Province				Province		`gorm:"foreignKey:ProvinceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Subdistricts			[]Subdistrict	`gorm:"foreignKey:DistrictID"`
	StoreAddresses			[]StoreAddress	`gorm:"foreignKey:DistrictID"`
}
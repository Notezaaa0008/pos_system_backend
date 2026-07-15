package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreAddress struct {
	ID        			uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IsActive  			bool			`gorm:"not null;default:true"`
	CreatedAt 			time.Time 		`gorm:"not null"`
	CreatedBy			uuid.UUID		`gorm:"type:uuid; not null"`
	UpdatedAt 			*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID		`gorm:"type:uuid"`
	DeletedAt 			*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy			*uuid.UUID		`gorm:"type:uuid"`

	StoreID				uuid.UUID		`gorm:"type:uuid;not null;uniqueIndex:idx_store_address_store_id"`
	Store				*Store			`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`	
	
	ProvinceID			uuid.UUID		`gorm:"type:uuid;not null;index"`
	Province			Province		`gorm:"foreignKey:ProvinceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	DistrictID			uuid.UUID		`gorm:"type:uuid;not null;index"`
	District			District		`gorm:"foreignKey:DistrictID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	SubdistrictID		uuid.UUID		`gorm:"type:uuid;not null;index"`
	Subdistrict			Subdistrict		`gorm:"foreignKey:SubdistrictID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	PostCodeID			uuid.UUID		`gorm:"type:uuid;not null;index"`
	PostCode			PostCode		`gorm:"foreignKey:PostCodeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

}
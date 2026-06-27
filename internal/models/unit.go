package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Unit struct {
	ID              	uuid.UUID       	`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UnitCode 			string 				`gorm:"type:varchar(20);not null;uniqueIndex:idx_store_unit_code"` // เช่น PCS, KG, L, M
	UnitName  			string          	`gorm:"type:varchar(50);not null"`
	IsActive  			bool				`gorm:"not null;default:true"`
	CreatedAt 			time.Time 			`gorm:"not null"`
	CreatedBy			uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 			*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 			*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy			*uuid.UUID			`gorm:"type:uuid"`

	StoreUnits			[]StoreUnit			`gorm:"foreignKey:UnitID"`
}
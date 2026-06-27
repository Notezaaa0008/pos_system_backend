package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreUnit struct {
	ID        			uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IsActive        	bool            	`gorm:"not null;default:true;"`
	CreatedAt       	time.Time       	`gorm:"not null"`
	CreatedBy       	uuid.UUID       	`gorm:"type:uuid;not null"`
	UpdatedAt       	*time.Time      	`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy       	*uuid.UUID      	`gorm:"type:uuid"`
	DeletedAt       	*gorm.DeletedAt 	`gorm:"index"`
	DeletedBy       	*uuid.UUID      	`gorm:"type:uuid"`

	StoreID 			uuid.UUID			`gorm:"type:uuid;not null;index"`
	Store				Store				`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UnitID				uuid.UUID			`gorm:"type:uuid;not null;index"`
	Unit				Unit				`gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
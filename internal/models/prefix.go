package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prefix struct {
	ID        			uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PrefixName			string			`gorm:"type:varchar(100);unique;not null"`
	IsCompany			bool			`gorm:"not null;default:false"`
	CreatedAt 			time.Time 		`gorm:"not null"`
	CreatedBy			uuid.UUID		`gorm:"type:uuid; not null"`
	UpdatedAt 			*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID		`gorm:"type:uuid;"`
	DeletedAt 			*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy			*uuid.UUID		`gorm:"type:uuid;"`

	Users	  			[]User    		`gorm:"foreignKey:PrefixID"`
}
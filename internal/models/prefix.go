package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Prefix struct {
	ID        			uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TitlePrint			string			`gorm:"type:varchar(100);not null"`
	TitleName			string			`gorm:"type:varchar(100);unique;not null"`
	IsActive  			bool			`gorm:"not null;default:true"`
	IsCompany			bool			`gorm:"not null;default:false"`
	CreatedAt 			time.Time 		`gorm:"not null;default:now()"`
	CreatedBy			uuid.UUID		`gorm:"type:uuid; not null;default:'00000000-0000-0000-0000-000000000000'"`
	UpdatedAt 			*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID		`gorm:"type:uuid;"`
	DeletedAt 			*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy			*uuid.UUID		`gorm:"type:uuid;"`

	Users	  			[]User    		`gorm:"foreignKey:PrefixID"`
}
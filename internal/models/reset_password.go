package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResetPassword struct {
	ID				uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Token			string			`gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiredAt 		time.Time 		`gorm:"not null"`
	IsUsed    		bool      		`gorm:"default:false;not null"`
	CreatedAt 		time.Time 		`gorm:"not null"`
	UpdatedAt 		time.Time
	DeletedAt 		gorm.DeletedAt	`gorm:"index"`

	UserID 			uuid.UUID		`gorm:"type:uuid;not null;index"`
	User      		User      		`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
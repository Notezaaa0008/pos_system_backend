package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID				uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TokenHash  		string    		`gorm:"type:varchar(255);not null;index"`
	ClientType 		string    		`gorm:"type:varchar(20);not null"` // "web" หรือ "mobile"
	DeviceInfo 		*string    		`gorm:"type:varchar(255)"`         // เช่น "Chrome / MacOS"
	IsRevoked  		bool      		`gorm:"default:false;not null"`
	ExpiresAt  		time.Time 		`gorm:"not null;index"`
	CreatedAt  		time.Time 		`gorm:"not null"`
	CreatedBy		uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt		*time.Time		`gorm:"autoUpdateTime:false;default:nil"`
	UpdatedBy		*uuid.UUID		`gorm:"type:uuid;"`
	DeletedAt		*gorm.DeletedAt	`gorm:"index"`
	DeletedBy		*uuid.UUID		`gorm:"type:uuid;"`	

	UserID			uuid.UUID 		`gorm:"type:uuid;not null;index"`
	User			User	  		`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
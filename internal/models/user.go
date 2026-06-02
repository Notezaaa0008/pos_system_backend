package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        		uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserName  		string			`gorm:"type:varchar(100);unique;not null;index"`
	FirstName 		string			`gorm:"type:varchar(100);"`
	LastName  		string			`gorm:"type:varchar(100);"`
	Email     		string    		`gorm:"type:varchar(100);unique;not null"`
	Password  		string    		`gorm:"type:varchar(255);not null"`
	IsActive  		bool			`goพm:"not null;default:true"`
	CreatedAt 		time.Time 		`gorm:"not null"`
	UpdatedAt 		time.Time
	DeletedAt 		gorm.DeletedAt	`gorm:"index"`

	// Belongs To: ตัวเชื่อมไปหา Role
	RoleID    		uuid.UUID 		`gorm:"type:uuid;not null;index"`
	Role      		Role      		`gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	RefreshToken 	[]RefreshToken	`gorm:"foreignKey:UserID"`
	ResetPassword	[]ResetPassword `gorm:"foreignKey:UserID"`
}
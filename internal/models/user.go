package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        			uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName 			string			`gorm:"type:varchar(100);not null"`
	LastName  			string			`gorm:"type:varchar(100);not null"`
	Email     			string    		`gorm:"type:varchar(100);not null;uniqueIndex:idx_user_email_unique"`
	Password  			string    		`gorm:"type:varchar(255);not null"`
	ImageName			*string			`gorm:"type:varchar(255)"`
	ImageOriginalName	*string			`gorm:"type:varchar(255)"`
	ImageUrl			*string			`gorm:"type:varchar(255)"` //ใช้ pointer เพื่อให้ค่าเริ่มต้นเป็น null
	SystemRole   		string          `gorm:"type:varchar(50);not null;default:'USER'"`
	IsActive  			bool			`gorm:"not null;default:true"`
	CreatedAt 			time.Time 		`gorm:"not null"`
	CreatedBy			uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt 			*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID		`gorm:"type:uuid"`
	DeletedAt 			*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy			*uuid.UUID		`gorm:"type:uuid"`

	PrefixID			uuid.UUID		`gorm:"type:uuid;not null;index"`
	Prefix				Prefix			`gorm:"foreignKey:PrefixID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	RefreshTokens 		[]RefreshToken	`gorm:"foreignKey:UserID"`
	ResetPasswords		[]ResetPassword `gorm:"foreignKey:UserID"`
	UserStores			[]UserStore		`gorm:"foreignKey:UserID"`
}
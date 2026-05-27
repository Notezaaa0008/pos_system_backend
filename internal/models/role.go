package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID        		uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoleName  		string			`gorm:"type:varchar(50);unique;not null;index"`
	Description		string			`gorm:"type:varchar(255)"`
	IsActive  		bool			`gorm:"not null;default:true"`
	CreatedAt  		time.Time 		`gorm:"not null"`
	UpdatedAt  		time.Time
	DeletedAt  		gorm.DeletedAt	`gorm:"index"`

	// Has Many: บอก GORM ว่า Role หนึ่งอัน สามารถเชื่อมไปหา User ได้หลายคน
	Users	  		[]User    		`gorm:"foreignKey:RoleID"`
}
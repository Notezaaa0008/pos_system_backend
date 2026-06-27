package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status struct {
	ID        				uuid.UUID 		`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// มี index เยอะได้เพราะส่วนใหญ่สำหรับอ่านอย่างเดียว
	StatusType				string			`gorm:"type:varchar(20);not null;uniqueIndex:idx_status_type_code;uniqueIndex:idx_status_type_name"`
	StatusCode 				string			`gorm:"type:varchar(20);not null;uniqueIndex:idx_status_type_code"`
	StatusName				string			`gorm:"type:varchar(100);not null;uniqueIndex:idx_status_type_name"`	
	IsActive  				bool			`gorm:"not null;default:true;"`
	CreatedAt 				time.Time 		`gorm:"not null"`
	CreatedBy				uuid.UUID		`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time		`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID		`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt	`gorm:"index"`	
	DeletedBy				*uuid.UUID		`gorm:"type:uuid"`

	OrderStatuses			[]Order			`gorm:"foreignKey:OrderStatusID"`
	OrderPaymentMethods		[]Order			`gorm:"foreignKey:OrderPaymentMethodID"`
}
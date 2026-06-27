package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogSendEmail struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Recipient   			string    			`gorm:"type:varchar(255);not null;index"` // อีเมลปลายทาง (ทำ Index ไว้เผื่อเสิร์ชหา)
    EmailType   			string    			`gorm:"type:varchar(50);not null"`        // ประเภท เช่น "FORGOT_PASSWORD", "INVOICE"
    Subject     			string    			`gorm:"type:varchar(255);not null"`       // หัวข้ออีเมล
    Status      			string    			`gorm:"type:varchar(20);not null;index"`  // สเตตัส: "PENDING", "SUCCESS", "FAILED"
    ErrorMessage 			*string   			`gorm:"type:text"`                       // เก็บข้อความเออเรอร์จาก SMTP (ถ้าส่งพัง)
    RetryCount  			int       			`gorm:"type:int;default:0"`               // นับจำนวนครั้งที่พยายามส่งซ้ำ
	CreatedAt 				time.Time 			`gorm:"not null"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`
}
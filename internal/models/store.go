package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID        			uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreCode			string				`gorm:"type:varchar(100);not null;uniqueIndex;"`
	StoreName			string				`gorm:"type:varchar(100);not null;index"`
	Description			*string				`gorm:"type:varchar(255);"`
	IsActive  			bool				`gorm:"not null;default:true"`
	CreatedAt 			time.Time 			`gorm:"not null"`
	CreatedBy			uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 			*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID			`gorm:"type:uuid;"`
	DeletedAt 			*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy			*uuid.UUID			`gorm:"type:uuid;"`

	StorePhones			[]StorePhone		`gorm:"foreignKey:StoreID"`
	Products 			[]Product			`gorm:"foreignKey:StoreID"`
	StoreUnits			[]StoreUnit			`gorm:"foreignKey:StoreID"`
	Orders 				[]Order				`gorm:"foreignKey:StoreID"`
	Discounts			[]Discount			`gorm:"foreignKey:StoreID"`
	StockTransactions	[]StockTransaction	`gorm:"foreignKey:StoreID"`
	UserStores			[]UserStore			`gorm:"foreignKey:StoreID"`

	StoreAddress		*StoreAddress		`gorm:"foreignKey:StoreID"`
}
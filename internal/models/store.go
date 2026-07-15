package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID        			uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StoreCode			string				`gorm:"type:varchar(100);not null;uniqueIndex"`
	StoreName			string				`gorm:"type:varchar(100);not null;uniqueIndex:idx_store_code_deleted_at"`
	BranchName			string				`gorm:"type:varchar(100);not null;index"`
	Description			*string				`gorm:"type:varchar(255)"`
	PrimaryPhone		string				`gorm:"type:varchar(50);not null;uniqueIndex:idx_primary_phone_deleted_at"`
	SecondaryPhone      *string             `gorm:"type:varchar(50);index:idx_secondary_phone"`
    LineID              *string             `gorm:"type:varchar(50);index:idx_line_id"`
	IsActive  			bool				`gorm:"not null;default:true"`
	CreatedAt 			time.Time 			`gorm:"not null"`
	CreatedBy			uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 			*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy			*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 			*gorm.DeletedAt		`gorm:"uniqueIndex:idx_primary_phone_deleted_at"`	
	DeletedBy			*uuid.UUID			`gorm:"type:uuid"`

	Products 			[]Product			`gorm:"foreignKey:StoreID"`
	StoreUnits			[]StoreUnit			`gorm:"foreignKey:StoreID"`
	Orders 				[]Order				`gorm:"foreignKey:StoreID"`
	Discounts			[]Discount			`gorm:"foreignKey:StoreID"`
	StockTransactions	[]StockTransaction	`gorm:"foreignKey:StoreID"`
	UserStores			[]UserStore			`gorm:"foreignKey:StoreID"`

	StoreAddress		*StoreAddress		`gorm:"foreignKey:StoreID"`
}
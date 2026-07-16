package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID        				uuid.UUID 			`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductCode				string				`gorm:"type:varchar(255);not null;uniqueIndex:idx_product_code_unique"`
	ProductName				string				`gorm:"type:varchar(255);not null"`
	Description				*string				`gorm:"type:text"`
	Price					decimal.Decimal		`gorm:"type:numeric(12,2);not null"`
	StockQuantity			decimal.Decimal		`gorm:"type:numeric(12,2);default:0"`
	ProductAvailable		bool				`gorm:"not null;default:true"`
	CreatedAt 				time.Time 			`gorm:"not null"`
	CreatedBy				uuid.UUID			`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time			`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID			`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt		`gorm:"index"`	
	DeletedBy				*uuid.UUID			`gorm:"type:uuid"`

	ProductCategoryID 		uuid.UUID			`gorm:"type:uuid;not null;index"`
	ProductCategory			ProductCategory		`gorm:"foreignKey:ProductCategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ProductSubcategoryID	*uuid.UUID			`gorm:"type:uuid;index"`
	ProductSubcategory		*ProductSubcategory	`gorm:"foreignKey:ProductSubcategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	StoreID					uuid.UUID			`gorm:"type:uuid;not null;index"`
	Store 					Store				`gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ProductPictures			[]ProductPicture	`gorm:"foreignKey:ProductID"`
	OrderItems				[]OrderItem			`gorm:"foreignKey:ProductID"`
	DiscountProducts		[]DiscountProduct	`gorm:"foreignKey:ProductID"`
	ProductOptions			[]ProductOption		`gorm:"foreignKey:ProductID"`
	StockTransactions		[]StockTransaction	`gorm:"foreignKey:ProductID"`	
}
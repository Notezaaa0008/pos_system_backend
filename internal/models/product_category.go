package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	ID        				uuid.UUID 				`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CategoryName			string					`gorm:"type:varchar(100);not null;uniqueIndex:idx_category_name"`
	Description 			*string					`gorm:"type:text"`
	IsActive  				bool					`gorm:"not null;default:true"`
	CreatedAt 				time.Time 				`gorm:"not null"`
	CreatedBy				uuid.UUID				`gorm:"type:uuid;not null"`
	UpdatedAt 				*time.Time				`gorm:"autoUpdateTime:false;default:null"`
	UpdatedBy				*uuid.UUID				`gorm:"type:uuid"`
	DeletedAt 				*gorm.DeletedAt			`gorm:"index"`	
	DeletedBy				*uuid.UUID				`gorm:"type:uuid"`

	ProductSubcategories	[]ProductSubcategory	`gorm:"foreignKey:ProductCategoryID"`
	Products				[]Product				`gorm:"foreignKey:ProductCategoryID"`
}
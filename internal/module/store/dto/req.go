package storeDto

import "github.com/google/uuid"

type GetStoreRequest struct {
	Page   int    `json:"page" binding:"gte=0"`
	Limit  int    `json:"limit" binding:"gte=0"` 
	Search string `json:"search"`           
}

type CreateStoreRequest struct {
	StoreName			string			`json:"store_name" binding:"required"`
	BranchName			string			`json:"branch_name_name" binding:"required"`
	Description			*string 		`json:"description" binding:"omitempty"`
	PrimaryPhone		string			`json:"primary_phone" binding:"required,numeric,min=9,max=10"`
	SecondaryPhone		*string			`json:"secondary_phone" binding:"omitempty,numeric,min=9,max=10"`
	LineID				*string			`json:"line_id" binding:"omitempty,max=50"`
	ProvinceID			uuid.UUID		`json:"province_id" binding:"required"`
	DistrictID			uuid.UUID		`json:"district_id" binding:"required"`
	SubdistrictID		uuid.UUID		`json:"subdistrict_id" binding:"required"`
	PostCodeID			uuid.UUID		`json:"post_code_id" binding:"required"`
}

type UpdateStoreRequest struct {
	StoreName			string			`json:"store_name" binding:"required"`
	BranchName			string			`json:"branch_name_name" binding:"required"`
	Description			*string 		`json:"description" binding:"omitempty"`
	PrimaryPhone		string			`json:"primary_phone" binding:"required,numeric,min=9,max=10"`
	SecondaryPhone		*string			`json:"secondary_phone" binding:"omitempty,numeric,min=9,max=10"`
	LineID				*string			`json:"line_id" binding:"omitempty,max=50"`
	ProvinceID			uuid.UUID		`json:"province_id" binding:"required"`
	DistrictID			uuid.UUID		`json:"district_id" binding:"required"`
	SubdistrictID		uuid.UUID		`json:"subdistrict_id" binding:"required"`
	PostCodeID			uuid.UUID		`json:"post_code_id" binding:"required"`
}

type UpdateStorestatus struct {
	IsActive			*bool			`json:"is_active" binding:"required"`
}
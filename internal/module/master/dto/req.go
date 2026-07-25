package masterDto

import "github.com/google/uuid"


type GetAllPrefixRequest struct {
	Search			string		`json:"search"` 
}

type GetAllUnitRequest struct {
	Search			string		`json:"search"` 
}

type GetAllProviceRequest struct {
	Search			string		`json:"search"`
}

type GetAllDistrictRequest struct {
	Search			string		`json:"search"`
	Province_ID		uuid.UUID	`json:"province_id" binding:"required"`
}

type GetAllSubdistrictRequest struct {
	Search			string		`json:"search"`
	District_ID		uuid.UUID	`json:"district_id" binding:"required"`
}

type GetAllPostcodeRequest struct {
	Search			string		`json:"search"`
	Subdistrict_id	uuid.UUID	`json:"subdistrict_id" binding:"required"`
}
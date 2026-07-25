package masterDto

import "github.com/google/uuid"

type GetAllPrefixResponse struct {
	ID			uuid.UUID	`json:"id"`
	PrefixName	string		`json:"prefix_name"`
}

type GetAllUnitResponse struct {
	ID			uuid.UUID	`json:"id"`
	UnitName	string		`json:"unit_name"`
}

type GetAllProvinceResponse struct {
	ID			uuid.UUID	`json:"id"`
	Province	string		`json:"province"`
}

type GetAllDistrictResponse struct {
	ID			uuid.UUID	`json:"id"`
	District	string		`json:"district"`
}

type GetAllSubdistrictResponse struct {
	ID			uuid.UUID	`json:"id"`
	Subdistrict string		`json:"subdistrict"`
}

type GetAllPostcodeResponse struct {
	ID			uuid.UUID	`json:"id"`
	Postcode    string		`json:"postcode"`
}
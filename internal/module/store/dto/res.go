package storeDto

import "github.com/google/uuid"

type GetStoreListResponse struct {
	ID				uuid.UUID	`json:"id"`
	StoreCode		string		`json:"store_code"`
	StoreName		string		`json:"store_name"`
	RoleName		string		`json:"role_name"`
}
package prefixDto

import "github.com/google/uuid"

type GetAllPrefixResponse struct {
	ID			uuid.UUID	`json:"id"`
	PrefixName	string		`json:"prefix_name"`
}
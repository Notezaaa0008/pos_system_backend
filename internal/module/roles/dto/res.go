package roledto

import "github.com/google/uuid"

type GetAllRoleResponse struct {
	ID			uuid.UUID	`json:"id"`
	RoleName	string		`json:"role_name"`
}
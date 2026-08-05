package usersdto

import "github.com/google/uuid"

type GetProfileResponse struct {
	FirstName		string		`json:"first_name"`
	LastName		string		`json:"last_name"`
	Email			string		`json:"email"`
	ImageName		string		`json:"image_name"`
	ImageUrl		string		`json:"image_url"`
	PrefixID		uuid.UUID	`json:"prefix_id"`
}
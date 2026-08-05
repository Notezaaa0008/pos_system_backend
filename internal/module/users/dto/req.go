package usersdto

import (
	"mime/multipart"
)

type UpdateProfileRequest struct {
	FirstName	string 					`form:"first_name" binding:"required"`
	LastName	string 					`form:"last_name" binding:"required"`
	Email 		string 					`form:"email" binding:"required,email"`
	PrefixID  	string                  `form:"prefix_id" binding:"required,uuid"`

	Files 		[]*multipart.FileHeader	`form:"files"`
}
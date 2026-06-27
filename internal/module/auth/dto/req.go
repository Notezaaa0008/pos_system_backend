package authdto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type RegisterSystemAdminRequest struct {
	FirstName	string 					`json:"first_name" binding:"required"`
	LastName	string 					`json:"last_name" binding:"required"`
	Email 		string 					`json:"email" binding:"required,email"`
	Password 	string 					`json:"password" binding:"required,min=8,max=16,strong_password"`
	PrefixID	uuid.UUID				`json:"prefix_id" binding:"required"`
}

type RegisterUserRequest struct {
	FirstName	string 					`form:"first_name" binding:"required"`
	LastName	string 					`form:"last_name" binding:"required"`
	Email 		string 					`form:"email" binding:"required,email"`
	Password 	string 					`form:"password" binding:"required,min=8,max=16,strong_password"`
	RoleID		uuid.UUID				`form:"role_id" binding:"required"`
	PrefixID	uuid.UUID				`form:"prefix_id" binding:"required"`

	Files 		[]*multipart.FileHeader	`form:"files"`
}

type LoginRequest struct {
	Email 		string 					`form:"email" binding:"required,email"`
	Password	string 					`json:"password" binding:"required,min=8,max=16,strong_password"`
	Client   	string 					`json:"client" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email		string 					`json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token		string					`json:"token_reset" binding:"required"`
	NewPassword	string					`json:"new_password" binding:"required,min=8,max=16,strong_password"`
}
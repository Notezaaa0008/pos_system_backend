package authdto

type RegisterSuperAdminRequest struct {
	UserName	string `json:"user_name" binding:"required"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Email 		string `json:"email" binding:"required,email"`
	Password 	string `json:"password" binding:"required,min=8,max=16,strong_password"`
}

type LoginRequest struct {
	UserName	string `json:"user_name" binding:"required"`
	Password	string `json:"password" binding:"required,min=8,max=16,strong_password"`
	Client   	string `json:"client" binding:"required"`
}
package roledto


type CreateRoleRequest struct {
	RoleName	string 		`json:"role_name" binding:"required"`
	Description	string 		`json:"description"`
}

type UpdateRoleRequest struct {
	RoleName	string 		`json:"role_name" binding:"required"`
	Description	string 		`json:"description"`
	IsActive	*bool		`json:"is_active" binding:"required"` //ถ้าไม่ใช้ pointer gin จะมองค่า false ว่าคือการไม่ส่งอะไรมาเลย จึงต้องใช้ pointer แยก
}
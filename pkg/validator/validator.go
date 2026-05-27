package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)


func CustomValidatorPassword() {
	var (
		hasUpperRegex   = regexp.MustCompile(`[A-Z]`)
		hasLowerRegex   = regexp.MustCompile(`[a-z]`)
		hasSpecialRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	)

	// การตรวจสอบเพื่อความปลอดภัย (Safe Check) ว่าการแปลงประเภทตัวแปรสำเร็จไหม
	validate, ok := binding.Validator.Engine().(*validator.Validate); 

	if ok {
		validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
			password := fl.Field().String()
			
			// ยังคงเรียกใช้ได้ปกติ เพราะฟังก์ชันย่อยมองเห็นตัวแปรของฟังก์ชันแม่ (Closure)
			return hasUpperRegex.MatchString(password) && 
				   hasLowerRegex.MatchString(password) && 
				   hasSpecialRegex.MatchString(password)
		})
	}
}
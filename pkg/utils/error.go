package utils

import "net/http"

// AppError คือ Struct ที่เราดีไซน์ขึ้นมาเพื่ออุ้ม HTTP Status ไว้คู่กับ Error จริง
type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

// ฟังก์ชันช่วยสร้าง Error ตามประเภท
func NewNotFoundError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusNotFound, Message: msg}
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Message: msg}
}

func NewUnauthorizedError(msg string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Message: msg}
}
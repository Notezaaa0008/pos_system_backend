package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func IsBlank(s string) (string, bool) {
	value := strings.TrimSpace(s)

	if value == "" {
		return "", true
	} 

	return value, false
}

func CleanInputPhoneNumber(input string) string {
    // 1. ตัดช่องว่าง หัว-ท้าย
    trimmed := strings.TrimSpace(input)
    // 2. ลบเครื่องหมายขีด (-) ออกให้หมด
    trimmed = strings.ReplaceAll(trimmed, "-", "")
    // 3. ลบช่องว่างที่อยู่ตรงกลางออกให้หมด
    trimmed = strings.ReplaceAll(trimmed, " ", "")
    return trimmed
}

func ValidateUploadFile(
	files []*multipart.FileHeader, 
	maxFiles int, 
	maxSizeInMB int64, 
	allowedTypes []string,
) error {
	if len(files) > maxFiles {
		return NewBadRequestError(fmt.Sprintf("you can upload a maximum of %d files", maxFiles)) 
	}

	maxByteSize := maxSizeInMB * 1024 * 1024

	for _, file := range files {
		
		// 2.1 Check file size dynamically
		if file.Size > maxByteSize {
			return  NewBadRequestError(fmt.Sprintf("file '%s' exceeds the maximum allowed size of %dMB", file.Filename, maxSizeInMB))
		}

		// 2.2 Get Content-Type and extension
		contentType := file.Header.Get("Content-Type")
		ext := strings.ToLower(filepath.Ext(file.Filename))

		// 2.3 Check if type/extension is allowed dynamically
		isValidType := false
		for _, allowedType := range allowedTypes {
			// standard allowedType inputs: "jpeg", "png", "jpg"
			if strings.Contains(contentType, allowedType) || strings.TrimPrefix(ext, ".") == allowedType {
				isValidType = true
				break
			}
		}

		if !isValidType {
			// Join allowed types array into a clean string for error message (e.g., "jpeg, png")
			supportedFormats := strings.Join(allowedTypes, ", ")
			return NewBadRequestError(fmt.Sprintf("file '%s' has an invalid format. Supported formats: %s", file.Filename, supportedFormats))
		}
	}

	return nil // Return nil if all files pass the validations
}
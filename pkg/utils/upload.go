package utils

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadResult struct {
	SecureURL    string // Link URL (https://...)
	OriginalName string // ชื่อไฟล์เดิมที่ผู้ใช้อัปโหลดมา (เช่น my-photo.png)
	CloudName    string // ชื่อไฟล์ใหม่ที่ระบุบน Cloudinary (Public ID)
}

func UploadToCloudinary(fileHeader *multipart.FileHeader) (*UploadResult, error) {
	ctx := context.Background()
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")

	if cloudinaryURL == "" {
		log.Println("Error: CLOUDINARY_URL is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return nil, errors.New("internal server error: security configuration is missing")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// ดึงชื่อไฟล์เดิมออกมา (ตัดนามสกุลออกถ้าต้องการเก็บแค่ชื่อเพียวๆ)
	// เช่น "my-avatar.png" -> "my-avatar"
	origFileName := fileHeader.Filename
	nameWithoutExt := strings.TrimSuffix(origFileName, filepath.Ext(origFileName))

	uploadParams := uploader.UploadParams{
		Folder: "my_app_uploads",
		// 💡 แนะนำ: บังคับให้ Cloudinary เติมสุ่มตัวเลขต่อท้ายชื่อเดิม เพื่อป้องกันไม่ให้ชื่อไฟล์ซ้ำกันแล้วไปทับกันเองบน Cloud
		PublicID: nameWithoutExt, 
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, err
	}

	// มัดรวมข้อมูลทั้ง 3 อย่างส่งกลับไป
	return &UploadResult{
		SecureURL:    result.SecureURL,
		OriginalName: origFileName,       // เช่น "my-avatar.png"
		CloudName:    result.PublicID,   // เช่น "my_app_uploads/my-avatar_v8s2da" (หรือตามที่ Cloudinary เจนให้)
	}, nil
}
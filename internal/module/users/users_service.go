package users

import (
	"log"
	"pos-system-backend/internal/models"
	usersdto "pos-system-backend/internal/module/users/dto"
	"pos-system-backend/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type userRepositoryInterface interface {
	GetProfileRepository(userID uuid.UUID) (*models.User, error)
	UpdateProfileRepository(userID uuid.UUID, userData *models.User) error
	DeleteUserRepository(userStoreData *models.UserStore, userData *models.User, deleteUserID uuid.UUID) error
}

type UsersService struct {
	repo userRepositoryInterface
}

func NewUsersService (repo userRepositoryInterface) *UsersService{
	return &UsersService{repo: repo}
}

func (service *UsersService) GetAllUserService() {
	
}

func (service *UsersService) GetProfileService(userID uuid.UUID) (*usersdto.GetProfileResponse, error){
	profile, err := service.repo.GetProfileRepository(userID)

	if err != nil {
		return nil, err
	}

	var imageName, imageUrl string
	if profile.ImageName != nil {
		imageName = *profile.ImageName
	}
	if profile.ImageUrl != nil {
		imageUrl = *profile.ImageUrl
	}

	responseProfile := usersdto.GetProfileResponse{
		FirstName: profile.FirstName,
		LastName: profile.LastName,
		Email: profile.Email,
		ImageName: imageName,
		ImageUrl: imageUrl,
		PrefixID: profile.PrefixID,
	}

	return &responseProfile, err
}

func (service *UsersService) UpdateProfileService(userID uuid.UUID, req *usersdto.UpdateProfileRequest) error{
	firstName, isBlank := utils.IsBlank(req.FirstName)
    if isBlank { return utils.NewBadRequestError("First name is required") }

    lastName, isBlank := utils.IsBlank(req.LastName)
    if isBlank { return utils.NewBadRequestError("Last name is required") }

	prefixUUID, err := uuid.Parse(req.PrefixID)
    if err != nil {
        return utils.NewBadRequestError("Invalid prefix ID format")
    }

	maxAllowedFiles := 1
	maxAllowedSizeMB := int64(5)
	allowedFormats := []string{"jpeg", "jpg", "png"}

	var uploadResults []*utils.UploadResult
	if len(req.Files) > 0 {
		err := utils.ValidateUploadFile(req.Files, maxAllowedFiles, maxAllowedSizeMB, allowedFormats)

		if err != nil {
			log.Printf("[Update Profile Service WARN] File validation failed for user %s: %v", req.Email, err)
			return err
		}

        for _, file := range req.Files {
            // 💡 เรียกใช้ฟังก์ชันที่ปรับปรุงใหม่ จะได้ข้อมูลกลับมาครบ 3 อย่าง
            res, err := utils.UploadToCloudinary(file)
            if err != nil {
				log.Printf("[Update Profile Service ERROR] Cloudinary upload crash for user %s: %v", req.Email, err)
                return err 
            }
            uploadResults = append(uploadResults, res)
        }
    }

	var imageUrl, originalName, fileName *string
    if len(uploadResults) > 0 {
        targetFile := uploadResults[0]
		imageUrl = &targetFile.SecureURL
		originalName = &targetFile.OriginalName
		fileName = &targetFile.CloudName
    }

	now := time.Now()
	userData := models.User{
		FirstName: firstName,
		LastName: lastName,
		Email: strings.TrimSpace(req.Email),
		ImageName: fileName,
		ImageOriginalName: originalName,
		ImageUrl: imageUrl,
		PrefixID: prefixUUID,
		UpdatedBy: &userID,
		UpdatedAt: &now,
	}

	err = service.repo.UpdateProfileRepository(userID, &userData)
	if err != nil {
		log.Printf("[Service Update Profile ERROR] Failed to Update Profile Error: %v", err)
		return err
	}

	return nil
}

func (service *UsersService) DeleteUserService(deleteUserID uuid.UUID, userID uuid.UUID) error{
	userData := models.User{
		IsActive:  false,
        DeletedBy: &userID,
	}

	userStoreData := models.UserStore{
		IsActive:  false,
        DeletedBy: &userID,
	}

	err := service.repo.DeleteUserRepository(&userStoreData, &userData, deleteUserID)
	if err != nil {
		log.Printf("[Service Delete User ERROR] Failed to Delete User Error: %v", err)
		return err
	}

	return nil
}

package auth

import (
	"errors"
	"gin-quickstart/internal/models"
	authdto "gin-quickstart/internal/module/auth/dto"
	"gin-quickstart/pkg/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type authRepositoryInterface interface {
	CheckSystemAdminExists(rolesID uuid.UUID) (bool, error)
	CheckRefreshTokenValid(hashRefreshToken string) (bool, error)
	FineUserByUserName(userName string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindValidResetToken(token string) (*models.ResetPassword, error)
	CreateUser(user *models.User) error
	CreateRefreshTokenRecord(refreshToken *models.RefreshToken) error
	CreateResetPassword(reset *models.ResetPassword) error
	RevokeRefreshToken(userID uuid.UUID, hashedToken string) error
	UpdatePasswordAndRevokeToken(userID uuid.UUID, hashedPwd string, resetID uuid.UUID) error
}

type roleServiceReader interface {
    GetRoleIdByRoleNameService(roleName string) (uuid.UUID, error)
}

type AuthService struct {
	repo authRepositoryInterface
	roleService roleServiceReader
}

func NewAuthService (repo authRepositoryInterface, roleService roleServiceReader) *AuthService {
	return &AuthService{
		repo: repo,
		roleService: roleService,
	}
}

func (service *AuthService) RegisterSystemAdminService(req *authdto.RegisterSystemAdminRequest) error {
	userName, isBlank := utils.IsBlank(req.UserName)
    if isBlank { return utils.NewBadRequestError("Username is required") }

    firstName, isBlank := utils.IsBlank(req.FirstName)
    if isBlank { return utils.NewBadRequestError("First name is required") }

    lastName, isBlank := utils.IsBlank(req.LastName)
    if isBlank { return utils.NewBadRequestError("Last name is required") }   

	roleID, err := service.roleService.GetRoleIdByRoleNameService("system_admin")

	if err != nil {
		return err
	}

	exists ,err := service.repo.CheckSystemAdminExists(roleID)

	if err != nil {
		return err
	}

	// ถ้า Repo บอกว่า exists == true แปลว่าระบบมี Super Admin อยู่แล้ว
	if exists {
		return errors.New("cannot create: a Super Admin account already exists in the system.")
	}

	hashedPassword, err := utils.HashedPassword(strings.TrimSpace(req.Password))

	if err != nil {
		return err
	}

	

	newUser := models.User{
		UserName: userName,
		FirstName: firstName,
		LastName: lastName,
		Email: strings.TrimSpace(req.Email),
		Password: string(hashedPassword),
		ImageUrl: nil,
		RoleID: roleID,
		PrefixID: req.PrefixID,
		IsActive: true,
		CreatedBy: uuid.Nil,
	}

	err = service.repo.CreateUser(&newUser)
	if err != nil {
		return err
	}

	return nil
} 

func (service *AuthService) RegisterUserService(req *authdto.RegisterUserRequest, userId uuid.UUID) error {
	userName, isBlank := utils.IsBlank(req.UserName)
    if isBlank { return utils.NewBadRequestError("Username is required") }

    firstName, isBlank := utils.IsBlank(req.FirstName)
    if isBlank { return utils.NewBadRequestError("First name is required") }

    lastName, isBlank := utils.IsBlank(req.LastName)
    if isBlank { return utils.NewBadRequestError("Last name is required") }   

	maxAllowedFiles := 1
	maxAllowedSizeMB := int64(5)
	allowedFormats := []string{"jpeg", "jpg", "png"}

	var uploadResults []*utils.UploadResult
    if len(req.Files) > 0 {
		err := utils.ValidateUploadFile(req.Files, maxAllowedFiles, maxAllowedSizeMB, allowedFormats)

		if err != nil {
			return err
		}

        for _, file := range req.Files {
            // 💡 เรียกใช้ฟังก์ชันที่ปรับปรุงใหม่ จะได้ข้อมูลกลับมาครบ 3 อย่าง
            res, err := utils.UploadToCloudinary(file)
            if err != nil {
                return err 
            }
            uploadResults = append(uploadResults, res)
        }
    }

	// 3. ดักจับค่ารูปภาพในรูปแบบ Pointer ที่ถูกต้อง ป้องกันสตริงว่างหลุดลง DB
    var imageUrl *string
    var originalName *string
    var fileName *string

    if len(uploadResults) > 0 {
        targetFile := uploadResults[0]

		imageUrl = &targetFile.SecureURL
		originalName = &targetFile.OriginalName
		fileName = &targetFile.CloudName
    }

	hashedPassword, err := utils.HashedPassword(req.Password)

	if err != nil {
		return err
	}

	userData := models.User{
		UserName: userName,
		FirstName: firstName,
		LastName: lastName,
		Email: strings.TrimSpace(req.Email),
		Password: string(hashedPassword),
		ImageName: fileName,
		ImageOriginalName: originalName,
		ImageUrl: imageUrl,
		RoleID: req.RoleID,
		PrefixID: req.PrefixID,
		IsActive: true,
		CreatedBy: userId,
	}

	err = service.repo.CreateUser(&userData)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) ValidateRefreshTokenService(refreshToken string) (bool, error) {
	if refreshToken == "" {
		return false, errors.New("refresh token is required.")
	}

	hashedToken := utils.HashToken(refreshToken)

	 isValid ,err := service.repo.CheckRefreshTokenValid(hashedToken)

	 if err != nil {
		return false, err 
	 }
	
	 return isValid, nil
}

func (service *AuthService) LoginService(req *authdto.LoginRequest, userAgent string) (string, string, *models.User, error) {
	userName, isBlank := utils.IsBlank(req.UserName)
    if isBlank { return "", "", nil, utils.NewBadRequestError("Username is required") }

	user, err := service.repo.FineUserByUserName(strings.TrimSpace(userName))

	if err != nil {
		return "", "", nil, err
	}

	err = utils.ComparePassword(user.Password, req.Password)

	if err != nil {
		return "", "", nil, err
	}

	userIDStr := user.ID.String()
	roleIDStr := user.Role.ID.String()

	timeAccessTokenStr := os.Getenv("TIME_ACC_TOKEN")
	timeRefreshTokenStr := os.Getenv("TIME_REF_TOKEN")

	if timeAccessTokenStr == "" || timeRefreshTokenStr == "" {
		log.Println("Error: TIME_ACC_TOKEN OR TIME_REF_TOKEN is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return "", "", nil, errors.New("Missing in environment configuration")
	}

	// แปลงจาก string เป็น int 
	timeAccessToken, err := strconv.Atoi(timeAccessTokenStr)
	
	if err != nil {
		log.Println("Admin Warning: TIME_ACC_TOKEN in .env must be a number:", err)
		return "", "", nil, errors.New("internal server error: invalid security configuration format")
	}

	timeRefreshToken, err := strconv.Atoi(timeRefreshTokenStr)

	if err != nil {
		log.Println("Admin Warning: TIME_REF_TOKEN in .env must be a number:", err)
		return "", "", nil, errors.New("internal server error: invalid security configuration format")
	}

	durationTimeAccessToken := time.Minute * time.Duration(timeAccessToken)
	durationTimeRefreshToken := time.Hour * time.Duration(timeRefreshToken)

	accessToken, err := utils.GenerateJWT(userIDStr, roleIDStr, durationTimeAccessToken)

	if err != nil {
		return "", "", nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateJWT(userIDStr, roleIDStr, durationTimeRefreshToken)

	hashedToken := utils.HashToken(refreshToken)

	refreshTokenRecord := models.RefreshToken{
		UserID: user.ID,
		TokenHash: hashedToken,
		ClientType: strings.TrimSpace(req.Client),
		DeviceInfo: &userAgent,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(durationTimeRefreshToken),
		CreatedBy: user.ID,
	}

	err = service.repo.CreateRefreshTokenRecord(&refreshTokenRecord)

	if err != nil {
		return "", "", nil, err
	} 

	return accessToken, refreshToken, user, nil
}

func (service *AuthService) LogoutService(userId uuid.UUID, rawRefreshToken string, allDevices bool) error {

	// 🚨 เคสที่ 1: สั่งลบทุกเครื่อง (ไม่สนใจ Token เครื่องปัจจุบัน)
	if allDevices {
		// ส่ง (userID, "") -> Repo จะไปใช้เงื่อนไขคัดเฉพาะอันที่ยังไม่หมดอายุของยูสเซอร์คนนี้
		return service.repo.RevokeRefreshToken(userId, "")
	}

	// 🚨 เคสที่ 2: สั่งลบเฉพาะเครื่องปัจจุบัน
	if rawRefreshToken == "" {
		return errors.New("missing refresh token for current device logout")
	}
	
	// 🔒 ทำการแฮช Token ดิบให้เป็นค่า SHA-256 อยู่ในชั้นนี้ตามกฎธุรกิจ
	hashedToken := utils.HashToken(rawRefreshToken)

	// ส่ง (uuid.Nil, hashedToken) -> Repo จะเจาะจงทำลายใบนี้ทันที
	return service.repo.RevokeRefreshToken(uuid.Nil, hashedToken)
}

func (service *AuthService) ForgotPasswordService(req *authdto.ForgotPasswordRequest) error{
	user, err := service.repo.FindUserByEmail(strings.TrimSpace(req.Email))

	if err != nil {
		return err
	}

	token := utils.RandomToken()
	
	timeResetPasswordStr := os.Getenv("TIME_RESET_PASSWORD")

	if timeResetPasswordStr == "" {
		log.Println("Error: TIME_RESET_PASSWORD is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return errors.New("Missing in environment configuration")
	}

	// แปลงจาก string เป็น int 
	timeResetPassword, err := strconv.Atoi(timeResetPasswordStr)
	
	if err != nil {
		log.Println("Admin Warning: TIME_ACC_TOKEN in .env must be a number:", err)
		return errors.New("internal server error: invalid security configuration format")
	}
	
	expiredAt := time.Now().Add(time.Duration(timeResetPassword) * time.Minute)

	reserData := models.ResetPassword{
		UserID: user.ID,
		Token: token,
		ExpiredAt: expiredAt,
		CrearedBy: user.ID,
	}

	err = service.repo.CreateResetPassword(&reserData)

	if err != nil {
		return err
	}

	err = utils.SendResetPasswordEmail(user.Email, token)
	if err != nil {
		log.Println("failed to send email.")
		return err
	}

	return nil
}

func (service *AuthService) ResetPasswordService(req *authdto.ResetPasswordRequest) error{
	token, isBlank := utils.IsBlank(req.Token)
    if isBlank { return utils.NewUnauthorizedError("Token is required.") }
	// 1. ตรวจสอบตั๋ว (Token) ว่าถูกต้อง/ไม่หมดอายุ ไหม
	resetToken, err := service.repo.FindValidResetToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// 2. แฮชรหัสผ่านใหม่ด้วย bcrypt (เหมือนตอนสมัครสมาชิก)
	hashedPassword, err := utils.HashedPassword(strings.TrimSpace(req.NewPassword))
	if err != nil {
		return err
	}

	// 3. สั่งบันทึกรหัสใหม่ลง User และสั่งติ๊กใช้ตั๋วใบนี้แล้ว
	err = service.repo.UpdatePasswordAndRevokeToken(resetToken.UserID, string(hashedPassword), resetToken.ID)
	if err != nil {
		return err
	}

	err = service.repo.RevokeRefreshToken(resetToken.UserID, "")
	if err != nil {
		return err
	}

	return nil
}
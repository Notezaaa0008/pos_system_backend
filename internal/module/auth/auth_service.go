package auth

import (
	"errors"
	"gin-quickstart/internal/models"
	authdto "gin-quickstart/internal/module/auth/dto"
	"gin-quickstart/pkg/utils"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AuthRepositoryInterface interface {
	CheckSuperAdminExists(rolesID uuid.UUID) (bool, error)
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

type RoleServiceReader interface {
    GetRoleIdByRoleNameService(roleName string) (uuid.UUID, error)
}

type AuthService struct {
	repo AuthRepositoryInterface
	roleService RoleServiceReader
}

func NewAuthService (repo AuthRepositoryInterface, roleService RoleServiceReader) *AuthService {
	return &AuthService{
		repo: repo,
		roleService: roleService,
	}
}

func (service *AuthService) RegisterSuperAdminService(req *authdto.RegisterSuperAdminRequest) error {
	roleID, err := service.roleService.GetRoleIdByRoleNameService("super_admin")

	if err != nil {
		return err
	}

	exists ,err := service.repo.CheckSuperAdminExists(roleID)

	if err != nil {
		return err
	}

	// ถ้า Repo บอกว่า exists == true แปลว่าระบบมี Super Admin อยู่แล้ว
	if exists {
		return errors.New("cannot create: a Super Admin account already exists in the system.")
	}

	hashedPassword, err := utils.HashedPassword(req.Password)

	if err != nil {
		return err
	}

	newUser := &models.User{
		UserName: req.UserName,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Password: string(hashedPassword),
		RoleID: roleID,
		IsActive: true,
	}

	err = service.repo.CreateUser(newUser)
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
	user, err := service.repo.FineUserByUserName(req.UserName)

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
		ClientType: req.Client,
		DeviceInfo: userAgent,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(durationTimeRefreshToken),
	}

	err = service.repo.CreateRefreshTokenRecord(&refreshTokenRecord)

	if err != nil {
		return "", "", nil, err
	} 

	return accessToken, refreshToken, user, nil
}

func (service *AuthService) LogoutService(userIDStr string, rawRefreshToken string, allDevices bool) error {
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return err
	}

	// 🚨 เคสที่ 1: สั่งลบทุกเครื่อง (ไม่สนใจ Token เครื่องปัจจุบัน)
	if allDevices {
		// ส่ง (userID, "") -> Repo จะไปใช้เงื่อนไขคัดเฉพาะอันที่ยังไม่หมดอายุของยูสเซอร์คนนี้
		return service.repo.RevokeRefreshToken(userUUID, "")
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

func (service *AuthService) ForgotPasswordService(req *authdto.ForgotPassword) error{
	user, err := service.repo.FindUserByEmail(req.Email)

	if err != nil {
		return nil
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

	reserData := &models.ResetPassword{
		UserID: user.ID,
		Token: token,
		ExpiredAt: expiredAt,
	}

	err = service.repo.CreateResetPassword(reserData)

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

func (service *AuthService) ResetPasswordService(req *authdto.ResetPassword) error{
	// 1. ตรวจสอบตั๋ว (Token) ว่าถูกต้อง/ไม่หมดอายุ ไหม
	resetToken, err := service.repo.FindValidResetToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// 2. แฮชรหัสผ่านใหม่ด้วย bcrypt (เหมือนตอนสมัครสมาชิก)
	hashedPassword, err := utils.HashedPassword(req.NewPassword)
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
package auth

import (
	"errors"
	"gin-quickstart/internal/models"
	authdto "gin-quickstart/internal/module/auth/dto"
	"gin-quickstart/internal/module/roles"
	"gin-quickstart/pkg/utils"
	"log"
	"os"
	"strconv"
	"time"
)

type AuthService struct {
	repo *AuthRepository
	roleService *roles.RolesService
}

func NewAuthService (repo *AuthRepository, roleService *roles.RolesService) *AuthService {
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


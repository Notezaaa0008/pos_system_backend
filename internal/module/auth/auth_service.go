package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"pos-system-backend/internal/models"
	authdto "pos-system-backend/internal/module/auth/dto"
	"pos-system-backend/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authRepositoryInterface interface {
	CheckSystemAdminExists(systemRole string) (bool, error)
	CheckRefreshTokenValid(hashRefreshToken string) (bool, error)
	FindUserByEmail(email string, findType string) (*models.User, error)
	FindValidResetToken(token string) (*models.ResetPassword, error)
	CreateUserSystemAdmin(user *models.User) error
	CreateUser(userData *models.User, userStoreData *models.UserStore) error
	CreateRefreshTokenRecord(refreshToken *models.RefreshToken) error
	CreateLogEmail(logEmail *models.LogSendEmail) error
	CheckPermission(userID uuid.UUID, storeID uuid.UUID) (*models.UserStore, error)
	CreateResetPassword(reset *models.ResetPassword) error
	RevokeRefreshToken(userID uuid.UUID, hashedToken string) error
	UpdatePasswordAndRevokeToken(userID uuid.UUID, hashedPwd string, resetID uuid.UUID) error
	UpdateLogEmailStatus(logID uuid.UUID, status string, errMsg *string, userID uuid.UUID) error
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
    firstName, isBlank := utils.IsBlank(req.FirstName)
    if isBlank { return utils.NewBadRequestError("First name is required") }

    lastName, isBlank := utils.IsBlank(req.LastName)
    if isBlank { return utils.NewBadRequestError("Last name is required") }   

	exists ,err := service.repo.CheckSystemAdminExists("SYSTEM_ADMIN")

	if err != nil {
		log.Printf("[ERROR] AuthService.RegisterSystemAdminService - CheckSystemAdminExists failed: %v", err)
		return err
	}

	// ถ้า Repo บอกว่า exists == true แปลว่าระบบมี Super Admin อยู่แล้ว
	if exists {
		log.Printf("[WARN] AuthService.RegisterSystemAdminService - Attempted to create duplicate Super Admin. ")
		return errors.New("cannot create: a Super Admin account already exists in the system.")
	}

	hashedPassword, err := utils.HashedPassword(strings.TrimSpace(req.Password))

	if err != nil {
		log.Printf("[ERROR] AuthService.RegisterSystemAdminService - Hashing password failed: %v", err)
		return err
	}

	newUser := models.User{
		FirstName: firstName,
		LastName: lastName,
		Email: strings.TrimSpace(req.Email),
		Password: string(hashedPassword),
		ImageUrl: nil,
		SystemRole: "SYSTEM_ADMIN",
		PrefixID: req.PrefixID,
		IsActive: true,
		CreatedBy: uuid.Nil,
	}

	err = service.repo.CreateUserSystemAdmin(&newUser)
	if err != nil {
		log.Printf("[ERROR] AuthService.RegisterSystemAdminService - CreateUser in DB failed Error: %v", err)
		return err
	}

	return nil
} 

func (service *AuthService) RegisterUserService(req *authdto.RegisterUserRequest, userID uuid.UUID, storeID uuid.UUID) error {
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
			log.Printf("[RegisterUser Service WARN] File validation failed for user %s: %v", req.Email, err)
			return err
		}

        for _, file := range req.Files {
            // 💡 เรียกใช้ฟังก์ชันที่ปรับปรุงใหม่ จะได้ข้อมูลกลับมาครบ 3 อย่าง
            res, err := utils.UploadToCloudinary(file)
            if err != nil {
				log.Printf("[RegisterUser Service ERROR] Cloudinary upload crash for user %s: %v", req.Email, err)
                return err 
            }
            uploadResults = append(uploadResults, res)
        }
    }

	// ดักจับค่ารูปภาพในรูปแบบ Pointer ป้องกันสตริงว่างหลุดลง DB
    var imageUrl, originalName, fileName *string
    if len(uploadResults) > 0 {
        targetFile := uploadResults[0]
		imageUrl = &targetFile.SecureURL
		originalName = &targetFile.OriginalName
		fileName = &targetFile.CloudName
    }

	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		log.Printf("[RegisterUser Service ERROR] Password hashing failed for user %s: %v", req.Email, err)
		return err
	}

	userData := models.User{
		FirstName: firstName,
		LastName: lastName,
		Email: strings.TrimSpace(req.Email),
		Password: string(hashedPassword),
		ImageName: fileName,
		ImageOriginalName: originalName,
		ImageUrl: imageUrl,
		SystemRole: "USER",
		PrefixID: req.PrefixID,
		IsActive: true,
		CreatedBy: userID,
	}

	userStoreData := models.UserStore{
        StoreID:  storeID, 
        RoleID:   req.RoleID,
        IsActive: true,
    }

	err = service.repo.CreateUser(&userData, &userStoreData)
	if err != nil {
		log.Printf("[RegisterUser Service DATABASE ERROR] Transaction failed for user %s: %v", req.Email, err)
		return err
	}

	return nil
}

func (service *AuthService) ValidateRefreshTokenService(refreshTokenStr string) (bool, error) {
	refreshToken, isBlank := utils.IsBlank(refreshTokenStr)
	if isBlank { return false, utils.NewBadRequestError("refresh token is required") }

	hashedToken := utils.HashToken(refreshToken)

	 isValid ,err := service.repo.CheckRefreshTokenValid(hashedToken)

	 if err != nil {
		return false, err 
	 }
	
	 return isValid, nil
}

func (service *AuthService) ValidatePermissionService(userIDStr string, storeIDStr string) (*models.UserStore, error) {
	userIDS, isBlank := utils.IsBlank(userIDStr)
	if isBlank { 
		log.Printf("[Permission Service WARN] Validation failed: userID is empty")
		return nil, utils.NewBadRequestError("user id is required.") 
	}

	storeIDS, isBlank := utils.IsBlank(storeIDStr)
	if isBlank { 
		log.Printf("[Permission Service WARN] Validation failed: storeID is empty")
		return nil, utils.NewBadRequestError("store id is required.") 
	}

	userID, err := uuid.Parse(userIDS)
    if err != nil {
		log.Printf("[Permission Service WARN] Invalid userID UUID format submitted: '%s'", userIDS)
		return nil, utils.NewBadRequestError("invalid user id format.")
    }

	storeID, err := uuid.Parse(storeIDS)
    if err != nil {
		log.Printf("[Permission Service WARN] Invalid storeID UUID format submitted: '%s'", storeIDS)
		return nil, utils.NewBadRequestError("invalid store id format.")
    }
	
	userStore ,err := service.repo.CheckPermission(userID, storeID)

	if err != nil {
		log.Printf("[Permission Service ERROR] Permission check failed in DB for UserID: %s, StoreID: %s, Error: %v", userID.String(), storeID.String(), err)
		return nil, err 
	}

	return userStore, nil
}

func (service *AuthService) LoginService(req *authdto.LoginRequest, userAgent string) (string, string, *models.User, error) {
	email, isBlank := utils.IsBlank(req.Email)
    if isBlank { return "", "", nil, utils.NewBadRequestError("email is required") }

	client, isBlank := utils.IsBlank(req.Client)
	if isBlank { return "", "", nil, utils.NewBadRequestError("client is required") }

	password, isBlank := utils.IsBlank(req.Password)
	if isBlank { return "", "", nil, utils.NewBadRequestError("password is required")}

	user, err := service.repo.FindUserByEmail(email, "LOGIN")

	if err != nil {
		log.Printf("[WARN] AuthService.LoginService - Login failed: User not found or account is inactive. Input Email: %s", email)
		return "", "", nil, errors.New("invalid email or password")
	}

	err = utils.ComparePassword(user.Password, password)
	
	if err != nil {
		log.Printf("[WARN] AuthService.LoginService - Login failed: Password mismatch for UserID: %s, Email: %s", user.ID, email)
		return "", "", nil, errors.New("invalid email or password")
	}
	
	timeAccessTokenStr := os.Getenv("TIME_ACC_TOKEN")
	timeRefreshTokenStr := os.Getenv("TIME_REF_TOKEN")
	
	if timeAccessTokenStr == "" || timeRefreshTokenStr == "" {
		log.Println("[CRITICAL] AuthService.LoginService - Infrastructure configuration error: TIME_ACC_TOKEN or TIME_REF_TOKEN is missing in environment variables")
		return "", "", nil, errors.New("internal server error: security configuration missing")
	}

	// แปลงจาก string เป็น int 
	timeAccessToken, err := strconv.Atoi(timeAccessTokenStr)
	
	if err != nil {
		log.Printf("[ERROR] AuthService.LoginService - Configuration invalid: Failed to parse TIME_ACC_TOKEN value '%s' to integer: %v", timeAccessTokenStr, err)
		return "", "", nil, errors.New("internal server error: invalid security configuration")
	}

	timeRefreshToken, err := strconv.Atoi(timeRefreshTokenStr)
	
	if err != nil {
		log.Printf("[ERROR] AuthService.LoginService - Configuration invalid: Failed to parse TIME_REF_TOKEN value '%s' to integer: %v", timeRefreshTokenStr, err)
		return "", "", nil, errors.New("internal server error: invalid security configuration")
	}

	durationTimeAccessToken := time.Minute * time.Duration(timeAccessToken)
	durationTimeRefreshToken := time.Hour * time.Duration(timeRefreshToken)

	userIDStr := user.ID.String()

	accessToken, err := utils.GenerateJWT(userIDStr, user.SystemRole, durationTimeAccessToken)
	
	if err != nil {
		log.Printf("[ERROR] AuthService.LoginService - Security Fault: Failed to generate Access Token for UserID: %s. Error: %v", userIDStr, err)
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateJWT(userIDStr, user.SystemRole, durationTimeRefreshToken)

	if err != nil {
		log.Printf("[ERROR] AuthService.LoginService - Security Fault: Failed to generate Refresh Token for UserID: %s. Error: %v", userIDStr, err)
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	hashedToken := utils.HashToken(refreshToken)

	refreshTokenRecord := models.RefreshToken{
		UserID: user.ID,
		TokenHash: hashedToken,
		ClientType: client,
		DeviceInfo: &userAgent,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(durationTimeRefreshToken),
		CreatedBy: user.ID,
	}
	
	err = service.repo.CreateRefreshTokenRecord(&refreshTokenRecord)

	if err != nil {
		log.Printf("[ERROR] AuthService.LoginService - Database Failure: Cannot save refresh token record for UserID: %s to DB. Error: %v", userIDStr, err)
		return "", "", nil, fmt.Errorf("failed to persist session data: %w", err)
	} 

	return accessToken, refreshToken, user, nil
}

func (service *AuthService) LogoutService(userId uuid.UUID, rawRefreshToken string, allDevices bool) error {

	// 🚨 เคสที่ 1: สั่งลบทุกเครื่อง (ไม่สนใจ Token เครื่องปัจจุบัน)
	if allDevices {
		// ส่ง (userID, "") -> Repo จะไปใช้เงื่อนไขคัดเฉพาะอันที่ยังไม่หมดอายุของยูสเซอร์คนนี้
		err := service.repo.RevokeRefreshToken(userId, "");
		if err != nil {
			log.Printf("[AUTH][LOGOUT_SERVICE][REVOKE_ALL_FAILED] userID=%s error=%v", userId, err)
			return err
		}
		return nil
	}

	// 🚨 เคสที่ 2: สั่งลบเฉพาะเครื่องปัจจุบัน
	if rawRefreshToken == "" {
		log.Printf("[AUTH][LOGOUT_SERVICE][VALIDATION_FAILED] missing refresh token for current device logout userID=%s", userId)
		return errors.New("missing refresh token for current device logout")
	}
	
	// 🔒 ทำการแฮช Token ดิบให้เป็นค่า SHA-256 อยู่ในชั้นนี้ตามกฎธุรกิจ
	hashedToken := utils.HashToken(rawRefreshToken)

	// ส่ง (uuid.Nil, hashedToken) -> Repo จะเจาะจงทำลายใบนี้ทันที
	err := service.repo.RevokeRefreshToken(uuid.Nil, hashedToken)
	if err != nil {
		log.Printf("[AUTH][LOGOUT_SERVICE][REVOKE_SINGLE_FAILED] userID=%s error=%v", userId, err)
		return err
	}

	return nil
}

func (service *AuthService) ForgotPasswordService(req *authdto.ForgotPasswordRequest) error{
	email, isBlank := utils.IsBlank(req.Email)
    if isBlank { return utils.NewBadRequestError("email is required") }
	user, err := service.repo.FindUserByEmail(email, "FORGOT")

	if err != nil {
		return err
	}

	token := utils.RandomToken()
	
	timeResetPasswordStr := os.Getenv("TIME_RESET_PASSWORD")

	if timeResetPasswordStr == "" {
		log.Println("[ForgotPassword Service CRITICAL] Config Error: TIME_RESET_PASSWORD is missing in .env")
		// ถ้าลืมประกาศตัวแปรนี้ใน .env เลย ให้เบรกระบบทันที
		return errors.New("Missing in environment configuration")
	}

	// แปลงจาก string เป็น int 
	timeResetPassword, err := strconv.Atoi(timeResetPasswordStr)
	
	if err != nil {
		log.Printf("[ForgotPassword Service CRITICAL] Config Error: TIME_RESET_PASSWORD in .env must be a number: %v", err)
		return errors.New("internal server error: invalid security configuration format")
	}
	
	expiredAt := time.Now().Add(time.Duration(timeResetPassword) * time.Minute)

	resetData := models.ResetPassword{
		UserID: user.ID,
		Token: token,
		ExpiredAt: expiredAt,
		CreatedBy: user.ID,
	}

	err = service.repo.CreateResetPassword(&resetData)

	if err != nil {
		log.Printf("[ForgotPassword Service DATABASE ERROR] Failed to save token for user %s: %v", email, err)
		return err
	}

	logEmail := models.LogSendEmail{
        Recipient: user.Email,
        EmailType: "FORGOT_PASSWORD",
        Subject:   "Reset Password",
        Status:    "PENDING",
        CreatedBy: user.ID, // ใช้ ID ยูสเซอร์เป็นคนสร้าง Log นี้
    }

	errLog := service.repo.CreateLogEmail(&logEmail) 
	if errLog != nil {
        log.Printf("[ForgotPassword Service ERROR] Cannot create email log record for %s: %v", email, errLog)
    }

	go func(logID uuid.UUID, emailToSend string, tokenToSend string, userID uuid.UUID) {
        
        // สั่งยิงเมลผ่าน SMTP Client (ท่อนี้จะบล็อกรอ I/O ประมาณ 1-3 วินาทีอยู่ในพื้นหลัง)
        err := utils.SendResetPasswordEmail(emailToSend, tokenToSend)
        
        if err != nil {
            log.Printf("[ForgotPassword Background MAIL ERROR] Failed to dispatch reset email to %s: %v", emailToSend, err)
            
            // เมลพัง -> อัปเดตสเตตัสใน DB เป็น FAILED พร้อมแนบสาเหตุการแครช
            errStr := err.Error()
            _ = service.repo.UpdateLogEmailStatus(logID, "FAILED", &errStr, userID)
            return
        }
        
        //ส่งเมลสำเร็จ -> อัปเดตสเตตัสใน DB เป็น SUCCESS สวยๆ ลิงก์ทำงานได้ปกติ
        log.Printf("[ForgotPassword Background MAIL SUCCESS] Reset link successfully dispatched to %s", emailToSend)
        _ = service.repo.UpdateLogEmailStatus(logID, "SUCCESS", nil, userID)

    }(logEmail.ID, user.Email, token, user.ID)

	return nil
}

func (service *AuthService) ResetPasswordService(req *authdto.ResetPasswordRequest) error{
	token, isBlank := utils.IsBlank(req.Token)
    if isBlank { 
        return utils.NewBadRequestError("Token is required.") 
    }
    
    newPassword, isBlank := utils.IsBlank(req.NewPassword)
    if isBlank { 
        return utils.NewBadRequestError("New password is required.") 
    }
	// ตรวจสอบตั๋ว (Token) ว่าถูกต้อง/ไม่หมดอายุ ไหม
	resetToken, err := service.repo.FindValidResetToken(token)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            log.Printf("[ResetPassword Service WARN] Attempted to use invalid or expired token: '%s'", token)
            return utils.NewBadRequestError("invalid or expired token")
        }
        log.Printf("[ResetPassword Service DATABASE ERROR] Error searching reset token: %v", err)
        return err
    }

	// แฮชรหัสผ่านใหม่ด้วย bcrypt (เหมือนตอนสมัครสมาชิก)
	hashedPassword, err := utils.HashedPassword(strings.TrimSpace(newPassword))
    if err != nil {
        log.Printf("[ResetPassword Service ERROR] Hashing password failed: %v", err)
        return err
    }

	// 3. สั่งบันทึกรหัสใหม่ลง User และสั่งติ๊กใช้ตั๋วใบนี้แล้ว
	err = service.repo.UpdatePasswordAndRevokeToken(resetToken.UserID, string(hashedPassword), resetToken.ID)
    if err != nil {
        log.Printf("[ResetPassword Service DATABASE ERROR] Transaction execution failed for UserID %s: %v", resetToken.UserID.String(), err)
        return err
    }

	return nil
}
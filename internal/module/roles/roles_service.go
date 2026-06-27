package roles

import (
	"errors"
	"pos-system-backend/internal/models"
	roledto "pos-system-backend/internal/module/roles/dto"
	"pos-system-backend/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoleRepositoryInterface interface {
	GetRoleIdByRoleName(roleName string) (uuid.UUID, error)
	GetAllRoles() ([]models.Role, error)
	CreateRole(roleData *models.Role) error
	UpdateRole(roleData *models.Role) error
}

type RolesService struct {
	repo RoleRepositoryInterface
}

func NewRolesService (repo RoleRepositoryInterface) *RolesService{
	return &RolesService{repo: repo}
}

func (service *RolesService) GetRoleIdByRoleNameService(roleName string) (uuid.UUID, error) {

	if roleName == "" {
		return uuid.Nil, errors.New("invalid input: role name cannot be empty")
	}
	
	roleID ,err := service.repo.GetRoleIdByRoleName(strings.TrimSpace(roleName))

	if err != nil {
		return uuid.Nil, err
	}

	return roleID, nil
}

func (service *RolesService) GetAllRolesService() ([]roledto.GetAllRoleResponse, error) {
	roles, err := service.repo.GetAllRoles()

	if err != nil {
        return nil, err
    }

	responseRoles := []roledto.GetAllRoleResponse{}

	for _, r := range roles {
        formatted := roledto.GetAllRoleResponse{
            ID:       r.ID,
            RoleName: r.RoleName,
        }
        responseRoles = append(responseRoles, formatted)
    }

	return responseRoles, nil
}

func (service *RolesService) CreateRoleService(req *roledto.CreateRoleRequest, userId uuid.UUID) error {
	roleName, isBlank :=utils.IsBlank(req.RoleName)

	if isBlank {
		return errors.New("Role name is required")
	}

	trimDescription := strings.TrimSpace(req.Description)

	roleData := models.Role{
		RoleName: roleName,
		Description: &trimDescription,
		IsActive: true,
		CreatedBy: userId,
	}

	err := service.repo.CreateRole(&roleData)

	if err != nil {
		return err
	}

	return nil
}

func (service *RolesService) UpadateRoleService(req *roledto.UpdateRoleRequest, userId uuid.UUID) error {
	roleName, isBlank :=utils.IsBlank(req.RoleName)

	if isBlank {
		return utils.NewBadRequestError("Role name is required")
	}

	trimDescription := strings.TrimSpace(req.Description)

	roleData := models.Role{
		ID: req.ID,
		RoleName: roleName,
		Description: &trimDescription,
		IsActive: *req.IsActive,
		UpdatedAt: time.Now(),
	}
	err := service.repo.UpdateRole(&roleData)

	if err != nil {
		return err
	}

	return nil
}
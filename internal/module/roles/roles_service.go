package roles

import (
	"errors"

	"github.com/google/uuid"
)

type RolesService struct {
	repo *RolesRepository
}

func NewRolesService (repo *RolesRepository) *RolesService{
	return &RolesService{repo: repo}
}

func (service *RolesService) GetRoleIdByRoleNameService(roleName string) (uuid.UUID, error) {

	if roleName == "" {
		return uuid.Nil, errors.New("invalid input: role name cannot be empty")
	}
	
	roleID ,err := service.repo.GetRoleIdByRoleName(roleName)

	if err != nil {
		return uuid.Nil, err
	}

	return roleID, nil
}
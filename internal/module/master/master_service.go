package master

import (
	"pos-system-backend/internal/models"
	masterDto "pos-system-backend/internal/module/master/dto"
)

type masterRepositoryInterface interface {
	GetAllPrefix() ([]models.Prefix, error)
}

type MasterService struct {
	repo masterRepositoryInterface
}

func NewMasterService(repo masterRepositoryInterface) *MasterService {
	return &MasterService{repo: repo}
}

func (service *MasterService) GetAllPrefixService() ([]masterDto.GetAllPrefixResponse, error){
	prefixs, err := service.repo.GetAllPrefix()

	if err != nil {
        return nil, err
    }

	responsePrefix := []masterDto.GetAllPrefixResponse{}

	for _, value := range prefixs {
        formatted := masterDto.GetAllPrefixResponse{
            ID:       value.ID,
            PrefixName: value.TitleName,
        }
        responsePrefix = append(responsePrefix, formatted)
    }

	return responsePrefix, nil
}
package prefix

import (
	"pos-system-backend/internal/models"
	prefixDto "pos-system-backend/internal/module/prefix/dto"
)

type prefixRepositoryInterface interface {
	GetAllPrefix() ([]models.Prefix, error)
}

type PrefixService struct {
	repo prefixRepositoryInterface
}

func NewPrefixService(repo prefixRepositoryInterface) *PrefixService {
	return &PrefixService{repo: repo}
}

func (service *PrefixService) GetAllPrefixService() ([]prefixDto.GetAllPrefixResponse, error){
	prefixs, err := service.repo.GetAllPrefix()

	if err != nil {
        return nil, err
    }

	responsePrefix := []prefixDto.GetAllPrefixResponse{}

	for _, value := range prefixs {
        formatted := prefixDto.GetAllPrefixResponse{
            ID:       value.ID,
            PrefixName: value.PrefixName,
        }
        responsePrefix = append(responsePrefix, formatted)
    }

	return responsePrefix, nil
}
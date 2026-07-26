package master

import (
	"pos-system-backend/internal/models"
	masterDto "pos-system-backend/internal/module/master/dto"
)

type masterRepositoryInterface interface {
	GetAllPrefix(req *masterDto.GetAllPrefixRequest) ([]models.Prefix, error)
	GetAllUnit(req *masterDto.GetAllUnitRequest) ([]models.Unit, error)
	GetAllProvince(req *masterDto.GetAllProviceRequest) ([]models.Province, error)
	GetAllDistrict(req *masterDto.GetAllDistrictRequest) ([]models.District, error)
	GetAllSubdistrict(req *masterDto.GetAllSubdistrictRequest) ([]models.Subdistrict, error)
	GetAllPostcode(req *masterDto.GetAllPostcodeRequest) ([]models.Postcode, error)
	GetAllRole(req *masterDto.GetAllRoleRequest) ([]models.Role, error)
}

type MasterService struct {
	repo masterRepositoryInterface
}

func NewMasterService(repo masterRepositoryInterface) *MasterService {
	return &MasterService{repo: repo}
}

func (service *MasterService) GetAllPrefixService(req *masterDto.GetAllPrefixRequest) ([]masterDto.GetAllPrefixResponse, error){
	prefixs, err := service.repo.GetAllPrefix(req)

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

func (service *MasterService) GetAllUnitService(req *masterDto.GetAllUnitRequest) ([]masterDto.GetAllUnitResponse, error) {
	units, err := service.repo.GetAllUnit(req)

	if err != nil {
		return nil, err
	}

	responseUnit := []masterDto.GetAllUnitResponse{}

	for _, value := range units {
		formatted := masterDto.GetAllUnitResponse{
			ID: 	  value.ID,
			UnitName: value.UnitName,
		}
		responseUnit = append(responseUnit, formatted)
	}

	return responseUnit, nil
}

func (service *MasterService) GetAllProvinceService(req *masterDto.GetAllProviceRequest) ([]masterDto.GetAllProvinceResponse, error){
	province, err := service.repo.GetAllProvince(req)

	if err != nil {
		return nil, err
	}

	responseProvince := []masterDto.GetAllProvinceResponse{}

	for _, value := range province {
		formatted := masterDto.GetAllProvinceResponse{
			ID: value.ID,
			Province: value.ProvinceName,
		}
		responseProvince = append(responseProvince, formatted)
	}

	return responseProvince, nil
}

func (service *MasterService) GetAllDistrictService(req *masterDto.GetAllDistrictRequest) ([]masterDto.GetAllDistrictResponse, error){
	district, err := service.repo.GetAllDistrict(req)
	if err != nil {
		return nil, err
	}

	responseDistrict := []masterDto.GetAllDistrictResponse{}

	for _, value := range district {
		formatted := masterDto.GetAllDistrictResponse{
			ID: value.ID,
			District: value.DistrictName,
		}
		responseDistrict = append(responseDistrict, formatted)
	}

	return  responseDistrict, nil
}

func (service *MasterService) GetAllSubdistrictService(req *masterDto.GetAllSubdistrictRequest) ([]masterDto.GetAllSubdistrictResponse, error){
	subdistrict, err := service.repo.GetAllSubdistrict(req)
	if err != nil {
		return nil, err
	}

	responseSubdistrict := []masterDto.GetAllSubdistrictResponse{}

	for _, value := range subdistrict {
		formatted := masterDto.GetAllSubdistrictResponse{
			ID: value.ID,
			Subdistrict: value.SubdistrictName,
		}
		responseSubdistrict = append(responseSubdistrict, formatted)
	}

	return responseSubdistrict, nil
}

func (service *MasterService) GetAllPostcodeService(req *masterDto.GetAllPostcodeRequest) ([]masterDto.GetAllPostcodeResponse, error){
	postcode, err := service.repo.GetAllPostcode(req)
	if err != nil {
		return nil, err
	}

	responsePostcode := []masterDto.GetAllPostcodeResponse{}

	for _, value := range postcode {
		formatted := masterDto.GetAllPostcodeResponse{
			ID: value.ID,
			Postcode: value.Postcode,
		}
		responsePostcode = append(responsePostcode, formatted)
	}

	return responsePostcode, nil

}

func (service *MasterService) GetAllRoleService(req *masterDto.GetAllRoleRequest) ([]masterDto.GetAllRoleResponse, error){
	role, err := service.repo.GetAllRole(req)
	if err != nil {
		return nil, err
	}

	responseRole := []masterDto.GetAllRoleResponse{}

	for _, value := range role {
		formatted := masterDto.GetAllRoleResponse{
			ID: value.ID,
			RoleName: value.RoleName,
		}
		responseRole = append(responseRole, formatted)
	}

	return responseRole, nil
}
package master

import (
	"fmt"
	"pos-system-backend/internal/models"
	masterDto "pos-system-backend/internal/module/master/dto"
	"strings"

	"gorm.io/gorm"
)

type MasterRepository struct {
	db  *gorm.DB
}

func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{db: db}
}

func (repo *MasterRepository) GetAllPrefix(req *masterDto.GetAllPrefixRequest) ([]models.Prefix, error) {
	var prefix []models.Prefix

	// ไม่ใส่ delete IS NULL เพราะ Gorm แอยใส่ให้
	query := repo.db.Model(&models.Prefix{}).Where("is_active = ?", true)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("(title_name LIKE ?)", searchPattern)
	}

	err := query.Limit(20).Find(&prefix).Error
	if err != nil {
		return nil, err
	}

	// return ในรูปแบบ slice []models.Role เพราะ จะได้ [] ตอนไม่มีค่า และ 
	// slice ไม่ได้เก็บข้อมูลดิบทั้งหมดเอาไว้โดยตรง แต่มันคือโครงสร้างข้อมูลขนาดเล็ก (Header) ที่ประกอบด้วย 3 อย่างนี้เท่านั้น:
	// 1. Pointer ที่ชี้ไปยัง Array จริงๆ ในหน่วยความจำ (Underlying Array)
	// 2. Length (ความยาวปัจจุบัน)
	// 3. Capacity (ความจุสูงสุด)
	return prefix, nil
}

func (repo *MasterRepository) GetAllUnit(req *masterDto.GetAllUnitRequest) ([]models.Unit, error) {
	var unit []models.Unit

	query := repo.db.Model(&models.Unit{}).Where("is_active = ?", true)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("(unit_name LIKE ? OR unit_code LIKE ?)", searchPattern, searchPattern)
	}

	err := query.Limit(20).Find(&unit).Error
	if err != nil {
		return nil, err
	}

	return unit, nil
}

func (repo *MasterRepository) GetAllProvince(req *masterDto.GetAllProviceRequest) ([]models.Province, error){
	var province []models.Province

	query := repo.db.Model(&models.Province{}).Where("is_active = ?", true)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("province_name LIKE ?", searchPattern)
	}

	err := query.Limit(20).Find(&province).Error
	if err != nil {
		return nil, err
	}

	return province, nil
}

func (repo *MasterRepository) GetAllDistrict(req *masterDto.GetAllDistrictRequest) ([]models.District, error){
	var district []models.District

	query := repo.db.Model(&models.District{}).Where("is_active = ? AND province_id", true, req.Province_ID)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("district_name LIKE ?", searchPattern)
	}

	err := query.Limit(10).Find(&district).Error
	if err != nil {
		return nil, err
	}

	return district, nil
}

func (repo *MasterRepository) GetAllSubdistrict(req *masterDto.GetAllSubdistrictRequest) ([]models.Subdistrict, error){
	var subdistrict []models.Subdistrict

	query := repo.db.Model(&models.Subdistrict{}).Where("is_active = ? AND district_id = ?", true, req.District_ID)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("subdistrict_name LIKE ?", searchPattern)
	}

	err := query.Limit(10).Find(&subdistrict).Error
	if err != nil {
		return nil, err
	}

	return subdistrict, nil
}

func (repo *MasterRepository) GetAllPostcode(req *masterDto.GetAllPostcodeRequest) ([]models.Postcode, error){
	var postcode []models.Postcode

	query := repo.db.Model(&models.Postcode{}).
			 Joins("JOIN postcode_areas ON postcode_areas.postcode_id = postcodes.id AND postcode_areas.deleted_at IS NULL").
			 Where("postcode_areas.subdistrict_id = ?", req.Subdistrict_id).
			 Where("postcodes.is_active = ?", true)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("postcodes.postcode LIKE ?", searchPattern)
	}

	err := query.Distinct("postcodes.*").
        Order("postcodes.postcode ASC").
        Limit(10).
        Find(&postcode).Error

    if err != nil {
        return nil, err
    }

    return postcode, nil
}

func (repo *MasterRepository) GetAllRole(req *masterDto.GetAllRoleRequest) ([]models.Role, error){
	var role []models.Role

	query := repo.db.Model(&models.Role{}).Where("is_active = ?", true)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("(role_name ILIKE ?)", searchPattern)
	}

	err := query.Limit(10).Find(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}
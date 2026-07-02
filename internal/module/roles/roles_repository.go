package roles

import (
	"errors"
	"fmt"
	"log"
	"pos-system-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolesRepository struct {
	db *gorm.DB
}

func NewRolesRepository (db *gorm.DB) *RolesRepository {
	return &RolesRepository{db: db}
}

// Get
func (repo *RolesRepository) GetRoleIdByRoleName(roleName string) (uuid.UUID, error) {
	if roleName == "" {
		return uuid.Nil, errors.New("role name is required.")
	}

	var role models.Role

	err := repo.db.Select("id").Where("role_name = ? AND is_active = ?", roleName, true).First(&role).Error

	if err != nil {
		// ไม่เจอ Role นี้เลยในตาราง ให้มองว่ายังไม่มี super_admin และไม่ถือว่าระบบพัง
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return uuid.Nil, fmt.Errorf("repository error: role %s not found: %w", roleName, err) 
        }
		// เกิดข้อผิดพลาดอื่น เช่น DB ล่ม
		return uuid.Nil, err
	}

	return role.ID, nil
}

func (repo *RolesRepository) GetAllRoles() ([]models.Role, error) {
	var role []models.Role

	err := repo.db.Where("is_active = ? AND role_name != ? AND deleted_at IS NULL", true, "super_admin").Find(&role).Error

	if err != nil {
		return nil, err
	}
	// return ในรูปแบบ slice []models.Role เพราะ จะได้ [] ตอนไม่มีค่า และ 
	// slice ไม่ได้เก็บข้อมูลดิบทั้งหมดเอาไว้โดยตรง แต่มันคือโครงสร้างข้อมูลขนาดเล็ก (Header) ที่ประกอบด้วย 3 อย่างนี้เท่านั้น:
	// 1. Pointer ที่ชี้ไปยัง Array จริงๆ ในหน่วยความจำ (Underlying Array)
	// 2. Length (ความยาวปัจจุบัน)
	// 3. Capacity (ความจุสูงสุด) 
	return role, nil
}


// Create
func (repo *RolesRepository) CreateRole(roleData *models.Role) error {
	if roleData == nil {
		return errors.New("data role is required.")
	}

	err := repo.db.Create(roleData).Error

	if err != nil {
		return err
	}

	return nil
}

// update
func (repo *RolesRepository) UpdateRole(roleData *models.Role) error {
	if roleData == nil {
		log.Println("data role is required.")
		return errors.New("data role is required.")
	}

	// 🚀 ใช้ Updates ของ GORM ในการบันทึกข้อมูลลง Database
    // .Model() บอก GORM ว่าจะอัปเดตแถวไหน โดยอ้างอิงจาก ID ที่อยู่ในตัวแปร roleData
    // .Updates() จะทำการอัปเดตเฉพาะฟิลด์ที่มีการเปลี่ยนแปลงส่งเข้ามา
	// .Select("*") หรือระบุฟิลด์ เพื่อบอก GORM ว่า "ฟิลด์ไหนที่เป็น false หรือค่าว่าง ก็ให้อัปเดตลงไปด้วยนะ"
    result := repo.db.Model(&models.Role{}).Where("id = ?", roleData.ID).
	Select("role_name", "description", "is_active", "updated_at").
	Updates(roleData)
    
    // 1. เช็กว่าเกิด Error ระหว่างยิง SQL ไหม (เช่น Database ล่ม หรือ Constraint พัง)
    if result.Error != nil {
        return result.Error
    }

    // 2. เช็กว่ามีแถวไหนโดนอัปเดตจริงไหม (RowsAffected == 0 แปลว่าส่ง ID ผิดมา หาในเบสไม่เจอ)
    if result.RowsAffected == 0 {
        return errors.New("role not found or no changes made")
    }

    return nil
}
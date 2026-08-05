package master

import (
	"log"
	"net/http"
	masterDto "pos-system-backend/internal/module/master/dto"
	"pos-system-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type masterServiceInterface interface {
	GetAllPrefixService(req *masterDto.GetAllPrefixRequest) ([]masterDto.GetAllPrefixResponse, error)
    GetAllUnitService(req *masterDto.GetAllUnitRequest) ([]masterDto.GetAllUnitResponse,error)
    GetAllProvinceService(req *masterDto.GetAllProviceRequest) ([]masterDto.GetAllProvinceResponse, error)
    GetAllDistrictService(req *masterDto.GetAllDistrictRequest) ([]masterDto.GetAllDistrictResponse, error)
    GetAllSubdistrictService(req *masterDto.GetAllSubdistrictRequest) ([]masterDto.GetAllSubdistrictResponse, error)
    GetAllPostcodeService(req *masterDto.GetAllPostcodeRequest) ([]masterDto.GetAllPostcodeResponse, error)
    GetAllRoleService(userID uuid.UUID, systemRole string, req *masterDto.GetAllRoleRequest) ([]masterDto.GetAllRoleResponse, error)
}

type MasterController struct {
	service masterServiceInterface
}

func NewMasterController(service masterServiceInterface) *MasterController {
	return &MasterController{service: service}
}

func (masterCtrl *MasterController) GetAllPrefixController(c *gin.Context) {
    var req masterDto.GetAllPrefixRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_PREFIX][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

	perfixs, err := masterCtrl.service.GetAllPrefixService(&req)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch roles data",
            // ถ้าขึ้น Production อาจจะซ่อนไว้เพื่อความปลอดภัย
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active prefixs",
        "data":  perfixs, // หน้าบ้านจะแกะข้อมูลจาก Field "data" อันนี้ไปใช้ต่อ
    })
}

func (masterCtrl *MasterController) GetAllUnitController(c *gin.Context) {
    var req masterDto.GetAllUnitRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_UNIT][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

    units, err := masterCtrl.service.GetAllUnitService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch units data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active units",
        "data":  units,
    })
}

func (masterCtrl *MasterController) GetAllProvinceController(c *gin.Context) {
    var req masterDto.GetAllProviceRequest
    err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_PROVINCE][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

    provinces, err := masterCtrl.service.GetAllProvinceService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch provinces data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active provinces",
        "data":  provinces,
    })
}

func (masterCtrl *MasterController) GetAllDistrictController(c *gin.Context) {
    var req masterDto.GetAllDistrictRequest
    err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_DISTRICT][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

    districts ,err := masterCtrl.service.GetAllDistrictService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch districts data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active districts",
        "data":  districts,
    })
}

func (masterCtrl *MasterController) GetAllSubdistrictController(c *gin.Context) {
    var req masterDto.GetAllSubdistrictRequest
    err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_SUBDISTRICT][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	} 

    subdistricts, err := masterCtrl.service.GetAllSubdistrictService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch subdistricts data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active subdistricts",
        "data":  subdistricts,
    })
}

func (masterCtrl *MasterController) GetAllPostcodeController(c *gin.Context) {
    var req masterDto.GetAllPostcodeRequest
    err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_POSTCODE][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	} 

    postecodes, err := masterCtrl.service.GetAllPostcodeService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch postcodes data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active subdistricts",
        "data":  postecodes,
    })
}

func (masterCtrl *MasterController) GetAllRoleController(c *gin.Context) {
    userID, err := utils.GetFromCtx(c, "userID")
	if err != nil {
		log.Printf("[Get All Role ERROR] systemRole missing or empty in contex.")
        c.JSON(http.StatusUnauthorized, gin.H{
            "status": "error", 
            "message": "Unauthorized: user identity not found",
        })
        return
    }

    systemRole := c.GetString("systemRole")
	if systemRole == "" {
        log.Printf("[Get All Role ERROR] systemRole missing or empty in context")
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  "error",
            "message": "Unauthorized: system role not found",
        })
        return
    }
    var req masterDto.GetAllRoleRequest
    err = c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("[STORE][GET_ROLE][INVALID_REQUEST] path=%s error=%v", c.Request.URL.Path, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid request payload or missing required parameters",
        })
        return
	}

    roles, err := masterCtrl.service.GetAllRoleService(userID, systemRole, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch roles data",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active subdistricts",
        "data":  roles,
    })
}
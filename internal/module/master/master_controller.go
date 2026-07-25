package master

import (
	"log"
	"net/http"
	masterDto "pos-system-backend/internal/module/master/dto"

	"github.com/gin-gonic/gin"
)

type masterServiceInterface interface {
	GetAllPrefixService(req *masterDto.GetAllPrefixRequest) ([]masterDto.GetAllPrefixResponse, error)
    GetAllUnitService(req *masterDto.GetAllUnitRequest) ([]masterDto.GetAllUnitResponse,error)
    GetAllProvinceService(req *masterDto.GetAllProviceRequest) ([]masterDto.GetAllProvinceResponse, error)
    GetAllDistrictService(req *masterDto.GetAllDistrictRequest) ([]masterDto.GetAllDistrictResponse, error)
    GetAllSubdistrictService(req *masterDto.GetAllSubdistrictRequest) ([]masterDto.GetAllSubdistrictResponse, error)
    GetAllPostcodeService(req *masterDto.GetAllPostcodeRequest) ([]masterDto.GetAllPostcodeResponse, error)
}

type MasterController struct {
	service masterServiceInterface
}

func NewMasterController(service masterServiceInterface) *MasterController {
	return &MasterController{service: service}
}

func (MasterCtrl *MasterController) GetAllPrefixController(c *gin.Context) {
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

	perfixs, err := MasterCtrl.service.GetAllPrefixService(&req)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch roles data",
            "error":   err.Error(), // ถ้าขึ้น Production อาจจะซ่อนไว้เพื่อความปลอดภัย
        })
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active prefixs",
        "data":  perfixs, // หน้าบ้านจะแกะข้อมูลจาก Field "data" อันนี้ไปใช้ต่อ
    })
}

func (MasterCtrl *MasterController) GetAllUnitController(c *gin.Context) {
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

    units, err := MasterCtrl.service.GetAllUnitService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch units data",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active units",
        "data":  units,
    })
}

func (MasterCtrl *MasterController) GetAllProvinceController(c *gin.Context) {
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

    provinces, err := MasterCtrl.service.GetAllProvinceService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch provinces data",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active provinces",
        "data":  provinces,
    })
}

func (MasterCtrl *MasterController) GetAllDistrictController(c *gin.Context) {
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

    districts ,err := MasterCtrl.service.GetAllDistrictService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch districts data",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active districts",
        "data":  districts,
    })
}

func (MasterCtrl *MasterController) GetAllSubdistrictController(c *gin.Context) {
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

    subdistricts, err := MasterCtrl.service.GetAllSubdistrictService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch subdistricts data",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active subdistricts",
        "data":  subdistricts,
    })
}

func (MasterCtrl *MasterController) GetAllPostcodeController(c *gin.Context) {
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

    postecodes, err := MasterCtrl.service.GetAllPostcodeService(&req)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "failed to fetch postcodes data",
            "error":   err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "successfully retrieved active subdistricts",
        "data":  postecodes,
    })
}
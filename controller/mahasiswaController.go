package controller

import (
	"net/http"

	"github.com/Lotsoo/GoDroidAPI/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MahasiswaController struct {
	DB *gorm.DB
}

func NewMahasiswaController(db *gorm.DB) *MahasiswaController {
	return &MahasiswaController{DB: db}
}

func (controller *MahasiswaController) CreateMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa

	if err := c.ShouldBindJSON(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := controller.DB.Create(&mahasiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to save data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Student data saved successfully",
		"data":    mahasiswa,
	})
}

func (controller *MahasiswaController) GetAllMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa

	if err := controller.DB.Find(&mahasiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get all data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"result": mahasiswa,
	})

}

func (controller *MahasiswaController) GetMahasiswaByID(c *gin.Context) {
	id := c.Param("id")
	var mahasiswa models.Mahasiswa

	if err := controller.DB.Find(&mahasiswa, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mahasiswa,
	})

}

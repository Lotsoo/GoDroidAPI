package controller

import (
	"net/http"

	"github.com/Lotsoo/GoDroidAPI/models"
	"github.com/Lotsoo/GoDroidAPI/websocket" // 1. Import package websocket
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 2. Tambahkan Hub ke dalam struct
type MahasiswaController struct {
	DB  *gorm.DB
	Hub *websocket.Hub
}

// 3. Update constructor untuk menerima Hub
func NewMahasiswaController(db *gorm.DB, hub *websocket.Hub) *MahasiswaController {
	return &MahasiswaController{
		DB:  db,
		Hub: hub,
	}
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

	// 4. Broadcast pesan update setelah data berhasil dibuat
	controller.Hub.Broadcast <- []byte("update")

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Student data saved successfully",
		"data":    mahasiswa,
	})
}

// Tidak ada perubahan di sini karena hanya membaca data
func (controller *MahasiswaController) GetAllMahasiswa(c *gin.Context) {
	var mahasiswa []models.Mahasiswa

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

// Tidak ada perubahan di sini karena hanya membaca data
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

func (controller *MahasiswaController) UpdateMahasiswa(c *gin.Context) {
	id := c.Param("id")
	var mahasiswa models.Mahasiswa

	if err := controller.DB.Find(&mahasiswa, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data not found",
		})
		return
	}

	var updateInput models.Mahasiswa
	if err := c.ShouldBindJSON(&updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := controller.DB.Model(&mahasiswa).Updates(updateInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update student",
		})
		return
	}

	// 4. Broadcast pesan update setelah data berhasil diubah
	controller.Hub.Broadcast <- []byte("update")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data student updated successfully",
		"data":    mahasiswa,
	})
}

func (controller *MahasiswaController) DeleteMahasiswa(c *gin.Context) {
	id := c.Param("id")
	var mahasiswa models.Mahasiswa

	// Cek apakah data ada sebelum dihapus
	if result := controller.DB.First(&mahasiswa, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data not found"})
		return
	}

	// Hapus data
	if err := controller.DB.Delete(&mahasiswa, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to delete student"})
		return
	}

	// 4. Broadcast pesan update setelah data berhasil dihapus
	controller.Hub.Broadcast <- []byte("update")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data student deleted successfully",
	})
}

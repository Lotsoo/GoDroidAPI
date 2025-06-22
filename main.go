package main

import (
	"log"

	"github.com/Lotsoo/GoDroidAPI/config"
	"github.com/Lotsoo/GoDroidAPI/controller"
	"github.com/Lotsoo/GoDroidAPI/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file", err)
	}

	var err error

	// Init database
	db, err = config.InitDb()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// migrasi
	if err := db.AutoMigrate(&models.Mahasiswa{}); err != nil {
		log.Fatal("Failed to migrate: ", err)
	}
}

func main() {
	r := gin.Default()

	mahasiswaController := controller.NewMahasiswaController(db)

	api := r.Group("/api/v1")
	{
		api.POST("/mahasiswa", mahasiswaController.CreateMahasiswa)
		api.GET("/mahasiswa", mahasiswaController.GetAllMahasiswa)
		api.GET("/mahasiswa/:id", mahasiswaController.GetMahasiswaByID)
		api.PUT("/mahasiswa/:id", mahasiswaController.UpdateMahasiswa)
	}

	r.Run(":8080")

}

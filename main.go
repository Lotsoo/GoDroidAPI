package main

import (
	"log"
	"os"

	"github.com/Lotsoo/GoDroidAPI/config"
	"github.com/Lotsoo/GoDroidAPI/controller"
	"github.com/Lotsoo/GoDroidAPI/models"
	"github.com/Lotsoo/GoDroidAPI/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	hub *websocket.Hub
)

func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load .env file", err)
	}

	var err error

	// Init database
	db, err = config.InitDb()
	if err != nil {
		log.Println("Failed to connect to the database: ", err)
	}

	// migrasi
	if err := db.AutoMigrate(&models.Mahasiswa{}); err != nil {
		log.Println("Failed to migrate: ", err)
	}
}

func main() {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Length"},
		AllowCredentials: true,
	}))

	hub = websocket.NewHub()
	go hub.Run()

	mahasiswaController := controller.NewMahasiswaController(db, hub)

	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	api := r.Group("/api/v1")
	{
		api.POST("/mahasiswa", mahasiswaController.CreateMahasiswa)
		api.GET("/mahasiswa", mahasiswaController.GetAllMahasiswa)
		api.GET("/mahasiswa/:id", mahasiswaController.GetMahasiswaByID)
		api.PUT("/mahasiswa/:id", mahasiswaController.UpdateMahasiswa)
		api.DELETE("/mahasiswa/:id", mahasiswaController.DeleteMahasiswa)
	}

	log.Println("Server running on port 3000 in", gin.Mode(), "mode")
	r.Run(":3000")

}

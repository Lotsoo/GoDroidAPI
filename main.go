package main

import (
	"log"

	"github.com/Lotsoo/GoDroidAPI/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error: ", err)
		return
	}

	db, err := config.InitDb()
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
}

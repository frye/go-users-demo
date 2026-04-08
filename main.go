package main

import (
	"log"

	"userprofile-api/api"
	"userprofile-api/database"
)

func main() {
	db, err := database.InitDB("users.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := database.SeedDB(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	router := api.SetupRouter(db)

	log.Println("Starting server on :8080")
	router.Run(":8080")
}

package main

import (
	"log"
	"os"
	"tugas-praktikum-crud/database"
	"tugas-praktikum-crud/route"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    database.ConnectDB()

    app := fiber.New()

    route.SetupRoutes(app)

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000" 
    }
    log.Fatal(app.Listen(":" + port))
}
package main

import (
	"ai.unlocker.app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// load config
	port := os.Getenv("PORT")

	// init handler
	handler := NewHandler(services.NewGeminiService(), services.NewAssemblyAIService(), services.NewFileService())

	// start app
	app := fiber.New()

	app.Get("/speech-to-text", handler.SpeechToText)
	app.Post("/bot", handler.GenerateAnswer)
	app.Post("/video-to-text", handler.VideoToText)

	err = app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
		return
	}
}

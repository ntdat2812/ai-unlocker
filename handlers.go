package main

import (
	"ai.unlocker.app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

type Handler struct {
	gemini     *services.GeminiService
	assemblyAI *services.AssemblyAIService
	file       *services.FileService
}

func NewHandler(geminiClient *services.GeminiService, assemblyAI *services.AssemblyAIService, fileService *services.FileService) *Handler {
	return &Handler{gemini: geminiClient, assemblyAI: assemblyAI, file: fileService}
}

type Request struct {
	Text string `json:"text"`
}

func (h *Handler) GenerateAnswer(c *fiber.Ctx) error {

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	// ask gemini
	resp, err := h.gemini.Generate(c.Context(), req.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(fiber.Map{
		"result": resp,
	})
}

func (h *Handler) SpeechToText(c *fiber.Ctx) error {

	text, err := h.assemblyAI.Transcribe(c.Context(), "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(fiber.Map{
		"result": text,
	})
}

func (h *Handler) VideoToText(c *fiber.Ctx) error {

	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to upload video file",
		})
	}

	// Define the upload directory and make sure it exists
	outputDir := os.Getenv("OUTPUT_DIR")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create upload directory",
			})
		}
	}

	// Save the video file to the uploads folder
	videoFilePath := fmt.Sprintf("%s/%s", outputDir, file.Filename)
	err = c.SaveFile(file, videoFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save video file",
		})
	}

	// Define the output path for the audio file
	audioFilePath := fmt.Sprintf("%s/%s.mp3", outputDir, file.Filename[:len(file.Filename)-4])

	// Convert video to audio using FFmpeg
	err = h.file.ConvertVideoToAudio(videoFilePath, audioFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to convert video to audio",
		})
	}

	// Convert to text
	text, err := h.assemblyAI.Transcribe(c.Context(), audioFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to transcribe audio tot text",
		})
	}

	h.file.RemoveFiles(videoFilePath, audioFilePath)

	return c.JSON(fiber.Map{
		"result": text,
	})
}

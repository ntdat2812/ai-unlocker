package services

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"os"
)

type GeminiService struct {
	Key string
}

func NewGeminiService() *GeminiService {
	return &GeminiService{Key: os.Getenv("GEMINI_KEY")}
}

func (g *GeminiService) Generate(ctx context.Context, input string) (genai.Part, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_KEY")))
	if err != nil {
		log.Errorf("failed to initialize gemini client: %v", err.Error())
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		log.Error("failed to generate result: %v", err.Error())
		return nil, err
	}

	return g.ExtractTextFromResponse(resp), nil
}

func (g *GeminiService) ExtractTextFromResponse(resp *genai.GenerateContentResponse) genai.Part {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("could not process ExtractTextFromResponse: %v", err)
		}
	}()

	return resp.Candidates[0].Content.Parts[0]
}

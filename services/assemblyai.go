package services

import (
	"context"
	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/gofiber/fiber/v2/log"
	"os"
)

type AssemblyAIService struct {
	client *aai.Client
}

func NewAssemblyAIService() *AssemblyAIService {
	return &AssemblyAIService{client: aai.NewClient(os.Getenv("ASSEMBLY_AI_KEY"))}
}

func (a *AssemblyAIService) Transcribe(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("failed to open file (%v), err: %v", path, err.Error())
		return "", err
	}
	defer f.Close()

	transcript, err := a.client.Transcripts.TranscribeFromReader(ctx, f, nil)
	if err != nil {
		log.Errorf("failed to transcript file (%v) to text, err: %v", path, err.Error())
		return "", err
	}

	return *transcript.Text, nil
}

package services

import (
	"context"
	"github.com/joho/godotenv"
	"testing"
)

func Test_AssemblyAITranscribe(t *testing.T) {

	// load env file
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal(err)
		return
	}

	s := NewAssemblyAIService()

	audioInputPath := "../uploads/videoplayback.mp3"

	text, err := s.Transcribe(context.Background(), audioInputPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(text)

}

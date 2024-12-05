package services

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"os/exec"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) ConvertVideoToAudio(videoFilePath string, audioOutputPath string) error {
	// check if videoFilePath existed
	if _, err := os.Stat(videoFilePath); os.IsNotExist(err) {
		log.Errorf("Input video file does not exist: %v", err.Error())
		return err
	}

	// exec ffmpeg
	cmd := exec.Command("ffmpeg", "-i", videoFilePath, "-vn", "-acodec", "mp3", audioOutputPath)

	// Run the command and get the output
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert video to audio: %v", err.Error())
	}

	return nil
}

func (f *FileService) RemoveFiles(paths ...string) {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			log.Errorf("Failed to remove file(%v) err: %v", path, err.Error())
		}
	}

	return
}

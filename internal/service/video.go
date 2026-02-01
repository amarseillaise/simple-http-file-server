package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/amarseillaise/simple-http-file-server/internal/storage"
	"github.com/amarseillaise/simple-http-file-server/pkg/config"
)

const ytdlp = "yt-dlp"

var ErrReelDoesNotExistOrYtdlpBroken = errors.New("reel doesn't exist or yt-dlp error")

type VideoService struct {
	storage    *storage.FileSystem
	downloader VideoDownloader
}

type VideoDownloader interface {
	Download(shortcode string) (err error)
}

type Downloader struct{}

func (m *Downloader) Download(shortcode string) error {
	err := executeCMD(shortcode)
	return err
}

func executeCMD(shortcode string) error {
	args := getScriptArgs(shortcode)
	cmd := exec.Command(ytdlp, args...)

	conf := config.Load()
	outputDir := filepath.Join(conf.ContentDir, shortcode)
	os.MkdirAll(outputDir, 0755)

	res, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing yt-dlp: %v\nOutput: %s", err, res)
		return err
	}
	return nil
}

func getScriptArgs(shortcode string) []string {
	conf := config.Load()
	url := fmt.Sprintf("https://www.instagram.com/reel/%s/", shortcode)
	outputPath := fmt.Sprintf("%s/%s/video.%%(ext)s", conf.ContentDir, shortcode)
	descriptionPath := fmt.Sprintf("%s/%s/description.txt", conf.ContentDir, shortcode)

	return []string{
		url,
		"-o", outputPath,
		"--print-to-file", "description", descriptionPath,
		"--no-warnings",
		"--quiet",
		"--no-playlist",
		"--cookies", "./cookies.txt",
		"--format", "bestvideo[height<=1280][width<=720][vcodec^=avc1]+bestaudio/best[height<=1280][width<=720][vcodec^=avc1]/best",
	}
}

func NewVideoService(storage *storage.FileSystem, downloader VideoDownloader) *VideoService {
	return &VideoService{
		storage:    storage,
		downloader: downloader,
	}
}

func (s *VideoService) validateShortcode(shortcode string) error {
	if err := s.storage.ValidateShortcode(shortcode); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

func (s *VideoService) CreateReel(shortcode string) error {
	exists, err := s.CheckExists(shortcode)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return storage.ErrDirectoryExists
	}

	if err := s.storage.CreateDirectory(shortcode); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = s.downloader.Download(shortcode)
	if err != nil {
		s.storage.DeleteDirectory(shortcode)
		return ErrReelDoesNotExistOrYtdlpBroken
	}

	log.Printf("Successfully created video entry for shortcode: %s", shortcode)
	return nil
}

func (s *VideoService) CheckExists(shortcode string) (bool, error) {
	if err := s.validateShortcode(shortcode); err != nil {
		return false, err
	}
	return s.storage.DirectoryExists(shortcode), nil
}

func (s *VideoService) GetVideoPath(shortcode string) (string, error) {
	if err := s.validateShortcode(shortcode); err != nil {
		return "", err
	}
	return s.storage.GetVideoPath(shortcode)
}

func (s *VideoService) GetDescriptionPath(shortcode string) (string, error) {
	if err := s.validateShortcode(shortcode); err != nil {
		return "", err
	}
	return s.storage.GetDescriptionPath(shortcode)
}

func (s *VideoService) GetReelDescription(reelPath string) string {
	descriptionBytes, err := os.ReadFile(reelPath)
	if err != nil {
		return ""
	}
	return string(descriptionBytes)
}

func (s *VideoService) DeleteReel(shortcode string) error {
	if err := s.validateShortcode(shortcode); err != nil {
		return err
	}
	if err := s.storage.DeleteDirectory(shortcode); err != nil {
		return err
	}
	log.Printf("Successfully deleted video entry for shortcode: %s", shortcode)
	return nil
}

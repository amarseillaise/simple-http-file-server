package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ErrDirectoryExists   = errors.New("directory already exists")
	ErrDirectoryNotFound = errors.New("directory not found")
	ErrVideoNotFound     = errors.New("video file not found")
	ErrInvalidShortcode  = errors.New("invalid shortcode format")
)

var shortcodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type FileSystem struct {
	contentDir string
}

func NewFileSystem(contentDir string) (*FileSystem, error) {
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create content directory: %w", err)
	}

	return &FileSystem{
		contentDir: contentDir,
	}, nil
}

func (fs *FileSystem) ValidateShortcode(shortcode string) error {
	if shortcode == "" {
		return ErrInvalidShortcode
	}
	if len(shortcode) > 100 {
		return ErrInvalidShortcode
	}
	if !shortcodeRegex.MatchString(shortcode) {
		return ErrInvalidShortcode
	}
	return nil
}

func (fs *FileSystem) DirectoryExists(shortcode string) bool {
	dirPath := fs.getDirectoryPath(shortcode)
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (fs *FileSystem) CreateDirectory(shortcode string) error {
	if fs.DirectoryExists(shortcode) {
		return ErrDirectoryExists
	}

	dirPath := fs.getDirectoryPath(shortcode)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

func (fs *FileSystem) DeleteDirectory(shortcode string) error {
	if !fs.DirectoryExists(shortcode) {
		return ErrDirectoryNotFound
	}

	dirPath := fs.getDirectoryPath(shortcode)
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("failed to delete directory: %w", err)
	}

	return nil
}

func (fs *FileSystem) GetVideoPath(shortcode string) (string, error) {
	if !fs.DirectoryExists(shortcode) {
		return "", ErrDirectoryNotFound
	}

	videoPath := fs.getVideoPath(shortcode)
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "", ErrVideoNotFound
	}

	return videoPath, nil
}

func (fs *FileSystem) GetDescriptionPath(shortcode string) (string, error) {
	if !fs.DirectoryExists(shortcode) {
		return "", ErrDirectoryNotFound
	}

	descPath := fs.getDescriptionPath(shortcode)
	return descPath, nil
}

func (fs *FileSystem) getDirectoryPath(shortcode string) string {
	return filepath.Join(fs.contentDir, shortcode)
}

func (fs *FileSystem) getVideoPath(shortcode string) string {
	return filepath.Join(fs.contentDir, shortcode, "video.mp4")
}

func (fs *FileSystem) getDescriptionPath(shortcode string) string {
	return filepath.Join(fs.contentDir, shortcode, "description.txt")
}

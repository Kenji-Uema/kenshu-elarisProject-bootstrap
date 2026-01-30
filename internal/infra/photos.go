package infra

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Kenji-Uema/bootstrap/internal/config"
)

func BootstrapPhotos(photosVolumeConfig config.PhotosVolumeConfig) error {
	volumePath := filepath.Clean(photosVolumeConfig.Path)
	imagesDir := filepath.Join("resources", "images")

	if volumePath == "" || volumePath == "." || volumePath == ".." || volumePath == string(filepath.Separator) {
		return fmt.Errorf("refusing to use unsafe volume path: %q", photosVolumeConfig.Path)
	}

	if err := os.MkdirAll(volumePath, 0o755); err != nil {
		return fmt.Errorf("failed to create volume directory %q: %w", volumePath, err)
	}

	imageEntries, err := os.ReadDir(imagesDir)
	if err != nil {
		return fmt.Errorf("failed to read images directory %q: %w", imagesDir, err)
	}

	for _, entry := range imageEntries {
		if entry.IsDir() {
			continue
		}
		srcPath := filepath.Join(imagesDir, entry.Name())
		dstPath := filepath.Join(volumePath, entry.Name())
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy image from %q to %q: %w", srcPath, dstPath, err)
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(src *os.File) {
		err := src.Close()
		if err != nil {
			slog.Error("failed to close file", "error", err)
		}
	}(src)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			slog.Error("failed to close file", "error", err)
		}
	}(dst)

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	if info, err := os.Stat(srcPath); err == nil {
		return os.Chmod(dstPath, info.Mode())
	}
	return nil
}

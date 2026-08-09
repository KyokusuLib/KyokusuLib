package files

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

func UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, folder string, width, height int) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return "", errors.New("unsupported file format (jpg, png, webp only)")
	}

	srcImage, err := imaging.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	var dstImage = srcImage
	if width > 0 && height > 0 {
		dstImage = imaging.Fill(srcImage, width, height, imaging.Center, imaging.Lanczos)
	} else if width > 0 {
		dstImage = imaging.Resize(srcImage, width, 0, imaging.Lanczos)
	}

	newFilename := fmt.Sprintf("%s.jpg", uuid.New().String())

	baseUploadDir := "./uploads"
	uploadDir := filepath.Join(baseUploadDir, folder)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	savePath := filepath.Join(uploadDir, newFilename)

	err = imaging.Save(dstImage, savePath, imaging.JPEGQuality(80))
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	publicURL := fmt.Sprintf("/uploads/%s/%s", folder, newFilename)
	return publicURL, nil
}

// DeleteImage removes a previously uploaded image by its public URL.
func DeleteImage(publicURL string) error {
	if publicURL == "" {
		return nil
	}

	path := filepath.Join(".", publicURL)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image: %w", err)
	}
	return nil
}

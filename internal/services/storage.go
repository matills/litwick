package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	appconfig "github.com/matills/litwick/internal/config"
)

type StorageService struct {
	supabaseURL string
	serviceKey  string
	bucket      string
}

func NewStorageService(ctx context.Context) (*StorageService, error) {
	return &StorageService{
		supabaseURL: appconfig.AppConfig.SupabaseURL,
		serviceKey:  appconfig.AppConfig.SupabaseServiceKey,
		bucket:      appconfig.AppConfig.StorageBucket,
	}, nil
}

// UploadFile uploads a file to Supabase Storage and returns the public URL
func (s *StorageService) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	path := fmt.Sprintf("uploads/%s", uniqueFilename)

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload to Supabase Storage
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.supabaseURL, s.bucket, path)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.supabaseURL, s.bucket, path)

	return publicURL, nil
}

// GetPresignedURL generates a signed URL for downloading a file (for private buckets)
func (s *StorageService) GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	// For Supabase, if the bucket is public, we just return the public URL
	// If it's private, we'd need to create a signed URL
	url := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s?expiresIn=%d",
		s.supabaseURL,
		s.bucket,
		key,
		int(duration.Seconds()),
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to generate signed URL with status %d: %s", resp.StatusCode, string(body))
	}

	// Read and return the signed URL response
	body, _ := io.ReadAll(resp.Body)
	// In production, you'd want to parse the JSON response properly
	// For now, return the raw response body which contains the signed URL
	return string(body), nil
}

// DeleteFile deletes a file from Supabase Storage
func (s *StorageService) DeleteFile(ctx context.Context, key string) error {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.supabaseURL, s.bucket, key)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.serviceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

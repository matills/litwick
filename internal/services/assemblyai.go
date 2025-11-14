package services

import (
	"context"
	"fmt"
	"time"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/matills/litwick/internal/config"
)

type AssemblyAIService struct {
	client *aai.Client
}

func NewAssemblyAIService() *AssemblyAIService {
	client := aai.NewClient(config.AppConfig.AssemblyAIAPIKey)
	return &AssemblyAIService{client: client}
}

type TranscriptionResult struct {
	ID       string
	Text     string
	Status   string
	Duration int // in milliseconds
	Error    string
}

// UploadFile uploads a file to AssemblyAI and returns the upload URL
func (s *AssemblyAIService) UploadFile(ctx context.Context, fileURL string) (string, error) {
	// If the file is already a URL (S3), AssemblyAI can access it directly
	return fileURL, nil
}

// CreateTranscription creates a new transcription job
func (s *AssemblyAIService) CreateTranscription(ctx context.Context, audioURL string, language string) (*TranscriptionResult, error) {
	params := &aai.TranscriptOptionalParams{
		LanguageCode: aai.TranscriptLanguageCode(language),
	}

	transcript, err := s.client.Transcripts.TranscribeFromURL(ctx, audioURL, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create transcription: %w", err)
	}

	result := &TranscriptionResult{
		ID:     *transcript.ID,
		Status: string(transcript.Status),
	}

	if transcript.Text != nil {
		result.Text = *transcript.Text
	}

	if transcript.Error != nil {
		result.Error = *transcript.Error
	}

	if transcript.AudioDuration != nil {
		result.Duration = int(*transcript.AudioDuration)
	}

	return result, nil
}

// GetTranscription retrieves the status and result of a transcription
func (s *AssemblyAIService) GetTranscription(ctx context.Context, transcriptID string) (*TranscriptionResult, error) {
	transcript, err := s.client.Transcripts.Get(ctx, transcriptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcription: %w", err)
	}

	result := &TranscriptionResult{
		ID:     *transcript.ID,
		Status: string(transcript.Status),
	}

	if transcript.Text != nil {
		result.Text = *transcript.Text
	}

	if transcript.Error != nil {
		result.Error = *transcript.Error
	}

	if transcript.AudioDuration != nil {
		result.Duration = int(*transcript.AudioDuration)
	}

	return result, nil
}

// GetSRT exports the transcription as SRT subtitles
func (s *AssemblyAIService) GetSRT(ctx context.Context, transcriptID string) (string, error) {
	srt, err := s.client.Transcripts.GetSubtitles(ctx, transcriptID, aai.SubtitleFormat("srt"), nil)
	if err != nil {
		return "", fmt.Errorf("failed to get SRT: %w", err)
	}
	return string(srt), nil
}

// WaitForCompletion polls the transcription until it's completed or failed
func (s *AssemblyAIService) WaitForCompletion(ctx context.Context, transcriptID string, maxWait time.Duration) (*TranscriptionResult, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	timeout := time.After(maxWait)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("transcription timeout after %v", maxWait)
		case <-ticker.C:
			result, err := s.GetTranscription(ctx, transcriptID)
			if err != nil {
				return nil, err
			}

			switch result.Status {
			case "completed":
				return result, nil
			case "error":
				return result, fmt.Errorf("transcription failed: %s", result.Error)
			}
			// Continue polling for "queued" and "processing" statuses
		}
	}
}

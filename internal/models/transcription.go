package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TranscriptionStatus string

const (
	StatusPending    TranscriptionStatus = "pending"
	StatusProcessing TranscriptionStatus = "processing"
	StatusCompleted  TranscriptionStatus = "completed"
	StatusFailed     TranscriptionStatus = "failed"
)

type Transcription struct {
	ID               uuid.UUID           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID           `gorm:"type:uuid;not null;index" json:"user_id"`
	User             User                `gorm:"foreignKey:UserID" json:"-"`
	FileName         string              `gorm:"not null" json:"file_name"`
	FileURL          string              `gorm:"not null" json:"file_url"`
	FileSize         int64               `json:"file_size"`         // in bytes
	Duration         int                 `json:"duration"`          // in seconds
	Status           TranscriptionStatus `gorm:"default:'pending'" json:"status"`
	AssemblyAIID     string              `json:"assemblyai_id,omitempty"`
	TranscriptText   *string             `gorm:"type:text" json:"transcript_text,omitempty"`
	TranscriptJSON   *string             `gorm:"type:jsonb" json:"transcript_json,omitempty"` // Full AssemblyAI response
	SRTContent       *string             `gorm:"type:text" json:"srt_content,omitempty"`
	ErrorMessage     string              `json:"error_message,omitempty"`
	Language         string              `gorm:"default:'es'" json:"language"` // detected or specified language
	CreditsUsed      int                 `json:"credits_used"`                 // minutes of audio processed
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	CompletedAt      *time.Time          `json:"completed_at,omitempty"`
}

func (t *Transcription) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

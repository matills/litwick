package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SupabaseUserID   string    `gorm:"uniqueIndex;not null" json:"supabase_user_id"`
	Email            string    `gorm:"uniqueIndex;not null" json:"email"`
	CreditsRemaining int       `gorm:"default:300" json:"credits_remaining"` // 5 hours * 60 minutes = 300 minutes
	Plan             string    `gorm:"default:'free'" json:"plan"`           // free, pro, enterprise
	StripeCustomerID string    `json:"stripe_customer_id,omitempty"`

	// Settings
	DefaultLanguage        string `gorm:"default:'es'" json:"default_language"`                // Default transcription language
	DefaultExportFormat    string `gorm:"default:'srt'" json:"default_export_format"`          // txt, srt, vtt
	IncludeTimestamps      bool   `gorm:"default:true" json:"include_timestamps"`              // Include timestamps in exports
	DetectSpeakers         bool   `gorm:"default:true" json:"detect_speakers"`                 // Detect multiple speakers
	EmailNotifications     bool   `gorm:"default:true" json:"email_notifications"`             // Send email when transcription completes
	PromotionalEmails      bool   `gorm:"default:false" json:"promotional_emails"`             // Send promotional emails

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// HasCredits checks if user has enough credits for a given duration in minutes
func (u *User) HasCredits(minutes int) bool {
	return u.CreditsRemaining >= minutes
}

// DeductCredits removes credits from user account
func (u *User) DeductCredits(minutes int) {
	u.CreditsRemaining -= minutes
	if u.CreditsRemaining < 0 {
		u.CreditsRemaining = 0
	}
}

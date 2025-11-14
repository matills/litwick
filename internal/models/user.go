package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SupabaseUserID  string    `gorm:"uniqueIndex;not null" json:"supabase_user_id"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	CreditsRemaining int      `gorm:"default:300" json:"credits_remaining"` // 5 hours * 60 minutes = 300 minutes
	Plan            string    `gorm:"default:'free'" json:"plan"`           // free, pro, enterprise
	StripeCustomerID string   `json:"stripe_customer_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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

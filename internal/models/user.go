package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SupabaseUserID        string     `gorm:"uniqueIndex;not null" json:"supabase_user_id"`
	Email                 string     `gorm:"uniqueIndex;not null" json:"email"`
	CreditsRemaining      int        `gorm:"default:120" json:"credits_remaining"`        // Free: 120 minutes/month
	Plan                  string     `gorm:"default:'free'" json:"plan"`                  // free, pro, team
	SubscriptionStatus    string     `gorm:"default:'active'" json:"subscription_status"` // active, cancelled, expired
	MercadoPagoCustomerID string     `json:"mercadopago_customer_id,omitempty"`
	SubscriptionID        string     `json:"subscription_id,omitempty"`
	BillingCycleResetAt   *time.Time `json:"billing_cycle_reset_at,omitempty"`   // When credits reset
	TeamID                *uuid.UUID `gorm:"type:uuid" json:"team_id,omitempty"` // For Team plan
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
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

// PlanLimits represents limits for each plan
type PlanLimits struct {
	MonthlyMinutes      int // Total minutes per month
	MaxMinutesPerFile   int // Max minutes per single file
	MaxFileUploadsMonth int // Max file uploads per month
}

// GetPlanLimits returns the limits for the user's current plan
func (u *User) GetPlanLimits() PlanLimits {
	switch u.Plan {
	case "pro":
		return PlanLimits{
			MonthlyMinutes:      1800, // 30 hours
			MaxMinutesPerFile:   180,  // 3 hours
			MaxFileUploadsMonth: 100,
		}
	case "team":
		return PlanLimits{
			MonthlyMinutes:      3000, // 50 hours per person
			MaxMinutesPerFile:   480,  // 8 hours
			MaxFileUploadsMonth: 200,
		}
	default: // free
		return PlanLimits{
			MonthlyMinutes:      120, // 2 hours
			MaxMinutesPerFile:   15,  // 15 minutes
			MaxFileUploadsMonth: 50,
		}
	}
}

// GetMonthlyCredits returns the total credits a user should receive per month
func (u *User) GetMonthlyCredits() int {
	return u.GetPlanLimits().MonthlyMinutes
}

// CanUploadFile checks if user can upload a file based on their plan
func (u *User) CanUploadFile(durationMinutes int) (bool, string) {
	limits := u.GetPlanLimits()

	// Check if file duration exceeds plan limit
	if durationMinutes > limits.MaxMinutesPerFile {
		return false, "file exceeds maximum duration for your plan"
	}

	// Check if user has enough credits
	if !u.HasCredits(durationMinutes) {
		return false, "insufficient credits"
	}

	return true, ""
}

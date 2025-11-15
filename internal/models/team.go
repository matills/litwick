package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID                        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name                      string    `gorm:"not null" json:"name"`
	OwnerID                   uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	MemberCount               int       `gorm:"default:1" json:"member_count"`
	MercadoPagoSubscriptionID string    `json:"mercadopago_subscription_id,omitempty"`
	SubscriptionStatus        string    `gorm:"default:'active'" json:"subscription_status"` // active, cancelled, expired
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// GetTotalMonthlyCredits returns total credits for the team based on member count
func (t *Team) GetTotalMonthlyCredits() int {
	return t.MemberCount * 3000 // 3000 minutes per member
}

// GetPriceUSD returns monthly price in USD
func (t *Team) GetPriceUSD() int {
	return t.MemberCount * 20 // $20 per member
}

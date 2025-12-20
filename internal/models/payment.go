package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentApproved  PaymentStatus = "approved"
	PaymentRejected  PaymentStatus = "rejected"
	PaymentCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID                   uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID               uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	User                 User          `gorm:"foreignKey:UserID" json:"-"`
	MercadoPagoPaymentID *string       `gorm:"index" json:"mercadopago_payment_id,omitempty"`
	PreferenceID         *string       `json:"preference_id,omitempty"`
	Status               PaymentStatus `gorm:"default:'pending'" json:"status"`
	Amount               float64       `json:"amount"`
	Currency             string        `gorm:"default:'ARS'" json:"currency"`
	CreditsAmount        int           `json:"credits_amount"`
	PackageName          string        `json:"package_name"`
	PaymentMethod        *string       `json:"payment_method,omitempty"`
	PaymentDetails       *string       `gorm:"type:jsonb" json:"payment_details,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	CompletedAt          *time.Time    `json:"completed_at,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type CreditPackage struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Credits     int     `json:"credits"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Popular     bool    `json:"popular"`
	Discount    int     `json:"discount"`
}

func GetCreditPackages() []CreditPackage {
	return []CreditPackage{
		{
			ID:          "basic",
			Name:        "Básico",
			Description: "Perfecto para empezar",
			Credits:     120,
			Price:       5,
			Currency:    "USD",
			Popular:     false,
			Discount:    0,
		},
		{
			ID:          "standard",
			Name:        "Estándar",
			Description: "Ideal para uso regular",
			Credits:     300,
			Price:       10,
			Currency:    "USD",
			Popular:     true,
			Discount:    17,
		},
		{
			ID:          "premium",
			Name:        "Premium",
			Description: "Para usuarios frecuentes",
			Credits:     600,
			Price:       18,
			Currency:    "USD",
			Popular:     false,
			Discount:    25,
		},
		{
			ID:          "max",
			Name:        "Max",
			Description: "Máxima capacidad",
			Credits:     1500,
			Price:       40,
			Currency:    "USD",
			Popular:     false,
			Discount:    33,
		},
	}
}

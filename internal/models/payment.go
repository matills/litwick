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
	MercadoPagoPaymentID string        `gorm:"index" json:"mercadopago_payment_id,omitempty"` // ID del pago en MercadoPago
	PreferenceID         string        `json:"preference_id,omitempty"`                        // ID de la preferencia de pago
	Status               PaymentStatus `gorm:"default:'pending'" json:"status"`
	Amount               float64       `json:"amount"`           // Monto en pesos/dólares
	Currency             string        `gorm:"default:'ARS'" json:"currency"` // ARS, USD, etc.
	CreditsAmount        int           `json:"credits_amount"`   // Cantidad de créditos (minutos)
	PackageName          string        `json:"package_name"`     // Nombre del paquete comprado
	PaymentMethod        string        `json:"payment_method,omitempty"` // Método de pago usado
	PaymentDetails       *string       `gorm:"type:jsonb" json:"payment_details,omitempty"` // Detalles completos del pago (nullable)
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

// CreditPackage represents a package of credits that can be purchased
type CreditPackage struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Credits     int     `json:"credits"`      // Minutes of transcription
	Price       float64 `json:"price"`        // Price in local currency
	Currency    string  `json:"currency"`     // ARS, USD, etc.
	Popular     bool    `json:"popular"`      // Highlight as popular
	Discount    int     `json:"discount"`     // Discount percentage (0-100)
}

// GetCreditPackages returns available credit packages
func GetCreditPackages() []CreditPackage {
	return []CreditPackage{
		{
			ID:          "starter",
			Name:        "Starter",
			Description: "Perfecto para empezar",
			Credits:     120,  // 2 horas
			Price:       2999, // $2999 ARS
			Currency:    "ARS",
			Popular:     false,
			Discount:    0,
		},
		{
			ID:          "professional",
			Name:        "Professional",
			Description: "Ideal para uso regular",
			Credits:     600,  // 10 horas
			Price:       12999, // $12999 ARS
			Currency:    "ARS",
			Popular:     true,
			Discount:    15,
		},
		{
			ID:          "business",
			Name:        "Business",
			Description: "Para equipos y empresas",
			Credits:     1800,  // 30 horas
			Price:       34999, // $34999 ARS
			Currency:    "ARS",
			Popular:     false,
			Discount:    25,
		},
		{
			ID:          "enterprise",
			Name:        "Enterprise",
			Description: "Máximo rendimiento",
			Credits:     6000,  // 100 horas
			Price:       99999, // $99999 ARS
			Currency:    "ARS",
			Popular:     false,
			Discount:    35,
		},
	}
}

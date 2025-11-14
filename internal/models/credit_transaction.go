package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionDebit  TransactionType = "debit"  // Used credits
	TransactionCredit TransactionType = "credit" // Added credits
)

type CreditTransaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	User            User            `gorm:"foreignKey:UserID" json:"-"`
	TranscriptionID *uuid.UUID      `gorm:"type:uuid" json:"transcription_id,omitempty"`
	Type            TransactionType `gorm:"not null" json:"type"`
	Amount          int             `gorm:"not null" json:"amount"` // minutes
	BalanceBefore   int             `json:"balance_before"`
	BalanceAfter    int             `json:"balance_after"`
	Description     string          `json:"description"`
	CreatedAt       time.Time       `json:"created_at"`
}

func (ct *CreditTransaction) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	return nil
}

package services

import (
	"context"
	"fmt"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
	appconfig "github.com/matills/litwick/internal/config"
	"github.com/matills/litwick/internal/models"
)

type MercadoPagoService struct {
	client preference.Client
}

func NewMercadoPagoService() *MercadoPagoService {
	cfg, err := config.New(appconfig.AppConfig.MercadoPagoAccessToken)
	if err != nil {
		panic(fmt.Sprintf("Failed to create MercadoPago config: %v", err))
	}

	return &MercadoPagoService{
		client: preference.NewClient(cfg),
	}
}

// CreatePreference creates a payment preference in MercadoPago
func (s *MercadoPagoService) CreatePreference(ctx context.Context, pkg models.CreditPackage, userEmail string, paymentID string) (string, error) {
	// Create preference request
	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				ID:          pkg.ID,
				Title:       pkg.Name + " - " + pkg.Description,
				Description: pkg.Description,
				Quantity:    1,
				UnitPrice:   pkg.Price,
				CurrencyID:  pkg.Currency,
			},
		},
		Payer: &preference.PayerRequest{
			Email: userEmail,
		},
		BackURLs: &preference.BackURLsRequest{
			Success: appconfig.AppConfig.FrontendURL + "/credits?status=success",
			Failure: appconfig.AppConfig.FrontendURL + "/credits?status=failure",
			Pending: appconfig.AppConfig.FrontendURL + "/credits?status=pending",
		},
		AutoReturn:          "approved",
		ExternalReference:   paymentID, // Our internal payment ID
		NotificationURL:     appconfig.AppConfig.WebhookURL + "/api/payments/webhook",
		StatementDescriptor: "LITWICK - Créditos",
	}

	// Create preference
	resp, err := s.client.Create(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to create preference: %w", err)
	}

	return resp.InitPoint, nil
}

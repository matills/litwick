package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

type PreferenceResponse struct {
	InitPoint    string
	PreferenceID string
}

func (s *MercadoPagoService) CreatePreference(ctx context.Context, pkg models.CreditPackage, userEmail string, paymentID string) (*PreferenceResponse, error) {
	baseURL := appconfig.AppConfig.FrontendURL + "/credits"

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
			Success: baseURL + "?payment_status=success",
			Failure: baseURL + "?payment_status=failure",
			Pending: baseURL + "?payment_status=pending",
		},
		BinaryMode:          true,
		ExternalReference:   paymentID,
		NotificationURL:     appconfig.AppConfig.WebhookURL + "/api/payments/webhook",
		StatementDescriptor: "LITWICK - Créditos",
		Purpose:             "wallet_purchase",
	}

	// Log the request for debugging
	reqJSON, _ := json.MarshalIndent(request, "", "  ")
	log.Printf("Creating MercadoPago preference with request: %s", string(reqJSON))

	resp, err := s.client.Create(ctx, request)
	if err != nil {
		log.Printf("MercadoPago preference creation failed: %v", err)
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	// Log the response
	log.Printf("MercadoPago preference created successfully: ID=%s, InitPoint=%s", resp.ID, resp.InitPoint)

	return &PreferenceResponse{
		InitPoint:    resp.InitPoint,
		PreferenceID: resp.ID,
	}, nil
}

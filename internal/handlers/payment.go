package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	appconfig "github.com/matills/litwick/internal/config"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/middleware"
	"github.com/matills/litwick/internal/models"
	"github.com/matills/litwick/internal/services"
)

func GetCreditPackages(c *fiber.Ctx) error {
	packages := models.GetCreditPackages()
	return c.JSON(fiber.Map{
		"packages": packages,
	})
}

func CreatePayment(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	type CreatePaymentRequest struct {
		PackageID string `json:"package_id"`
	}

	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	packages := models.GetCreditPackages()
	var selectedPackage *models.CreditPackage
	for _, pkg := range packages {
		if pkg.ID == req.PackageID {
			selectedPackage = &pkg
			break
		}
	}

	if selectedPackage == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid package ID",
		})
	}

	payment := models.Payment{
		UserID:        user.ID,
		Status:        models.PaymentPending,
		Amount:        selectedPackage.Price,
		Currency:      selectedPackage.Currency,
		CreditsAmount: selectedPackage.Credits,
		PackageName:   selectedPackage.Name,
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create payment",
		})
	}

	mpService := services.NewMercadoPagoService()
	ctx := context.Background()

	prefResp, err := mpService.CreatePreference(ctx, *selectedPackage, user.Email, payment.ID.String())
	if err != nil {
		log.Printf("Failed to create MercadoPago preference: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create payment preference",
		})
	}

	// Save preference ID to payment
	payment.PreferenceID = &prefResp.PreferenceID
	if err := database.DB.Save(&payment).Error; err != nil {
		log.Printf("Failed to save payment: %v", err)
	}

	return c.JSON(fiber.Map{
		"payment_id":   payment.ID,
		"init_point":   prefResp.InitPoint,
		"preference_id": prefResp.PreferenceID,
	})
}

// verifyWebhookSignature verifies the MercadoPago webhook signature
func verifyWebhookSignature(c *fiber.Ctx) bool {
	// Get the webhook secret from config
	secret := appconfig.AppConfig.MercadoPagoWebhookSecret
	if secret == "" {
		log.Printf("WARNING: MercadoPago webhook secret not configured - skipping signature verification")
		return true // Allow in development, but log warning
	}

	// Get headers
	xSignature := c.Get("x-signature")
	xRequestID := c.Get("x-request-id")

	if xSignature == "" {
		log.Printf("Missing x-signature header")
		return false
	}

	// Get data.id from query params
	dataID := c.Query("data.id")
	if dataID == "" {
		log.Printf("Missing data.id query parameter")
		return false
	}

	// Parse x-signature to extract ts and v1
	parts := strings.Split(xSignature, ",")
	var ts, hash string

	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			if key == "ts" {
				ts = value
			} else if key == "v1" {
				hash = value
			}
		}
	}

	if ts == "" || hash == "" {
		log.Printf("Invalid x-signature format: missing ts or v1")
		return false
	}

	// Generate the manifest string
	manifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", dataID, xRequestID, ts)

	// Create HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(manifest))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Compare signatures
	isValid := hmac.Equal([]byte(expectedHash), []byte(hash))

	if !isValid {
		log.Printf("Webhook signature verification failed. Expected: %s, Got: %s", expectedHash, hash)
	}

	return isValid
}

func WebhookMercadoPago(c *fiber.Ctx) error {
	log.Printf("Received MercadoPago webhook notification")

	// Verify webhook signature
	if !verifyWebhookSignature(c) {
		log.Printf("Webhook signature verification failed - rejecting webhook")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	type WebhookData struct {
		Action      string `json:"action"`
		APIVersion  string `json:"api_version"`
		Data        struct {
			ID string `json:"id"`
		} `json:"data"`
		DateCreated string `json:"date_created"`
		ID          int64  `json:"id"`
		LiveMode    bool   `json:"live_mode"`
		Type        string `json:"type"`
		UserID      string `json:"user_id"`
	}

	var webhook WebhookData
	if err := c.BodyParser(&webhook); err != nil {
		log.Printf("Failed to parse webhook body: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Printf("Webhook verified - type=%s, action=%s, payment_id=%s", webhook.Type, webhook.Action, webhook.Data.ID)

	// Only process payment notifications
	if webhook.Type != "payment" {
		log.Printf("Ignoring non-payment webhook type: %s", webhook.Type)
		return c.SendStatus(fiber.StatusOK)
	}

	// Get payment details from MercadoPago
	cfg, err := config.New(appconfig.AppConfig.MercadoPagoAccessToken)
	if err != nil {
		log.Printf("Failed to create MercadoPago config: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	paymentClient := payment.NewClient(cfg)
	ctx := context.Background()

	// Convert payment ID from string to int
	paymentID, err := strconv.Atoi(webhook.Data.ID)
	if err != nil {
		log.Printf("Invalid payment ID format: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	mpPayment, err := paymentClient.Get(ctx, paymentID)
	if err != nil {
		log.Printf("Failed to get payment from MercadoPago: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Retrieved payment from MercadoPago: id=%s, status=%s, external_reference=%s",
		mpPayment.ID, mpPayment.Status, mpPayment.ExternalReference)

	// Find our payment record using external_reference (which is our payment ID)
	if mpPayment.ExternalReference == "" {
		log.Printf("Payment has no external_reference - cannot process")
		return c.SendStatus(fiber.StatusOK)
	}

	paymentUUID, err := uuid.Parse(mpPayment.ExternalReference)
	if err != nil {
		log.Printf("Invalid external_reference UUID: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var ourPayment models.Payment
	if err := database.DB.Where("id = ?", paymentUUID).First(&ourPayment).Error; err != nil {
		log.Printf("Payment not found in database: %s", paymentUUID)
		return c.SendStatus(fiber.StatusNotFound)
	}

	// Check if already processed to avoid double-processing
	if ourPayment.Status != models.PaymentPending {
		log.Printf("Payment already processed with status: %s", ourPayment.Status)
		return c.SendStatus(fiber.StatusOK)
	}

	// Update payment based on MercadoPago status
	mpPaymentID := fmt.Sprintf("%d", mpPayment.ID)
	now := time.Now()
	ourPayment.MercadoPagoPaymentID = &mpPaymentID

	switch mpPayment.Status {
	case "approved":
		log.Printf("Payment approved - adding %d credits to user %s", ourPayment.CreditsAmount, ourPayment.UserID)

		ourPayment.Status = models.PaymentApproved
		ourPayment.CompletedAt = &now

		// Get user and add credits
		var user models.User
		if err := database.DB.Where("id = ?", ourPayment.UserID).First(&user).Error; err != nil {
			log.Printf("Failed to find user: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Add credits to user
		balanceBefore := user.CreditsRemaining
		user.CreditsRemaining += ourPayment.CreditsAmount

		if err := database.DB.Save(&user).Error; err != nil {
			log.Printf("Failed to update user credits: %v", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Create credit transaction record
		transaction := models.CreditTransaction{
			UserID:        user.ID,
			Type:          models.TransactionCredit,
			Amount:        ourPayment.CreditsAmount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  user.CreditsRemaining,
			Description:   fmt.Sprintf("Compra de paquete %s (Webhook)", ourPayment.PackageName),
		}
		if err := database.DB.Create(&transaction).Error; err != nil {
			log.Printf("Failed to create credit transaction: %v", err)
		}

		log.Printf("Successfully added %d credits to user %s. New balance: %d",
			ourPayment.CreditsAmount, user.ID, user.CreditsRemaining)

	case "rejected", "cancelled":
		log.Printf("Payment %s - status: %s", mpPayment.Status, mpPayment.Status)
		if mpPayment.Status == "rejected" {
			ourPayment.Status = models.PaymentRejected
		} else {
			ourPayment.Status = models.PaymentCancelled
		}

	case "pending", "in_process", "in_mediation":
		log.Printf("Payment still pending - status: %s", mpPayment.Status)
		// Keep as pending, don't update status yet
		return c.SendStatus(fiber.StatusOK)

	default:
		log.Printf("Unknown payment status: %s", mpPayment.Status)
		return c.SendStatus(fiber.StatusOK)
	}

	// Save payment method if available
	if mpPayment.PaymentMethodID != "" {
		ourPayment.PaymentMethod = &mpPayment.PaymentMethodID
	}

	// Save payment details as JSON
	detailsMap := map[string]interface{}{
		"mercadopago_payment_id": mpPayment.ID,
		"status":                 mpPayment.Status,
		"status_detail":          mpPayment.StatusDetail,
		"payment_type":           mpPayment.PaymentTypeID,
		"payment_method":         mpPayment.PaymentMethodID,
		"transaction_amount":     mpPayment.TransactionAmount,
		"processed_at":           time.Now().Format(time.RFC3339),
	}
	detailsJSON, _ := json.Marshal(detailsMap)
	detailsStr := string(detailsJSON)
	ourPayment.PaymentDetails = &detailsStr

	// Save updated payment
	if err := database.DB.Save(&ourPayment).Error; err != nil {
		log.Printf("Failed to save payment: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("Webhook processed successfully for payment %s", ourPayment.ID)
	return c.SendStatus(fiber.StatusOK)
}

func GetPaymentHistory(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var payments []models.Payment
	if err := database.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&payments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch payment history",
		})
	}

	return c.JSON(fiber.Map{
		"payments": payments,
	})
}

func ProcessPaymentSuccess(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	paymentID := c.Query("payment_id")
	status := c.Query("status")
	externalReference := c.Query("external_reference")
	preferenceID := c.Query("preference_id")

	log.Printf("Payment success callback: payment_id=%s, status=%s, external_ref=%s",
		paymentID, status, externalReference)

	if externalReference == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing external reference",
		})
	}

	var payment models.Payment
	if err := database.DB.Where("id = ? AND user_id = ?", externalReference, user.ID).First(&payment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "payment not found",
		})
	}

	if payment.Status != models.PaymentPending {
		return c.JSON(fiber.Map{
			"payment": payment,
			"message": "payment already processed",
		})
	}

	now := time.Now()
	payment.MercadoPagoPaymentID = &paymentID
	payment.PreferenceID = &preferenceID

	if status == "approved" {
		payment.Status = models.PaymentApproved
		payment.CompletedAt = &now

		user.CreditsRemaining += payment.CreditsAmount
		if err := database.DB.Save(user).Error; err != nil {
			log.Printf("Failed to add credits to user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to add credits",
			})
		}

		transaction := models.CreditTransaction{
			UserID:        user.ID,
			Type:          models.TransactionCredit,
			Amount:        payment.CreditsAmount,
			BalanceBefore: user.CreditsRemaining - payment.CreditsAmount,
			BalanceAfter:  user.CreditsRemaining,
			Description:   fmt.Sprintf("Compra de paquete %s", payment.PackageName),
		}
		database.DB.Create(&transaction)

		log.Printf("Payment approved: user_id=%s, credits_added=%d", user.ID, payment.CreditsAmount)
	} else if status == "rejected" || status == "cancelled" {
		if status == "rejected" {
			payment.Status = models.PaymentRejected
		} else {
			payment.Status = models.PaymentCancelled
		}
	}

	if paymentID != "" {
		details := map[string]string{
			"payment_id":    paymentID,
			"status":        status,
			"preference_id": preferenceID,
		}
		detailsJSON, _ := json.Marshal(details)
		detailsStr := string(detailsJSON)
		payment.PaymentDetails = &detailsStr
	}

	if err := database.DB.Save(&payment).Error; err != nil {
		log.Printf("Failed to save payment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save payment",
		})
	}

	return c.JSON(fiber.Map{
		"payment": payment,
		"user": user,
	})
}

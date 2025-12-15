package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
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

	initPoint, err := mpService.CreatePreference(ctx, *selectedPackage, user.Email, payment.ID.String())
	if err != nil {
		log.Printf("Failed to create MercadoPago preference: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create payment preference",
		})
	}

	if err := database.DB.Save(&payment).Error; err != nil {
		log.Printf("Failed to save payment: %v", err)
	}

	return c.JSON(fiber.Map{
		"payment_id": payment.ID,
		"init_point": initPoint,
	})
}

func WebhookMercadoPago(c *fiber.Ctx) error {
	type WebhookData struct {
		Action string `json:"action"`
		APIVersion string `json:"api_version"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
		DateCreated string `json:"date_created"`
		ID int64 `json:"id"`
		LiveMode bool `json:"live_mode"`
		Type string `json:"type"`
		UserID string `json:"user_id"`
	}

	var webhook WebhookData
	if err := c.BodyParser(&webhook); err != nil {
		log.Printf("Failed to parse webhook: %v", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	log.Printf("Received MercadoPago webhook: type=%s, action=%s, id=%s", webhook.Type, webhook.Action, webhook.Data.ID)

	if webhook.Type != "payment" {
		return c.SendStatus(fiber.StatusOK)
	}

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
	payment.MercadoPagoPaymentID = paymentID
	payment.PreferenceID = preferenceID

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

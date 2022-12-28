package onfido

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Constants
const (
	WebhookSignatureHeader = "X-Signature"
	WebhookTokenEnv        = "ONFIDO_WEBHOOK_TOKEN"
)

// Webhook errors
var (
	ErrInvalidWebhookSignature = errors.New("invalid request, payload hash doesn't match signature")
	ErrMissingWebhookToken     = errors.New("webhook token not found in environmental variable")
)

// Webhook represents a webhook handler
type Webhook struct {
	Token                   string
	SkipSignatureValidation bool
}

// WebhookRequest represents an incoming webhook request from Onfido
type WebhookRequest struct {
	Payload struct {
		ResourceType string `json:"resource_type"`
		Action       string `json:"action"`
		Object       struct {
			ID                 string    `json:"id"`
			Status             string    `json:"status"`
			CompletedAtIso8601 time.Time `json:"completed_at_iso8601"`
			Href               string    `json:"href"`
		} `json:"object"`
		Resource struct {
			CreatedAt time.Time   `json:"created_at"`
			Error     interface{} `json:"error"`
			Link      struct {
				CompletedRedirectURL interface{} `json:"completed_redirect_url"`
				URL                  string      `json:"url"`
				Language             interface{} `json:"language"`
				ExpiresAt            interface{} `json:"expires_at"`
				ExpiredRedirectURL   interface{} `json:"expired_redirect_url"`
			} `json:"link"`
			UpdatedAt  time.Time `json:"updated_at"`
			WorkflowID string    `json:"workflow_id"`
			Status     string    `json:"status"`
			Output     struct {
				DateOfBirth            string `json:"date_of_birth"`
				DocumentIssuingCountry string `json:"document_issuing_country"`
				DocumentMediaIds       []struct {
					ID string `json:"id"`
				} `json:"document_media_ids"`
				DocumentNumber string `json:"document_number"`
				DocumentType   string `json:"document_type"`
				FirstName      string `json:"first_name"`
				LastName       string `json:"last_name"`
				SelfieMediaIds []struct {
					ID string `json:"id"`
				} `json:"selfie_media_ids"`
			} `json:"output"`
			DashboardURL      string        `json:"dashboard_url"`
			WorkflowVersionID int           `json:"workflow_version_id"`
			ID                string        `json:"id"`
			ApplicantID       string        `json:"applicant_id"`
			Reasons           []interface{} `json:"reasons"`
		} `json:"resource"`
	} `json:"payload"`
}

// NewWebhookFromEnv creates a new webhook handler using
// configuration from environment variables.
func NewWebhookFromEnv() (*Webhook, error) {
	token := os.Getenv(WebhookTokenEnv)
	if token == "" {
		return nil, ErrMissingWebhookToken
	}
	return NewWebhook(token), nil
}

// NewWebhook creates a new webhook handler
func NewWebhook(token string) *Webhook {
	return &Webhook{
		Token: token,
	}
}

// ValidateSignature validates the request body against the signature header.
func (wh *Webhook) ValidateSignature(body []byte, signature string) error {
	mac := hmac.New(sha1.New, []byte(wh.Token))
	if _, err := mac.Write(body); err != nil {
		return err
	}

	sig, err := hex.DecodeString(signature)
	if err != nil || !hmac.Equal(sig, mac.Sum(nil)) {
		return ErrInvalidWebhookSignature
	}

	return nil
}

// ParseFromRequest parses the webhook request body and returns
// it as WebhookRequest if the request signature is valid.
func (wh *Webhook) ParseFromRequest(req *http.Request) (*WebhookRequest, error) {
	signature := req.Header.Get(WebhookSignatureHeader)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if !wh.SkipSignatureValidation {
		if err := wh.ValidateSignature(body, signature); err != nil {
			return nil, err
		}
	}

	var wr WebhookRequest
	if err := json.Unmarshal(body, &wr); err != nil {
		return nil, err
	}

	return &wr, nil
}

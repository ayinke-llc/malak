package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ayinke-llc/malak"
	malak_mocks "github.com/ayinke-llc/malak/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateStripeSignature(payload []byte, secret string, timestamp int64) string {
	// Format the payload as Stripe expects: timestamp.payload
	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))

	signature := hex.EncodeToString(mac.Sum(nil))

	// Format the full signature header as Stripe expects: t=timestamp,v1=signature
	return fmt.Sprintf("t=%d,v1=%s", timestamp, signature)
}

func TestStripeHandler_HandleWebhook(t *testing.T) {
	for _, v := range generateStripeWebhookTestTable() {
		t.Run(v.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userRepo := malak_mocks.NewMockUserRepository(controller)
			planRepo := malak_mocks.NewMockPlanRepository(controller)
			workspaceRepo := malak_mocks.NewMockWorkspaceRepository(controller)
			preferencesRepo := malak_mocks.NewMockPreferenceRepository(controller)
			billingClient := malak_mocks.NewMockClient(controller)
			queueRepo := malak_mocks.NewMockQueueHandler(controller)

			v.mockFn(userRepo, planRepo, workspaceRepo, preferencesRepo, billingClient)

			cfg := getConfig()
			h := &stripeHandler{
				cfg:             cfg,
				user:            userRepo,
				planRepo:        planRepo,
				workRepo:        workspaceRepo,
				preferencesRepo: preferencesRepo,
				billingClient:   billingClient,
				taskQueue:       queueRepo,
				logger:          getLogger(t),
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(v.payload))

			if v.signature == "valid_signature" {
				timestamp := time.Now().Unix()
				signature := generateStripeSignature(v.payload, cfg.Billing.Stripe.WebhookSecret, timestamp)
				req.Header.Set("Stripe-Signature", signature)
			} else if v.signature != "" {
				req.Header.Set("Stripe-Signature", v.signature)
			}

			h.handleWebhook(rr, req)

			require.Equal(t, v.expectedStatusCode, rr.Code)
			verifyMatch(t, rr)
		})
	}
}

func generateStripeWebhookTestTable() []struct {
	name               string
	mockFn             func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient)
	expectedStatusCode int
	payload            []byte
	signature          string
} {
	// Create test webhook payloads with API version
	customerCreatedPayload := []byte(`{
		"id": "evt_123",
		"type": "customer.created",
		"api_version": "2025-01-27.acacia",
		"created": 1706962800,
		"livemode": false,
		"pending_webhooks": 1,
		"request": {
			"id": null,
			"idempotency_key": null
		},
		"data": {
			"object": {
				"id": "cus_123",
				"object": "customer",
				"created": 1706962800,
				"livemode": false
			}
		}
	}`)

	invoicePaidPayload := []byte(`{
		"id": "evt_123",
		"type": "invoice.paid",
		"api_version": "2025-01-27.acacia",
		"created": 1706962800,
		"livemode": false,
		"pending_webhooks": 1,
		"request": {
			"id": null,
			"idempotency_key": null
		},
		"data": {
			"object": {
				"id": "in_123",
				"object": "invoice",
				"customer": "cus_123",
				"created": 1706962800,
				"livemode": false,
				"lines": {
					"object": "list",
					"data": [{
						"id": "il_123",
						"object": "line_item",
						"plan": {
							"id": "plan_123",
							"object": "plan",
							"product": "prod_123"
						}
					}]
				}
			}
		}
	}`)

	return []struct {
		name               string
		mockFn             func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient)
		expectedStatusCode int
		payload            []byte
		signature          string
	}{
		{
			name: "invalid webhook signature",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
			},
			expectedStatusCode: http.StatusBadRequest,
			payload:            customerCreatedPayload,
			signature:          "invalid_signature",
		},
		{
			name: "missing webhook signature",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
			},
			expectedStatusCode: http.StatusBadRequest,
			payload:            customerCreatedPayload,
			signature:          "",
		},
		{
			name: "customer.created - workspace not found",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
				workspaceRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("workspace not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			payload:            customerCreatedPayload,
			signature:          "valid_signature", // In real test this would be properly generated
		},
		{
			name: "customer.created - add plan fails",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
				workspaceRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Workspace{ID: uuid.New()}, nil)

				billingClient.EXPECT().
					AddPlanToCustomer(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", errors.New("failed to add plan"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			payload:            customerCreatedPayload,
			signature:          "valid_signature",
		},
		{
			name: "customer.created - success",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
				workspaceRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Workspace{ID: uuid.New()}, nil)

				billingClient.EXPECT().
					AddPlanToCustomer(gomock.Any(), gomock.Any()).
					Times(1).
					Return("sub_123", nil)

				workspaceRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			payload:            customerCreatedPayload,
			signature:          "valid_signature",
		},
		{
			name: "invoice.paid - workspace not found",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
				workspaceRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("workspace not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			payload:            invoicePaidPayload,
			signature:          "valid_signature",
		},
		{
			name: "invoice.paid - plan not found",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
				workspaceRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Workspace{ID: uuid.New()}, nil)

				planRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("plan not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			payload:            invoicePaidPayload,
			signature:          "valid_signature",
		},
		{
			name: "invoice.paid - success",
			mockFn: func(userRepo *malak_mocks.MockUserRepository, planRepo *malak_mocks.MockPlanRepository, workspaceRepo *malak_mocks.MockWorkspaceRepository, preferencesRepo *malak_mocks.MockPreferenceRepository, billingClient *malak_mocks.MockClient) {
				workspaceRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Workspace{ID: uuid.New()}, nil)

				planRepo.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&malak.Plan{ID: uuid.New()}, nil)

				workspaceRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			payload:            invoicePaidPayload,
			signature:          "valid_signature",
		},
	}
}

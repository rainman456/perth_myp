package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PaystackVerifyResponse represents the response from verifying a transaction via Paystack.
type PaystackVerifyResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    PaystackTransactionData `json:"data"`
}

// PaystackTransactionData contains detailed information about a verified transaction.
type PaystackTransactionData struct {
	ID              int64                  `json:"id"`
	Domain          string                 `json:"domain"`
	Status          string                 `json:"status"`
	Reference       string                 `json:"reference"`
	Amount          int64                  `json:"amount"`
	Message         string                 `json:"message"`
	GatewayResponse string                 `json:"gateway_response"`
	PaidAt          *time.Time             `json:"paid_at"`
	CreatedAt       *time.Time             `json:"created_at"`
	Channel         string                 `json:"channel"`
	Currency        string                 `json:"currency"`
	IPAddress       string                 `json:"ip_address"`
	Metadata json.RawMessage `json:"metadata"`
	Log             *PaystackTransactionLog `json:"log"`
	Fees            int64                  `json:"fees"`
	Customer        PaystackCustomer        `json:"customer"`
	Authorization   PaystackAuthorization   `json:"authorization"`
	Plan            string                  `json:"plan"`
	Split           interface{}             `json:"split"`
	OrderID         interface{}             `json:"order_id"`
	POSTransactionData interface{}          `json:"pos_transaction_data"`
	Source          interface{}             `json:"source"`
	FeeBreakdown    []interface{}           `json:"fee_breakdown"`
	PaidAtString    string                  `json:"paidAt"`
	CreatedAtString string                  `json:"createdAt"`
}

// PaystackTransactionLog represents details of the transaction logs.
type PaystackTransactionLog struct {
	StartTime int64       `json:"start_time"`
	TimeSpent int64       `json:"time_spent"`
	Attempts  int64       `json:"attempts"`
	Errors    int64       `json:"errors"`
	Success   bool        `json:"success"`
	Mobile    bool        `json:"mobile"`
	Input     []string    `json:"input"`
	History   []LogHistory `json:"history"`
}

// LogHistory represents a single step in the transaction log.
type LogHistory struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	Time        int64  `json:"time"`
	Stage       string `json:"stage"`
	Amount      int64  `json:"amount"`
	IPAddress   string `json:"ip_address"`
	Attempt     int64  `json:"attempt"`
}

// PaystackCustomer represents customer information from Paystack.
type PaystackCustomer struct {
	ID           int64       `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	Metadata     interface{} `json:"metadata"`
	CustomerCode string      `json:"customer_code"`
	RiskAction   string      `json:"risk_action"`
}

// PaystackAuthorization represents card or payment method details.
type PaystackAuthorization struct {
	AuthorizationCode string      `json:"authorization_code"`
	Bin               string      `json:"bin"`
	Last4             string      `json:"last4"`
	ExpMonth          string      `json:"exp_month"`
	ExpYear           string      `json:"exp_year"`
	Channel           string      `json:"channel"`
	CardType          string      `json:"card_type"`
	Bank              string      `json:"bank"`
	CountryCode       string      `json:"country_code"`
	Brand             string      `json:"brand"`
	Reusable          bool        `json:"reusable"`
	Signature         string      `json:"signature"`
	AccountName       interface{} `json:"account_name"`
}


func (s *PaymentService) verifyPaystack(ctx context.Context, reference string) (*PaystackVerifyResponse, error) {
    url := fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+s.config.PaystackSecretKey)
    req.Header.Set("Accept", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
        return nil, fmt.Errorf("paystack verify returned non-200: %d - %s", resp.StatusCode, string(body))
    }

    var pr PaystackVerifyResponse
    dec := json.NewDecoder(resp.Body)
    if err := dec.Decode(&pr); err != nil {
        return nil, fmt.Errorf("failed to decode paystack verify response: %w", err)
    }
    return &pr, nil
}

// normalizeMetadata handles object, string-encoded-json, empty string, and other primitives.
func normalizeMetadata(raw json.RawMessage) (map[string]interface{}, error) {
    if len(raw) == 0 {
        return nil, nil
    }

    // If it's JSON null -> return nil
    if string(raw) == "null" {
        return nil, nil
    }

    // Try decoding as object
    var asObj map[string]interface{}
    if err := json.Unmarshal(raw, &asObj); err == nil {
        // If object is empty, return nil (optional)
        if len(asObj) == 0 {
            return nil, nil
        }
        return asObj, nil
    }

    // Try decoding as a string (e.g. "" or "{\"a\":1}")
    var asStr string
    if err := json.Unmarshal(raw, &asStr); err == nil {
        // empty string -> no metadata
        if asStr == "" {
            return nil, nil
        }
        // maybe it's a JSON-encoded object inside the string
        var inner map[string]interface{}
        if err := json.Unmarshal([]byte(asStr), &inner); err == nil {
            if len(inner) == 0 {
                return nil, nil
            }
            return inner, nil
        }
        // otherwise return the string under a key
        return map[string]interface{}{"value": asStr}, nil
    }

    // Last resort: unmarshal into any and coerce to map if possible
    var any interface{}
    if err := json.Unmarshal(raw, &any); err == nil {
        if mm, ok := any.(map[string]interface{}); ok && len(mm) > 0 {
            return mm, nil
        }
        if s, ok := any.(string); ok && s == "" {
            return nil, nil
        }
        return map[string]interface{}{"value": any}, nil
    }

    return nil, fmt.Errorf("unknown metadata format")
}
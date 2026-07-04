package webhook

import "github.com/google/uuid"

type Customer struct {
	Email    string `json:"email"`
	Document string `json:"document"`
}

type Card struct {
	Bin   string `json:"bin"`
	Last4 string `json:"last4"`
}

type WebhookRequest struct {
	EventID   uuid.UUID `json:"event_id"`
	Type      string    `json:"type"`
	PaymentID uuid.UUID `json:"payment_id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	Currency  string    `json:"currency"`
	Customer  Customer  `json:"customer"`
	Card      Card      `json:"card"`
	IP        string    `json:"ip"`
	DeviceID  uuid.UUID `json:"device_id"`
}

type WebhookResponse struct {
	Status string `json:"status"`
}

package payments

import "time"

type PaymentEventItem struct {
	ID         string    `json:"id"`
	Event_ID   string    `json:"event_id"`
	Payment_ID string    `json:"payment_id"`
	Type       string    `json:"type"`
	Created_At time.Time `json:"created_at"`
}

type PaymentEventListResponse struct {
	Data        []PaymentEventItem `json:"data"`
	Page        int                `json:"page"`
	Limit       int                `json:"limit"`
	Total       int                `json:"total"`
	Total_Pages int                `json:"total_pages"`
}

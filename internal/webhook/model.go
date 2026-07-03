package webhook

type Webhook struct {
	Event_id   string `json:"event_id"`
	Type       string `json:"type"`
	Payment_id string `json:"payment_id"`
	Amount     struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"amount"`
	Status   string `json:"status"`
	Currency string `json:"currency"`
	Customer string `json:"customer"`
	Card     struct {
		Card_id   string `json:"card_id"`
		Last4     string `json:"last4"`
		Exp_month string `json:"exp_month"`
		Exp_year  string `json:"exp_year"`
	} `json:"card"`
	Ip        string `json:"ip"`
	Device_id string `json:"device_id"`
}

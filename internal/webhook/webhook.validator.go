package webhook

import (
	"errors"

	"github.com/google/uuid"
)

func validateWebhook(input WebhookRequest) error {
	switch {
	case input.EventID == uuid.Nil:
		return errors.New("event_id é obrigatório")
	case input.PaymentID == uuid.Nil:
		return errors.New("payment_id é obrigatório")
	case input.Amount <= 0:
		return errors.New("amount deve ser maior que zero")
	case input.Status == "":
		return errors.New("status é obrigatório")
	case input.Currency == "":
		return errors.New("currency é obrigatório")
	case input.Customer.Email == "":
		return errors.New("customer.email é obrigatório")
	case input.Customer.Document == "":
		return errors.New("customer.document é obrigatório")
	case input.Card.Bin == "":
		return errors.New("card.bin é obrigatório")
	case input.Card.Last4 == "":
		return errors.New("card.last4 é obrigatório")
	case input.IP == "":
		return errors.New("ip é obrigatório")
	case input.DeviceID == uuid.Nil:
		return errors.New("device_id é obrigatório")
	default:
		return nil
	}
}

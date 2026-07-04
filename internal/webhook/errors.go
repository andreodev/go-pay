package webhook

import "errors"

var (
	ErrEventAlreadyExists = errors.New("webhook already exists")
	ErrInvalidWebhook     = errors.New("invalid webhook")
)

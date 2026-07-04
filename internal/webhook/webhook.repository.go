package webhook

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ExistsByEventID(ctx context.Context, eventID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM payment_events
			WHERE event_id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, eventID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) CreateEvent(ctx context.Context, webhook *WebhookRequest) error {
	payload, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO payment_events (
			event_id,
			type,
			payment_id,
			payload
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err = r.db.Exec(ctx, query,
		webhook.EventID,
		webhook.Type,
		webhook.PaymentID,
		payload,
	)

	return err
}

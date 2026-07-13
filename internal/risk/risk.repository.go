package risk

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Risk struct {
	PaymentID uuid.UUID `json:"payment_id"`
	EventID   uuid.UUID `json:"event_id"`
	Score     int       `json:"score"`
	Level     string    `json:"level"`
	Reasons   []string  `json:"reasons"`
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRisk(ctx context.Context, payload *Risk) error {
	reasonsJSON, err := json.Marshal(payload.Reasons)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO payment_risks (
			payment_id,
			event_id,
			score,
			level,
			reasons
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`

	_, err = r.db.Exec(
		ctx,
		query,
		payload.PaymentID,
		payload.EventID,
		payload.Score,
		payload.Level,
		reasonsJSON,
	)

	return err
}

func (r *Repository) GetRiskByPaymentID(ctx context.Context, paymentID uuid.UUID) (*Risk, error) {
	query := `
		SELECT payment_id, event_id, score, level, reasons
		FROM payment_risks
		WHERE payment_id = $1
	`

	row := r.db.QueryRow(ctx, query, paymentID)

	var risk Risk
	var reasonsJSON []byte

	err := row.Scan(
		&risk.PaymentID,
		&risk.EventID,
		&risk.Score,
		&risk.Level,
		&reasonsJSON,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(reasonsJSON, &risk.Reasons)
	if err != nil {
		return nil, err
	}

	return &risk, nil
}

package payments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPaymentsEvents(
	ctx context.Context,
	limit int,
	offset int,
) ([]PaymentEventItem, int, error) {
	const query = `
			SELECT
					event_id,
					payment_id,
					type,
					created_at
				FROM payment_events
				ORDER BY created_at DESC
				LIMIT $1
				OFFSET $2
				`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []PaymentEventItem
	for rows.Next() {
		var event PaymentEventItem
		if err := rows.Scan(&event.Event_ID, &event.Payment_ID, &event.Type, &event.Created_At); err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}

	const countQuery = `
			SELECT COUNT(*) FROM payment_events
			`
	var total int
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

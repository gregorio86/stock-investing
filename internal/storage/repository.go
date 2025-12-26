package storage

import (
	"context"
	"time"

	"stock-investing/internal/models"
)

type Repository interface {
	InsertTrade(ctx context.Context, t *models.Trade) error
	ListTrades(ctx context.Context, limit int) ([]*models.Trade, error)
}

type repo struct {
	store *SQLiteStore
}

func NewRepository(store *SQLiteStore) Repository {
	return &repo{store: store}
}

func (r *repo) InsertTrade(ctx context.Context, t *models.Trade) error {
	const q = `
INSERT INTO trades (code, side, quantity, price, time, strategy)
VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.store.DB.ExecContext(
		ctx,
		q,
		t.Code,
		t.Side,
		t.Quantity,
		t.Price,
		t.Time.UTC().Format(time.RFC3339),
		t.Strategy,
	)
	return err
}

func (r *repo) ListTrades(ctx context.Context, limit int) ([]*models.Trade, error) {
	const q = `
SELECT id, code, side, quantity, price, time, strategy
FROM trades
ORDER BY id DESC
LIMIT ?`
	rows, err := r.store.DB.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Trade
	for rows.Next() {
		var t models.Trade
		var ts string
		if err := rows.Scan(
			&t.ID,
			&t.Code,
			&t.Side,
			&t.Quantity,
			&t.Price,
			&ts,
			&t.Strategy,
		); err != nil {
			return nil, err
		}
		parsed, _ := time.Parse(time.RFC3339, ts)
		t.Time = parsed
		out = append(out, &t)
	}
	return out, rows.Err()
}

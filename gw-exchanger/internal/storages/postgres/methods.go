package postgres

import (
	"context"
	"errors"
)

func (s *PostgresStorage) GetAllRates(ctx context.Context) (map[string]float32, error) {
	rows, err := s.pool.Query(ctx, "SELECT currency_from, currency_to, rate FROM rates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float32)
	for rows.Next() {
		var from, to string
		var rate float32
		if err := rows.Scan(&from, &to, &rate); err != nil {
			return nil, err
		}
		key := from + "_" + to
		result[key] = rate
	}

	return result, nil
}

func (s *PostgresStorage) GetRate(ctx context.Context, from, to string) (float32, error) {
	var rate float32
	err := s.pool.QueryRow(ctx,
		"SELECT rate FROM rates WHERE currency_from=$1 AND currency_to=$2 LIMIT 1",
		from, to,
	).Scan(&rate)
	if err != nil {
		return 0, errors.New("rate not found")
	}
	return rate, nil
}

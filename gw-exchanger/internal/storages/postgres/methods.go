package postgres

import (
	"context"
)

func (s *PostgresStorage) GetAllRates(ctx context.Context) (map[string]float32, error) {
	rows, err := s.pool.Query(ctx, "SELECT currency, rate FROM rates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float32)
	for rows.Next() {
		var currency string
		var rate float32
		if err := rows.Scan(&currency, &rate); err != nil {
			return nil, err
		}
		result[currency] = rate
	}

	return result, nil
}

func (s *PostgresStorage) GetRate(ctx context.Context, currency string) (float32, error) {
	var rate float32
	err := s.pool.QueryRow(ctx, "SELECT rate FROM rates WHERE currency=$1", currency).Scan(&rate)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

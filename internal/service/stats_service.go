package service

import (
	"context"
	"fmt"

	"github.com/YannisMuminov/internal/domain"
	"github.com/jmoiron/sqlx"
)

type statsService struct {
	db *sqlx.DB
}

func NewStatsService(db *sqlx.DB) domain.StatsService {
	return &statsService{db: db}
}

func (s *statsService) GetStats(ctx context.Context) (*domain.Stats, error) {
	var stats domain.Stats
	var err error

	stats.Deals, err = s.getDealStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get deal stats: %w", err)
	}

	stats.Clients, err = s.getClientStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get client stats: %w", err)
	}

	stats.Task, err = s.getTaskStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get task stats: %w", err)
	}

	stats.Users, err = s.getUserStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user stats: %w", err)
	}

	return &stats, nil
}

func (s *statsService) getDealStats(ctx context.Context) (domain.DealStats, error) {
	query := `
		 SELECT
    COUNT(*) AS total,
    COUNT(CASE WHEN status = 'won' THEN  1 END) AS won,
    COUNT(CASE WHEN  status = 'lost' THEN 1 END) AS lost,
    COUNT(CASE WHEN  status IN ('new', 'negotiation') THEN 1 END) AS is_progress,
    COALESCE(SUM(amount), 0) AS total_amount,
    COALESCE(SUM(CASE WHEN status = 'won' THEN amount ELSE 0 END), 0) AS won_amount
		FROM deals
`

	var stats domain.DealStats

	err := s.db.QueryRowContext(ctx, query).Scan(
		&stats.Total,
		&stats.Won,
		&stats.Lost,
		&stats.InProgress,
		&stats.TotalAmount,
		&stats.WonAmount,
	)

	if err != nil {
		return domain.DealStats{}, fmt.Errorf("get deal stats: %w", err)
	}
	return stats, nil
}

func (s *statsService) getClientStats(ctx context.Context) (domain.ClientStats, error) {
	query := `
		SELECT
    COUNT(*) AS total,
    COUNT(CASE WHEN created_at >= DATE_TRUNC('month', NOW()) THEN 1 END) AS new_this_month
	FROM clients
	`

	var stats domain.ClientStats

	err := s.db.QueryRowContext(ctx, query).Scan(
		&stats.Total,
		&stats.NewThisMonth,
	)

	if err != nil {
		return domain.ClientStats{}, fmt.Errorf("get client stats: %w", err)
	}
	return stats, nil
}

func (s *statsService) getTaskStats(ctx context.Context) (domain.TaskStats, error) {
	query := `
		SELECT
    COUNT(*) AS total,
    COUNT(CASE WHEN status = 'done' THEN 1 END) AS done,
    COUNT(CASE WHEN due_data < NOW() AND status != 'done' AND due_data IS NOT NULL THEN 1 END) AS overdue
		FROM tasks
		`

	var stats domain.TaskStats
	err := s.db.QueryRowContext(ctx, query).Scan(
		&stats.Total,
		&stats.Done,
		&stats.Overdue,
	)

	if err != nil {
		return domain.TaskStats{}, fmt.Errorf("get tasks stats: %w", err)
	}
	return stats, nil
}

func (s *statsService) getUserStats(ctx context.Context) (domain.UserStats, error) {
	query := `
		SELECT 
		COUNT(*) AS total,
		COUNT(CASE WHEN is_active = true THEN 1 END) AS is_active
		FROM users
`
	var stats domain.UserStats
	err := s.db.QueryRowContext(ctx, query).Scan(

		&stats.Total,
		&stats.Active,
	)
	if err != nil {
		return domain.UserStats{}, fmt.Errorf("get users stats: %w", err)
	}

	return stats, nil
}

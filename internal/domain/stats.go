package domain

import "context"

type DealStats struct {
	Total       int `json:"total"`
	Won         int `json:"won"`
	Lost        int `json:"lost"`
	InProgress  int `json:"in_progress"`
	TotalAmount int `json:"total_amount"`
	WonAmount   int `json:"won_amount"`
}

type ClientStats struct {
	Total        int `json:"total"`
	NewThisMonth int `json:"new_this_month"`
}

type TaskStats struct {
	Total   int `json:"total"`
	Done    int `json:"done"`
	Overdue int `json:"over_due"`
}

type UserStats struct {
	Total  int `json:"total"`
	Active int `json:"active"`
}

type Stats struct {
	Deals   DealStats   `json:"deals"`
	Clients ClientStats `json:"clients"`
	Task    TaskStats   `json:"task"`
	Users   UserStats   `json:"users"`
}

type StatsService interface {
	GetStats(ctx context.Context) (*Stats, error)
}

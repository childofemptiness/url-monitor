package monitor

import (
	"url-monitor/internal/storage/postgres"
)

type Service struct {
	repo *postgres.Repository
}

func NewMonitorService(repo *postgres.Repository) *Service {
	return &Service{repo: repo}
}

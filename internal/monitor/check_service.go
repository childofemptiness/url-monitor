package monitor

import (
	"context"
	"time"
)

type CheckRepository interface {
	CompleteCheck(ctx context.Context, check MonitorCheck, nextCheckAt time.Time) error
}

type CheckStoreService struct {
	repo CheckRepository
}

func NewCheckStoreService(repo CheckRepository) *CheckStoreService {
	return &CheckStoreService{repo: repo}
}

func (c *CheckStoreService) SaveCheckResult(ctx context.Context, check MonitorCheck, nextCheckAt time.Time) error {
	return c.repo.CompleteCheck(ctx, check, nextCheckAt)
}

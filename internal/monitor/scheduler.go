package monitor

import (
	"context"
	"time"
)

type SchedulerRepository interface {
	ListDue(ctx context.Context, now time.Time, limit int) ([]Monitor, error)
}

type Scheduler struct {
	repo         SchedulerRepository
	monitorsChan chan Monitor
}

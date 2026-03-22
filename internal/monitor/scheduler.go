package monitor

import (
	"context"
	"time"
)

type SchedulerRepository interface {
	ListDue(ctx context.Context, now time.Time, liimt int) ([]Monitor, error)
}

type Scheduler struct {
	repo 		 SchedulerRepository
	monitorsChan chan Monitor
}

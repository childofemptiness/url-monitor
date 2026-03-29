package monitor

import (
	"context"
	"log"
	"time"
)

const (
	checksLimit = 5
)

type SchedulerRepository interface {
	ListDue(ctx context.Context, now time.Time, limit int) ([]Monitor, error)
}

type Publisher interface {
	Submit(ctx context.Context, m Monitor) error
}
type Scheduler struct {
	repo         SchedulerRepository
	publisher    Publisher
	timeInterval time.Duration
}

func NewScheduler(
	schedulerRepo SchedulerRepository,
	publisher Publisher,
	timeInterval time.Duration,
) *Scheduler {
	return &Scheduler{
		repo:         schedulerRepo,
		publisher:    publisher,
		timeInterval: timeInterval,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	log.Printf("scheduler started")

	ticker := time.NewTicker(s.timeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err := s.runOnce(ctx)
			if err == nil {
				continue
			}

			log.Printf("scheduler temporary error: %v", err)
			return err
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) error {
	monitors, err := s.repo.ListDue(ctx, time.Now(), checksLimit)
	if err != nil {
		return err
	}

	for _, monitor := range monitors {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := s.publisher.Submit(ctx, monitor)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

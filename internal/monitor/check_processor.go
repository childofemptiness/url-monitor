package monitor

import (
	"context"
	"log"
	"time"
)

type Checker interface {
	Check(ctx context.Context, monitor Monitor) MonitorCheck
}
type CheckService interface {
	SaveCheckResult(ctx context.Context, check MonitorCheck, nextCheckAt time.Time) error
}
type CheckProcessor struct {
	checker      Checker
	checkService CheckService
}

func NewCheckProcessor(checker Checker, checkService CheckService) *CheckProcessor {
	return &CheckProcessor{
		checker:      checker,
		checkService: checkService,
	}
}

func (cp *CheckProcessor) Process(ctx context.Context, monitor Monitor) error {
	log.Printf("process check for monitor m: %+v\n", monitor)
	check := cp.checker.Check(ctx, monitor)
	nextCheckAt := check.FinishedAt.Add(time.Duration(monitor.IntervalSeconds) * time.Second)

	err := cp.checkService.SaveCheckResult(ctx, check, nextCheckAt)
	if err != nil {
		return err
	}

	return nil
}

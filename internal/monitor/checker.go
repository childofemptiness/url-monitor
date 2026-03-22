package monitor

import "context"

type Checker interface {
	Check(ctx context.Context, m Monitor) MonitorCheck
}

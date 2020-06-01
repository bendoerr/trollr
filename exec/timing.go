package exec

import (
	"context"

	servertiming "github.com/mitchellh/go-server-timing"
)

type TimingExecutor struct {
	delegate Executor
}

func NewTimingExecutor(delegate Executor) *TimingExecutor {
	return &TimingExecutor{
		delegate: delegate,
	}
}

func (tx *TimingExecutor) Run(ctx context.Context, command string, stdin *string, args ...string) Execution {
	t := servertiming.FromContext(ctx)
	m := t.NewMetric("external-process").Start()
	e := tx.delegate(ctx, command, stdin, args...)
	m.Stop()
	return e
}

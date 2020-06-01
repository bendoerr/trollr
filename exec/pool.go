package exec

import (
	"context"
	"time"

	"github.com/panjf2000/ants/v2"
)

type PoolExecutor struct {
	pool     *ants.Pool
	delegate Executor
	timeout  time.Duration
}

func NewPoolExecutor(delegate Executor) *PoolExecutor {
	pool, _ := ants.NewPool(10, ants.WithMaxBlockingTasks(100))

	return &PoolExecutor{
		pool:     pool,
		timeout:  30 * time.Second,
		delegate: delegate,
	}
}

func (px *PoolExecutor) Run(ctx context.Context, command string, stdin *string, args ...string) Execution {
	rchan := make(chan Execution)

	ctx, cancel := context.WithTimeout(ctx, px.timeout)
	defer cancel()

	err := px.pool.Submit(func() {
		rchan <- px.delegate(ctx, command, stdin, args...)
	})

	if err != nil {
		return &CommandExec{
			err: err,
		}
	}

	select {
	case r := <-rchan:
		return r
	case <-ctx.Done():
		return &CommandExec{
			err: ctx.Err(),
		}
	}
}

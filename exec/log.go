package exec

import (
	"context"

	"go.uber.org/zap"
)

type LoggingExecutor struct {
	delegate Executor
	logger   *zap.Logger
}

func NewLoggingExecutor(delegate Executor, logger *zap.Logger) *LoggingExecutor {
	return &LoggingExecutor{
		delegate: delegate,
		logger:   logger.Named("executor"),
	}
}

func (lx *LoggingExecutor) Run(ctx context.Context, command string, stdin *string, args ...string) Execution {
	lx.logger.Info("executing",
		zap.String("command", command),
		zap.String("stdin", *stdin),
		zap.Strings("args", args),
	)

	return lx.delegate(ctx, command, stdin, args...)
}

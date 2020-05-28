package util

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

type Executor = func(ctx context.Context, command string, stdin *string, args ...string) Execution

type ExecutionExitError = exec.ExitError

type Execution interface {
	Command() string
	Args() []string
	Stdout() *bytes.Buffer
	Stderr() *bytes.Buffer
	ExitCode() int
	Err() error
}

type CommandExec struct {
	command  string
	args     []string
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
	exitCode int
	err      error
}

func (c *CommandExec) Command() string {
	return c.command
}

func (c *CommandExec) Args() []string {
	return c.args
}

func (c *CommandExec) Stdout() *bytes.Buffer {
	return c.stdout
}

func (c *CommandExec) Stderr() *bytes.Buffer {
	return c.stderr
}

func (c *CommandExec) ExitCode() int {
	return c.exitCode
}

func (c *CommandExec) Err() error {
	return c.err
}

func Run(ctx context.Context, command string, stdin *string, args ...string) Execution {
	var result CommandExec
	var cmd *exec.Cmd
	var cmdOut, cmdErr bytes.Buffer

	// Create the resulting execution
	result = CommandExec{
		command: command,
		args:    args,
	}

	// Create a new Command
	cmd = exec.CommandContext(
		ctx,
		command,
		args...,
	)

	// Attach stdin
	if stdin != nil {
		cmd.Stdin = bytes.NewBufferString(*stdin)
	}

	// Attach stdout
	cmd.Stdout = &cmdOut
	result.stdout = &cmdOut

	// Attach stderr
	cmd.Stderr = &cmdErr
	result.stderr = &cmdErr

	// Execute the command
	result.err = cmd.Run()

	// Error Checking
	if result.err != nil {
		var exitError *ExecutionExitError
		if errors.As(result.err, &exitError) {
			result.exitCode = exitError.ExitCode()
		} else {
			result.exitCode = 1
		}
	}

	return &result
}

type TestExecution struct {
	TestCommand  string
	TestArgs     []string
	TestStdin    string
	TestStdout   *bytes.Buffer
	TestStderr   *bytes.Buffer
	TestExitCode int
	TestErr      error
}

func (e *TestExecution) Command() string {
	return e.TestCommand
}

func (e *TestExecution) Args() []string {
	return e.TestArgs
}

func (e *TestExecution) Stdout() *bytes.Buffer {
	return e.TestStdout
}

func (e *TestExecution) Stderr() *bytes.Buffer {
	return e.TestStderr
}

func (e *TestExecution) ExitCode() int {
	return e.TestExitCode
}

func (e *TestExecution) Err() error {
	return e.TestErr
}

func NewTestExecutor(stdout string, stderr string, exitCode int, err error) Executor {
	return func(ctx context.Context, command string, stdin *string, args ...string) Execution {
		readIn := ""

		if stdin != nil {
			readIn = *stdin
		}

		return &TestExecution{
			TestCommand:  command,
			TestArgs:     args,
			TestStdin:    readIn,
			TestStdout:   bytes.NewBufferString(stdout),
			TestStderr:   bytes.NewBufferString(stderr),
			TestExitCode: exitCode,
			TestErr:      err,
		}
	}
}

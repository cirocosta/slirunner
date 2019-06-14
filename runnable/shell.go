package runnable

import (
	"context"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

// ShellCommand executes a given command using a shell.
//
type ShellCommand struct {
	command string
	stderr  io.Writer
}

var _ Runnable = &ShellCommand{}

// NewShellCommand instantiates a new ShellCommand that
// is meant to always run with a specified `command`.
//
func NewShellCommand(command string, stderr io.Writer) (runnable *ShellCommand) {
	runnable = &ShellCommand{
		command: command,
	}

	return
}

// Run runs the command either until completion or context
// cancellation.
//
func (r *ShellCommand) Run(ctx context.Context) (err error) {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", r.command)

	_, err = cmd.CombinedOutput()

	switch ctx.Err() {
	case context.DeadlineExceeded:
		err = errors.Wrapf(ctx.Err(),
			"command didn't finish on time")
	case context.Canceled:
		err = errors.Wrapf(ctx.Err(),
			"command execution cancelled")
	}

	if err != nil {
		err = errors.Wrapf(err,
			"command execution failed")
		return
	}

	return
}

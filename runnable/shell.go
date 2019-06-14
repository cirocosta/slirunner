package runnable

import (
	"context"
	"io"
	"os/exec"
)

type ShellRunnable struct {
	command string
}

func NewShellRunnable(command string, stdout, stderr io.Writer) (runnable *ShellRunnable) {
	runnable = &ShellRunnable{
		command: command,
	}

	return
}

func (r *ShellRunnable) Run(ctx context.Context) (succeeded bool) {
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", r.command)

	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return
	}

	if err != nil {
		succeeded = false
		// send ALL output to stderr
	}

	return
}

package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cirocosta/slirunner/probes"
)

type runCommand struct{}

func (c *runCommand) Execute(args []string) (err error) {
	err = probes.All.Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return
}

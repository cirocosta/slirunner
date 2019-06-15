package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cirocosta/slirunner/probes"
)

type runCommand struct {
	Target          string `long:"target" required:"true"`
	PipelinesPrefix string `long:"prefix" default:"slirunner-"`
}

func (c *runCommand) Execute(args []string) (err error) {
	err = probes.New(c.Target, c.PipelinesPrefix).Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return
}

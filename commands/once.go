package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/cirocosta/slirunner/probes"
)

type onceCommand struct {
	Target          string `long:"target" required:"true"`
	PipelinesPrefix string `long:"prefix" default:"slirunner-"`
}

func (c *onceCommand) Execute(args []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go onTerminationSignal(cancel)

	err = probes.New(c.Target, c.PipelinesPrefix).Run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return
}

package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cirocosta/slirunner/exporter"
	"github.com/cirocosta/slirunner/probes"
)

type startCommand struct {
	Target          string        `long:"target" required:"true"`
	PipelinesPrefix string        `long:"prefix" default:"slirunner-"`
	Interval        time.Duration `long:"interval" default:"1m"`

	Prometheus exporter.Exporter `group:"Prometheus configuration"`
}

func (c *startCommand) Execute(args []string) (err error) {
	var (
		allProbes = probes.New(c.Target, c.PipelinesPrefix)
		ticker    = time.NewTicker(c.Interval)
	)

	ctx, cancel := context.WithCancel(context.Background())
	go onTerminationSignal(func() {
		cancel()
		c.Prometheus.Close()
	})

	go func() {
		err := c.Prometheus.Listen()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}

		cancel()
	}()

	f := func() {
		err = allProbes.Run(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	f()

	for {
		select {
		case <-ticker.C:
			f()
		case <-ctx.Done():
			c.Prometheus.Close()
		}
	}

	return
}

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

	allProbes.Run(ctx)

	for {
		select {
		case <-ticker.C:
			allProbes.Run(ctx)
		case <-ctx.Done():
			c.Prometheus.Close()
			return
		}
	}

	return
}

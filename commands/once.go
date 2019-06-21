package commands

import (
	"context"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/cirocosta/slirunner/probes"
)

type onceCommand struct {
	Target          string `long:"target" required:"true"`
	PipelinesPrefix string `long:"prefix" default:"slirunner-"`

	Username     string `long:"username"      short:"u" required:"true"`
	Password     string `long:"password"      short:"p" required:"true"`
	ConcourseUrl string `long:"concourse-url" short:"c" required:"true"`
}

func (c *onceCommand) Execute(args []string) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	go onTerminationSignal(cancel)

	logger := lager.NewLogger("my-app")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))

	ctx = lagerctx.NewContext(ctx, logger)

	err = probes.NewAll(
		c.Target,
		c.Username, c.Password,
		c.ConcourseUrl,
		c.PipelinesPrefix,
	).Run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return
}

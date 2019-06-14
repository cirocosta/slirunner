package main

import (
	"os"

	"github.com/cirocosta/slirunner/commands"
	"github.com/jessevdk/go-flags"
)

func main() {
	_, err := flags.Parse(&commands.SLIRunner)
	if err != nil {
		os.Exit(1)
	}

	return
}

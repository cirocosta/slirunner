package main

import (
	"fmt"
	"os"

	"github.com/cirocosta/slirunner/commands"
	"github.com/jessevdk/go-flags"
	"github.com/vito/twentythousandtonnesofcrudeoil"
)

func main() {
	parser := flags.NewParser(&commands.SLIRunner, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	twentythousandtonnesofcrudeoil.TheEnvironmentIsPerfectlySafe(parser, "SR_")

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return
}

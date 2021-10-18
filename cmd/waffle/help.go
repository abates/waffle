package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/abates/waffle"
)

var helpCmd *waffle.Command

func init() {
	helpCmd = app.AddCommand("help", "obtain more information about a command", helpCmdFunc)
}

func helpCmdFunc(args ...string) error {
	if len(args) < 1 {
		fmt.Fprintf(helpCmd.Output(), "Usage: %s <command>\n", strings.Join(helpCmd.Path(), " "))
		os.Exit(1)
	} else if cmd, found := app.Lookup(args[0]); found {
		cmd.Usage()
	} else {
		fmt.Fprintf(helpCmd.Output(), "Command %q not found\n", args[0])
		fmt.Fprintf(helpCmd.Output(), "Usage: %s <command>\n", strings.Join(helpCmd.Path(), " "))
		os.Exit(2)
	}
	return nil
}

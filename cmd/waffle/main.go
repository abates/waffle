package main

import (
	"os"
	"path/filepath"

	"github.com/abates/waffle"
)

/*type command struct {
	desc  string
	run   func([]string) error
	flags *flag.FlagSet
}

var commands = map[string]command{}

func addCommand(name, desc string, run func([]string) error) command {
	cmd := command{
		desc: desc,
		run:  run,
	}

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.SetOutput(output)
	flags.Usage = func() { helpUsage(name, cmd) }

	cmd.flags = flags
	commands[name] = cmd
	return cmd
}

var output = os.Stderr

func cmdUsage(fmtstr string, a ...interface{}) {
	fmt.Fprintf(output, "Usage: %s %s\n", filepath.Base(os.Args[0]), fmt.Sprintf(fmtstr, a...))
}

func usage() {
	l := 0
	names := []string{}
	for name := range commands {
		names = append(names, name)
		if len(name) > l {
			l = len(name)
		}
	}

	sort.Strings(names)
	cmdUsage("<command> [arguments]")
	fmt.Fprintf(output, "\nAvailable Commands:\n")
	format := fmt.Sprintf("     %%%ds %%s\n", l)
	for _, name := range names {
		cmd := commands[name]
		fmt.Fprintf(output, format, name, cmd.desc)
	}
	fmt.Fprintf(output, "\nUse \"%s help <command>\" for more information about a command\n", filepath.Base(os.Args[0]))
}*/

var app = waffle.NewCommand()

func main() {
	app.Name = filepath.Base(os.Args[0])
	app.Run(os.Args)
	/*if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	} else if cmd, found := commands[os.Args[1]]; found {
		cmd.flags.Parse(os.Args[2:])
		err := cmd.run(cmd.flags.Args())
		if err == nil {
			os.Exit(0)
		}
		fmt.Fprintf(output, "Command %q failed with %v\n", os.Args[1], err)
		os.Exit(2)
	}

	fmt.Fprintf(output, "Command %q not found\n", os.Args[1])
	usage()
	os.Exit(3)*/
}

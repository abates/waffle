package waffle

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type CommandFunc func([]string) error

type Command struct {
	Name  string
	Desc  string
	Run   CommandFunc
	Flags *flag.FlagSet
	Usage func()

	parent   *Command
	commands map[string]*Command
	output   io.Writer
}

func NewCommand() *Command {
	cmd := &Command{
		commands: make(map[string]*Command),
		output:   os.Stderr,
	}
	cmd.Usage = func() {
		fmt.Fprintf(cmd.output, "Usage: %s", cmd.Name)
		cmd.PrintHelp()
	}
	cmd.Run = cmd.Runner
	return cmd
}

func (cmd *Command) setOutput(output io.Writer) {
	cmd.output = output
	for _, subcmd := range cmd.commands {
		subcmd.setOutput(output)
	}
}

func (cmd *Command) SetOutput(output io.Writer) {
	if cmd.parent == nil {
		cmd.setOutput(output)
	} else {
		cmd.parent.SetOutput(output)
	}
}

func (cmd *Command) PrintHelp() {
	path := []string{}
	if cmd.parent != nil {
		for parent := cmd.parent; parent != nil; parent = parent.parent {
			path = append([]string{parent.Name}, path...)
		}
	}

	prefix := ""
	suffix := ""

	if cmd.Flags != nil {
		prefix = " [arguments]"
	}

	if len(cmd.commands) > 0 {
		suffix = " <command>"
	}

	fmt.Fprintf(cmd.output, "%s%s%s\n", strings.Join(path, " "), prefix, suffix)

	cmd.printFlags()
	cmd.printCommands()
}

func (cmd *Command) hasFlags() bool {

	hasFlags := false
	if cmd.Flags != nil {
		cmd.Flags.VisitAll(func(*flag.Flag) {
			hasFlags = true
		})
	}
	return hasFlags
}

func (cmd *Command) printFlags() {
	if cmd.hasFlags() {
		fmt.Fprintf(cmd.output, "Flags:\n")
		cmd.Flags.PrintDefaults()
	}
}

func (cmd *Command) printCommands() {
	if len(cmd.commands) == 0 {
		return
	}

	l := 0
	names := []string{}
	for name := range cmd.commands {
		names = append(names, name)
		if len(name) > l {
			l = len(name)
		}
	}

	sort.Strings(names)
	fmt.Fprintf(cmd.output, "\nAvailable Commands:\n")
	format := fmt.Sprintf("     %%%ds %%s\n", l)
	for _, name := range names {
		cmd := cmd.commands[name]
		fmt.Fprintf(cmd.output, format, name, cmd.Desc)
	}
}

func (cmd *Command) AddCommand(name, desc string, run func([]string) error) *Command {
	subcmd := &Command{
		Name: name,
		Desc: desc,
		Run:  run,

		commands: make(map[string]*Command),
		output:   cmd.output,
	}

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.SetOutput(cmd.output)
	flags.Usage = cmd.Usage
	subcmd.Flags = flags

	cmd.commands[name] = subcmd
	return subcmd
}

func (cmd *Command) Runner(args []string) error {
	if len(args) < 2 {
		cmd.Usage()
		os.Exit(1)
	} else if cmd, found := cmd.commands[args[1]]; found {
		cmd.Flags.Parse(args[2:])
		err := cmd.Run(cmd.Flags.Args())
		if err == nil {
			os.Exit(0)
		}
		fmt.Fprintf(cmd.output, "Command %q failed with %v\n", args[1], err)
		os.Exit(2)
	}

	fmt.Fprintf(cmd.output, "Command %q not found\n", args[1])
	cmd.Usage()
	os.Exit(3)
	return nil
}

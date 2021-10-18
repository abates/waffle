package waffle

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/abates/formatter"
)

type commandError struct {
	error
	cmd *Command
}

func (ce commandError) Unwrap() error { return ce.error }

var ErrUsage = errors.New("Invalid command usage")

var Logger = formatter.ColorLogger()

type CommandFunc func(...string) error

type Command struct {
	Name     string
	Desc     string
	Run      CommandFunc
	Flags    *flag.FlagSet
	Usage    func()
	UsageStr string

	parent   *Command
	commands map[string]*Command

	output io.Writer
}

func Usage(cmd *Command) func() {
	return func() {
		cmd.PrintUsage()
		cmd.PrintHelp()
	}

}

func NewCommand() *Command {
	cmd := &Command{
		commands: make(map[string]*Command),
		output:   os.Stderr,
	}

	cmd.Usage = Usage(cmd)
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

func (cmd *Command) Path() []string {
	path := []string{cmd.Name}
	if cmd.parent != nil {
		for parent := cmd.parent; parent != nil; parent = parent.parent {
			path = append([]string{parent.Name}, path...)
		}
	}
	return path
}

func (cmd *Command) PrintHelp() {
	cmd.printFlags()
	cmd.printCommands()
}

func (cmd *Command) PrintUsage() {
	prefix := ""
	suffix := ""

	if cmd.hasFlags() {
		prefix = " [flags]"
	}

	if len(cmd.commands) > 0 {
		suffix = " <command>"
	} else if cmd.UsageStr != "" {
		suffix = " " + cmd.UsageStr
	}

	fmt.Fprintf(cmd.output, "Usage: %s%s%s\n", strings.Join(cmd.Path(), " "), prefix, suffix)
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

func (cmd *Command) AddCommand(name, desc string, run CommandFunc) *Command {
	subcmd := &Command{
		Name: name,
		Desc: desc,
		Run:  run,

		commands: make(map[string]*Command),
		output:   cmd.output,
		parent:   cmd,
	}

	if subcmd.Run == nil {
		subcmd.Run = subcmd.Runner
	}

	subcmd.Usage = Usage(subcmd)

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.SetOutput(cmd.output)
	flags.Usage = subcmd.Usage
	subcmd.Flags = flags

	cmd.commands[name] = subcmd
	return subcmd
}

func (cmd *Command) Output() io.Writer {
	return cmd.output
}

func (cmd *Command) Lookup(name string) (subcmd *Command, found bool) {
	cmd, found = cmd.commands[name]
	return cmd, found
}

func (cmd *Command) Runner(args ...string) (err error) {
	if len(args) < 1 {
		err = commandError{fmt.Errorf("%w: expecting sub-command", ErrUsage), cmd}
	} else if subcmd, found := cmd.Lookup(args[0]); found {
		subcmd.Flags.Parse(args[1:])
		err = subcmd.Run(subcmd.Flags.Args()...)
		if err != nil {
			// don't re-wrap the error
			if _, ok := err.(commandError); !ok {
				err = commandError{err, subcmd}
			}
		}
	} else {
		err = commandError{fmt.Errorf("%w: Unknown command %q", ErrUsage, args[0]), cmd}
	}

	if cmd.parent == nil && err != nil {
		if ce, ok := err.(commandError); ok {
			if errors.Is(ce.error, ErrUsage) {
				ce := err.(commandError)
				Logger.Logf("%v", ce.error)
				ce.cmd.Usage()
			} else {
				Logger.Logf("Command %s <fail>failed</fail>: %v", ce.cmd.Name, ce.error)
			}
		} else {
			Logger.Logf("Command %s <fail>failed</fail>: %v", cmd.Name, err)
		}
	}

	return err
}

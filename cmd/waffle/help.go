package main

func init() {
	app.AddCommand("help", "obtain more information about a command", helpCmd)
}

/*func helpUsage(name string, cmd command) {
	if cmd.flags == nil {
		cmdUsage("%s\n", name)
	} else {
		cmdUsage("%s [arguments]\n", name)
		fmt.Fprintf(output, "Arguments:\n")
		cmd.flags.PrintDefaults()
	}
}*/

func helpCmd(args []string) error {
	/*if len(args) < 1 {
		cmdUsage("help <command>")
		os.Exit(1)
	} else if cmd, found := commands[args[0]]; found {
		helpUsage(args[0], cmd)
	} else {
		fmt.Fprintf(output, "Command %q not found\n", args[0])
		usage()
		os.Exit(2)
	}*/
	return nil
}

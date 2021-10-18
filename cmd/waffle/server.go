package main

import "github.com/abates/waffle"

func init() {
	serverCmd := app.AddCommand("server", "manage api server controllers and endpoints", nil)

	addCmd := serverCmd.AddCommand("add", "add controllers, endpoints and security", nil)
	ctrlCmd := addCmd.AddCommand("controller", "add a controller to the server", addController)
	ctrlCmd.UsageStr = "<name> <path>"
	addCmd.AddCommand("endpoint", "add a controller to the server", addEndpoint)

	removeCmd := serverCmd.AddCommand("remove", "remove controllers, endpoints and security", nil)
	removeCmd.AddCommand("controller", "remove a controller from the server", rmController)
	removeCmd.AddCommand("endpoint", "remove a controller from the server", rmEndpoint)

	serverCmd.AddCommand("generate", "generate code specified in openapi.json", generateServer)
}

func addController(args ...string) error {
	if len(args) < 2 {
		return waffle.ErrUsage
	}

	config().AddController(args[0])
	err := config().SaveDef()
	if err == nil {
		err = genCmd()
	}
	return err
}

func addEndpoint(args ...string) error {
	return nil
}

func rmController(args ...string) error {
	return nil
}

func rmEndpoint(args ...string) error {
	return nil
}

func generateServer(args ...string) error {
	return nil
}

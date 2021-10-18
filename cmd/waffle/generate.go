package main

import (
	"fmt"

	"github.com/abates/waffle"
)

func init() {
	app.AddCommand("generate", "(re)generate all code for the project", genCmd)
}

func genCmd(args ...string) (err error) {
	if len(args) > 0 {
		return fmt.Errorf("unexpected argument %q", args)
	}

	return waffle.ExecuteTemplates("generate", ".", *config())
}

package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/abates/waffle"
)

var app = waffle.NewCommand()

var log = waffle.Logger

var c *waffle.Config

func exit(format string, v ...interface{}) {
	log.Logf(format, v...)
	os.Exit(1)
}

func config() *waffle.Config {
	if c == nil {
		c = &waffle.Config{}
		err := c.LoadDef()
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			exit("Failed to load <fail>default config</fail>: %v", err.Error())
		}
	}
	return c
}

func main() {
	app.Name = filepath.Base(os.Args[0])
	err := app.Run(os.Args[1:]...)
	if err != nil {
		os.Exit(1)
	}
}

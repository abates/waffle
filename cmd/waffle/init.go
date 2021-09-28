package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/abates/waffle"
)

var config waffle.Config

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to determine local directory name: %v", err)
		dir = "main"
	}

	cmd := app.AddCommand("init", "initialize current directory with new project tree", initCmd)
	cmd.Flags.StringVar(&config.Pkg, "pkg", filepath.Base(dir), "Package name to use for generated code")
}

func initCmd(args []string) error {
	return waffle.ExecuteTemplates(".", config)
}

package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/abates/waffle"
)

var initRepo *waffle.GitRepo
var gitRemote string

func init() {
	config().LoadDef()

	dir, err := os.Getwd()
	if err != nil {
		exit("Failed to determine <fail>local directory name</fail>: %v", err.Error())
	}

	initRepo, err = waffle.OpenGit(".")
	if err == nil {
		// load version from git
		if config().Module.Version == (waffle.Version{0, 0, 0}) {
			config().Module.Version, err = initRepo.CurrentVersion()
			if errors.Is(err, waffle.ErrNoGitVersion) {
				err = nil
			}
		}
	} else if errors.Is(err, waffle.ErrNoGitRepo) {
		err = nil
	}

	if err == nil {
		config().Maintainer, err = waffle.LoadGitMaintainer()
	}

	if err != nil {
		exit("Failed to initialize command: %v", err)
	}

	cmd := app.AddCommand("init", "initialize current directory with new project tree", initCmd)
	cmd.Flags.StringVar(&config().Name, "name", filepath.Base(dir), "Project name")
	cmd.Flags.StringVar(&config().Desc, "desc", "", "Project description")
	cmd.Flags.Var(&config().Module.Version, "version", "Current version")
	cmd.Flags.StringVar(&config().Maintainer.Name, "maintainer", config().Maintainer.Name, "Maintainer name")
	cmd.Flags.StringVar(&config().Maintainer.Email, "email", config().Maintainer.Email, "Maintainer email")
	cmd.Flags.StringVar(&config().URL, "url", "", "Project Webpage URL")
	cmd.Flags.StringVar(&config().Module.Path, "mod", config().Module.Path, "Go Module Path")
	cmd.Flags.StringVar(&gitRemote, "origin", "", "Go Module Path")
}

func initCmd(args ...string) (err error) {
	if initRepo == nil {
		initRepo, err = waffle.InitGit(".")
		if err == nil {
			if config().Module.Path == "" {
				config().Module.Path = waffle.PromptStr("Module Path: ")
			}

			if gitRemote == "" {
				gitRemote = waffle.PromptStr("Git Remote URL: ")
			}

			if initRepo.SetOrigin(gitRemote) != nil {
				exit("Couldn't set git remote: %v", err)
			}
		} else {
			exit("Could not initialize <fail>empty git repo</fail>: %v", err)
		}
	}

	err = config().SaveDef()
	if err == nil {
		err = genCmd()
	}
	return err
}

// Copyright 2021 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type autogenCommand struct {
	inputFilename  string
	outputFilename string
	data           func() interface{}
}

var autogenCommands = make(map[string]autogenCommand)

func main() {
	var maintainer, pkg string

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <command>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Available Flags:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&maintainer, "m", "Snakeoil Ltd", "the name/info about the maintainer")
	flag.StringVar(&pkg, "p", "main", "the package to be used in the source file")
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	command, found := autogenCommands[flag.Args()[0]]
	if !found {
		fmt.Fprintf(os.Stderr, "Command %q not found\n", flag.Args()[0])
		os.Exit(2)
	}

	cmdsTmpl := template.Must(template.New("").ParseFiles("internal/license.tmpl", command.inputFilename))

	buf := bytes.NewBuffer(make([]byte, 0))

	err := cmdsTmpl.ExecuteTemplate(buf, filepath.Base(command.inputFilename), struct {
		Copyright string
		Owner     string
		Package   string
		Data      interface{}
	}{
		Copyright:  fmt.Sprintf("%4d", time.Now().Year()),
		Maintainer: maintainer,
		Package:    pkg,
		Data:       command.data(),
	})

	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	f, err := os.Create(command.outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := format.Source(buf.Bytes())
	if err != nil {
		f.Write(buf.Bytes()) // This is here to debug bad format
		log.Fatalf("error formatting: %s", err)
	}

	f.Write(b)
}

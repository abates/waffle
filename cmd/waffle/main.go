package main

import (
	"fmt"
	"github.com/abates/waffle"
	"path/filepath"
)

func main() {
	directories.walk(func(dir string) {
		fmt.Printf("%q\n", dir)
	})
}

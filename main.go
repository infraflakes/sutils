package main

import (
	"fmt"
	"os"

	"github.com/infraflakes/sutils/cmd"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("sutils version: %s\n", version)
		return
	}
	cmd.Execute()
}

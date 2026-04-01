package main

import (
	"fmt"
	"os"

	_ "go.uber.org/automaxprocs"

	"<xpfo{ .ModulePath }xpfo>/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

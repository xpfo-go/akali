package main

import (
	"fmt"
	"os"

	"github.com/xpfo-go/akali/cmd/akali"
)

func main() {
	err := akali.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "execute error:", err)
		os.Exit(1)
	}
}

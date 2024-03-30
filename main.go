package main

import (
	"fmt"
	"github.com/xfpo-go/akali/cmd/akali"
)

func main() {
	err := akali.Execute()
	if err != nil {
		fmt.Println("execute error: ", err.Error())
	}
}

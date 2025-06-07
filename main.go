package main

import (
	"fmt"
	"os"

	"github.com/vvvkkkggg/kubeconomist-core/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}

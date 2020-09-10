package main

import (
	"fmt"
	"os"

	"github.com/sebdah/markdown-toc/cmd"
)

var version = "undefined"

func main() {
	cmd.Version = version
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

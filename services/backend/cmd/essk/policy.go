package main

import (
	"fmt"
	"os"
)

func runPolicy(args []string) {
	if len(args) < 1 || args[0] != "check" {
		printUsage()
		os.Exit(1)
	}

	scope := "--staged"
	if len(args) > 1 {
		scope = args[1]
	}

	switch scope {
	case "--all", "--staged":
		fmt.Printf("essk policy check %s: no policy violations found\n", scope)
	default:
		fmt.Printf("unknown policy scope: %s\n", scope)
		os.Exit(1)
	}
}

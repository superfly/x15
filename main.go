package main

import (
	"os"

	x15cmd "github.com/superfly/x15/cmd"
)

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main() {
	defer x15cmd.HandleExit()

	argsWithoutProg := os.Args[1:]

	x15clicmds := []string{"version", "help"}

	if len(argsWithoutProg) == 0 || find(x15clicmds, argsWithoutProg[0]) {
		x15cli()
	} else {
		// Coming
	}
}

func x15cli() {
	x15cmd.Execute()
}

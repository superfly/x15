package cmd

import "github.com/spf13/cobra"

func AllFlySubCommands() []*cobra.Command {
	return []*cobra.Command{
		{Use: "status", Short: "Get application status"},
		{Use: "deploy", Short: "Deploy an application"},
		{Use: "init", Short: "Initialise a new application"},
	}
}

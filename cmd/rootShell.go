package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var ShellCmd = &cobra.Command{
	Use:   "x15",
	Short: "X15 is a test plane for new ideas",
	Long:  "v0.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		_, x15CmdMap := AllX15SubCommands(cmd)
		allX15Cmds := AllX15AndFlySubCommands(cmd)
		appSuggestions := AppListSuggestions()
		completerSuggestionMap := map[string][]prompt.Suggest{
			"":       {},
			"shell":  CobraCommandToSuggestions(allX15Cmds),
			"status": appSuggestions,
		}
		resp := SuggestionPrompt("> x15 ", shellCommandCompleter(completerSuggestionMap))
		subCommand := resp

		if subCommand == "" {
			return
		}

		if strings.Index(resp, " ") > 0 {
			subCommand = subCommand[0:strings.Index(resp, " ")]
		}

		parsedArgs, err := parseCommandLine(resp)
		if err != nil {
			fmt.Println(err)
			return
		}

		if x15CmdMap[subCommand] == nil {
			RunX15CommandWithArgs(parsedArgs)
			return
		}
	},
}

func shellCommandCompleter(suggestionMap map[string][]prompt.Suggest) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return promptCompleter(suggestionMap, d.Text)
	}
}

func promptCompleter(suggestionMap map[string][]prompt.Suggest, text string) []prompt.Suggest {
	var suggestions []prompt.Suggest
	split := strings.Split(text, " ")
	filterFlags := make([]string, 0, len(split))
	for i, v := range split {
		if !strings.HasPrefix(v, "-") || i == len(split)-1 {
			filterFlags = append(filterFlags, v)
		}
	}
	prev := filterFlags[0] // in git commit -m "hello"  commit is prev
	if len(prev) == len(text) {
		suggestions = suggestionMap["shell"]
		return prompt.FilterContains(suggestions, prev, true)
	}
	curr := filterFlags[1] // in git commit -m "hello"  "hello" is curr
	// if strings.HasPrefix(curr, "--") {
	// 	suggestions = FlagSuggestionsForCommand(prev, "--")
	// } else if strings.HasPrefix(curr, "-") {
	// 	suggestions = FlagSuggestionsForCommand(prev, "-")
	// } else if suggestionMap[prev] != nil {
	// 	suggestions = suggestionMap[prev]
	// }
	return prompt.FilterContains(suggestions, curr, true)
}

func Execute() {
	if err := ShellCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CreateSuggestionMap(cmd *cobra.Command) (map[string][]prompt.Suggest, map[string]*cobra.Command) {
	_, x15CmdMap := AllX15SubCommands(cmd)
	allX15Cmds := AllX15AndFlySubCommands(cmd)
	appSuggestions := AppListSuggestions()

	completerSuggestionMap := map[string][]prompt.Suggest{
		"":       {},
		"shell":  CobraCommandToSuggestions(allX15Cmds),
		"status": appSuggestions,
	}

	return completerSuggestionMap, x15CmdMap
}

func RunX15CommandWithArgs(args []string) {
	var err error
	err = RunInTerminalWithColor("flyctl ", args)
	if err != nil {
		fmt.Println("Command may not exist", err)
	}
	return
}

func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("Unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/c-bata/go-prompt"
)

func SuggestionPrompt(prefix string, completer func(d prompt.Document) []prompt.Suggest) string {
	result := prompt.Input(prefix, completer,
		prompt.OptionTitle(""),
		prompt.OptionHistory([]string{""}),
		prompt.OptionPrefixTextColor(prompt.Yellow), // fine
		//prompt.OptionPreviewSuggestionBGColor(prompt.Yellow),
		//prompt.OptionPreviewSuggestionTextColor(prompt.Yellow),
		prompt.OptionSelectedSuggestionBGColor(prompt.Yellow),
		prompt.OptionSuggestionBGColor(prompt.Yellow),
		prompt.OptionSuggestionTextColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.Blue),
		prompt.OptionDescriptionBGColor(prompt.Black),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionCompletionOnDown(),
		prompt.OptionSwitchKeyBindMode(prompt.CommonKeyBind),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn:  exit,
		}),
	)
	branchName := strings.TrimSpace(result)
	if strings.HasPrefix(result, "origin/") {
		branchName = branchName[7:]
	}
	return branchName
}

type Exit int

func exit(_ *prompt.Buffer) {
	panic(Exit(0))
}

func HandleExit() {
	switch v := recover().(type) {
	case nil:
		return
	case Exit:
		os.Exit(int(v))
	default:
		fmt.Println(v)
		fmt.Println(string(debug.Stack()))
	}
}

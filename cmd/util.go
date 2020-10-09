package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

func RunInTerminalWithColor(cmdName string, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	return RunInTerminalWithColorInDir(cmdName, dir, args)
}

func RunInTerminalWithColorInDir(cmdName string, dir string, args []string) error {
	_, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runtime.GOOS != "windows" {
		cmd.ExtraFiles = []*os.File{w}
	}

	err = cmd.Run()
	return err
}

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

func AllX15SubCommands(rootCmd *cobra.Command) ([]*cobra.Command, map[string]*cobra.Command) {
	bitCmds := rootCmd.Commands()
	bitCmdMap := map[string]*cobra.Command{}
	for _, bitCmd := range bitCmds {
		bitCmdMap[bitCmd.Name()] = bitCmd
	}
	return bitCmds, bitCmdMap
}

type App struct {
	Name    string
	Created string
	Owner   string
}

func AppsList() []App {
	apps := make([]App, 0)
	apps = append(apps, App{Name: "Testing", Created: "Today", Owner: "Dj"})
	apps = append(apps, App{Name: "Probe", Created: "Yesterday", Owner: "Dj"})
	apps = append(apps, App{Name: "Builder", Created: "Two days ago", Owner: "Dj"})
	return apps
}

func AppListSuggestions() []prompt.Suggest {
	apps := AppsList()
	var suggestions []prompt.Suggest
	for _, app := range apps {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        app.Name,
			Description: app.Created + " " + app.Owner,
		})
	}
	return suggestions
}

func AllX15AndFlySubCommands(rootCmd *cobra.Command) (cc []*cobra.Command) {
	x15Cmds, _ := AllX15SubCommands(rootCmd)
	flyCmds := AllFlySubCommands()
	//	commonCommands := CommonCommandsList()
	return concatCopyPreAllocate([][]*cobra.Command{x15Cmds, flyCmds})
}

func concatCopyPreAllocate(slices [][]*cobra.Command) []*cobra.Command {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]*cobra.Command, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func CobraCommandToSuggestions(cmds []*cobra.Command) []prompt.Suggest {
	var suggestions []prompt.Suggest
	for _, branch := range cmds {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        branch.Use,
			Description: branch.Short,
		})
	}
	return suggestions
}

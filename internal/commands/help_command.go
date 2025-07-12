package commands

import (
	"fmt"
	"github.com/D3rise/gator/internal/state"
	"strings"
)

func NewHelpCommand(commands []Command) Command {
	return Command{
		Name:        "help",
		Description: "Displays list of available commands",
		Handler:     newHelpCommandHandler(commands),
	}
}

func newHelpCommandHandler(commands []Command) func(state *state.State, args ...string) error {
	return func(state *state.State, args ...string) error {
		var loggedInText string
		if state.Config.CurrentUserName != "" {
			loggedInText = fmt.Sprintf("You are currently logged in as %s", state.Config.CurrentUserName)
		} else {
			loggedInText = "You are currently not logged in"
		}

		fmt.Println("Welcome to Gator!")
		fmt.Print(loggedInText, "\n\n")

		fmt.Println("Available commands:")

		for _, command := range commands {
			fmt.Printf(" - %s%s: %s\n", command.Name, formatCommandArgs(command.Args), command.Description)
		}

		fmt.Println()
		return nil
	}
}

func formatCommandArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return " <" + strings.Join(args, "> <") + ">"
}

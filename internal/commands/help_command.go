package commands

import (
	"fmt"
	"github.com/D3rise/gator/internal/state"
	"strings"
)

func NewHelpCommand(commands []Command) Command {
	return Command{
		Name:        "help",
		Args:        []string{},
		Description: "Displays list of available commands",
		Handler:     newHelpCommandHandler(commands),
	}
}

func newHelpCommandHandler(commands []Command) func(state *state.State, args ...string) error {
	return func(state *state.State, args ...string) error {
		fmt.Print("Welcome to Gator!\n\n")
		fmt.Println("Available commands:")

		for _, command := range commands {
			fmt.Printf(" - %s %s: %s", command.Name, formatCommandArgs(command.Args), command.Description)
		}

		fmt.Println()
		return nil
	}
}

func formatCommandArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	return "<" + strings.Join(args, "> <") + ">"
}

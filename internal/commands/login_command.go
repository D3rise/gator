package commands

import (
	"fmt"
	"github.com/D3rise/gator/internal/state"
)

func NewLoginCommand() Command {
	return Command{
		Name:        "login",
		Args:        []string{"username"},
		Handler:     loginCommandHandler,
		Description: "Change current user",
	}
}

func loginCommandHandler(state *state.State, args ...string) error {
	newUserName := args[0]
	err := state.Config.SetCurrentUserName(newUserName)
	if err != nil {
		return err
	}

	fmt.Printf("Sucessfully logged in as %s!\n", newUserName)
	return nil
}

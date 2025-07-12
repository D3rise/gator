package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/state"
)

func NewUsersCommand() Command {
	return Command{
		Name:        "users",
		Handler:     usersCommandHandler,
		Description: "List all registered users",
	}
}

func usersCommandHandler(state *state.State, _ ...string) error {
	users, err := state.Queries.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Print(" * ", user.Name)
		if state.Config.CurrentUserName == user.Name {
			fmt.Print(" (current)")
		}
		fmt.Print("\n")
	}

	return nil
}

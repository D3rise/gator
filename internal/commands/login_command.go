package commands

import (
	"context"
	"database/sql"
	"errors"
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
	user, err := state.Queries.GetUserByName(context.Background(), newUserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist")
		}
		return err
	}

	err = state.Config.SetCurrentUserName(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Sucessfully logged in as %s!\n", newUserName)
	return nil
}

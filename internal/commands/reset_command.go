package commands

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/state"
)

func NewResetCommand() Command {
	return Command{
		Name:        "reset",
		Handler:     resetCommandHandler,
		Description: "Resets database state for further tests",
	}
}

func resetCommandHandler(state *state.State, _ ...string) error {
	err := state.Queries.ResetUsersTable(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Successfully reset database state")
	return nil
}

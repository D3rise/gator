package commands

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/state"
	"time"
)

func NewRegisterCommand() Command {
	return Command{
		Name:        "register",
		Args:        []string{"username"},
		Handler:     registerCommandHandler,
		Description: "Register a new user",
	}
}

func registerCommandHandler(state *state.State, args ...string) error {
	username := args[0]
	userExists, err := state.Queries.CheckUserExistenceByName(context.Background(), username)
	if err != nil {
		return err
	}

	if userExists {
		return fmt.Errorf("user already exists")
	}

	newUser, err := state.Queries.CreateUser(
		context.Background(),
		database.CreateUserParams{Name: username, UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true}},
	)
	if err != nil {
		return err
	}

	err = state.Config.SetCurrentUserName(newUser.Name)
	if err != nil {
		return err
	}

	fmt.Println("Successfully registered and logged in as", newUser.Name)

	return nil
}

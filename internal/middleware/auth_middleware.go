package middleware

import (
	"context"
	"fmt"
	"github.com/D3rise/gator/internal/state"
)

var AuthMiddleware = Middleware{
	Handler: authMiddlewareHandler,
}

func authMiddlewareHandler(state *state.State, _ ...string) error {
	if state.Config.CurrentUserName == "" {
		return fmt.Errorf("this command requires authentication! Use login command to authenticate, then try again")
	}

	userExists, err := state.Queries.CheckUserExistenceByName(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error while checking for existence of user: %w", err)
	}

	if !userExists {
		_ = state.Config.SetCurrentUserName("")
		return fmt.Errorf("your authentication data is incorrect! try to login again")
	}

	return nil
}

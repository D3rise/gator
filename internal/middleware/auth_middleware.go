package middleware

import (
	"context"
	"database/sql"
	"errors"
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

	user, err := state.Queries.GetUserByName(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = state.Config.SetCurrentUserName("")
			return fmt.Errorf("your authentication data is incorrect! try to login again")
		}
		return fmt.Errorf("error while checking for existence of user: %w", err)
	}

	_ = state.Config.SetCurrentUserId(user.ID)

	return nil
}

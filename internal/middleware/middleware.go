package middleware

import "github.com/D3rise/gator/internal/state"

type middlewareHandler func(state *state.State, args ...string) error

type Middleware struct {
	Handler middlewareHandler
}

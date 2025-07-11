package commands

import "gator/internal/state"

type cliCommandHandler func(state *state.State, args ...string) error

type Command struct {
	Name        string
	Args        []string
	Handler     cliCommandHandler
	Description string
}

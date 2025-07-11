package state

import "gator/internal/config"

type State struct {
	Config *config.Config
}

func NewState(config *config.Config) *State {
	return &State{
		Config: config,
	}
}

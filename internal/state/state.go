package state

import (
	"github.com/D3rise/gator/internal/config"
	"github.com/D3rise/gator/internal/database"
)

type State struct {
	Config  *config.Config
	Queries *database.Queries
}

func NewState(config *config.Config, queries *database.Queries) *State {
	return &State{
		Config:  config,
		Queries: queries,
	}
}

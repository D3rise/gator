package state

import (
	"github.com/D3rise/gator/internal/config"
	"github.com/D3rise/gator/internal/database"
	"net/http"
)

type State struct {
	Config     *config.Config
	Queries    *database.Queries
	HttpClient *http.Client
}

func NewState(config *config.Config, queries *database.Queries, httpClient *http.Client) *State {
	return &State{
		Config:     config,
		Queries:    queries,
		HttpClient: httpClient,
	}
}

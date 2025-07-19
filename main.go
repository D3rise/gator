package main

import (
	"database/sql"
	"fmt"
	"github.com/D3rise/gator/internal/cli"
	"github.com/D3rise/gator/internal/commands"
	"github.com/D3rise/gator/internal/config"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/middleware"
	"github.com/D3rise/gator/internal/state"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	configPathEnv = "GATOR_CONFIG_PATH"
)

func main() {
	// Initialize config
	conf := initConfig()

	// Initialize database connection
	db := initDB(conf)
	defer db.Close()

	dbQueries := database.New(db)

	// Initialize http client for the state
	httpClient := http.Client{}

	// Initialize state and CLI itself
	appState := state.NewState(&conf, dbQueries, &httpClient)
	appCLI := cli.NewCLI(appState)

	// Register all the required commands to the CLI
	appCLI.Register(commands.NewLoginCommand())
	appCLI.Register(commands.NewRegisterCommand())
	appCLI.Register(commands.NewResetCommand())
	appCLI.Register(commands.NewUsersCommand(), middleware.AuthMiddleware)
	appCLI.Register(commands.NewAggCommand(), middleware.AuthMiddleware)
	appCLI.Register(commands.NewAddFeedCommand(), middleware.AuthMiddleware)
	appCLI.Register(commands.NewFeedsCommand())
	appCLI.Register(commands.NewFollowCommand(), middleware.AuthMiddleware)
	appCLI.Register(commands.NewFollowingCommand(), middleware.AuthMiddleware)
	appCLI.Register(commands.NewUnfollowCommand(), middleware.AuthMiddleware)
	appCLI.Register(commands.NewBrowseCommand(), middleware.AuthMiddleware)

	// Help command must be registered last as
	// it requires list of all registered commands
	appCLI.RegisterDefaultCommand(commands.NewHelpCommand(appCLI.GetCommandList()))

	runCommand(appCLI)
}

func initConfig() config.Config {
	conf, err := config.NewConfig(os.Getenv(configPathEnv))
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

func initDB(conf config.Config) *sql.DB {
	db, err := sql.Open("postgres", conf.DbUrl)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func runCommand(appCLI *cli.CLI) {
	if len(os.Args) <= 1 {
		err := appCLI.RunDefaultCommand(os.Args[1:])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	err := appCLI.RunCommand(command, args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

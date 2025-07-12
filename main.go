package main

import (
	"database/sql"
	"fmt"
	"github.com/D3rise/gator/internal/cli"
	"github.com/D3rise/gator/internal/commands"
	"github.com/D3rise/gator/internal/config"
	"github.com/D3rise/gator/internal/database"
	"github.com/D3rise/gator/internal/state"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	configPathEnv = "CONFIG_PATH"
)

func main() {
	conf := initConfig()
	db := initDB(conf)
	defer db.Close()

	dbQueries := database.New(db)

	appState := state.NewState(&conf, dbQueries)
	appCLI := cli.NewCLI(appState)

	appCLI.Register(commands.NewLoginCommand())
	appCLI.Register(commands.NewRegisterCommand())
	appCLI.Register(commands.NewResetCommand())
	appCLI.Register(commands.NewUsersCommand())

	// Help command must be registered last as it requires list of commands
	appCLI.Register(commands.NewHelpCommand(appCLI.GetCommandList()))

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
		_ = appCLI.RunCommand("help", []string{})
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

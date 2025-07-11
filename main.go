package main

import (
	"fmt"
	"gator/internal/cli"
	"gator/internal/commands"
	"gator/internal/config"
	"gator/internal/state"
	"log"
	"os"
)

const (
	configPathEnv = "gator.yml"
)

func main() {
	conf, err := config.NewConfig(os.Getenv(configPathEnv))
	if err != nil {
		log.Fatal(err)
	}

	appState := state.NewState(&conf)
	appCLI := cli.NewCLI(appState)

	appCLI.Register(commands.NewLoginCommand())
	appCLI.Register(commands.NewHelpCommand(appCLI.GetCommandList()))

	if len(os.Args) <= 1 {
		_ = appCLI.RunCommand("help", []string{})
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	err = appCLI.RunCommand(command, args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

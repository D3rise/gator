package cli

import (
	"fmt"
	"github.com/D3rise/gator/internal/commands"
	"github.com/D3rise/gator/internal/state"
	"maps"
	"slices"
)

type CLI struct {
	state          *state.State
	commands       map[string]commands.Command
	defaultCommand *commands.Command
}

func NewCLI(state *state.State) *CLI {
	return &CLI{state: state, commands: make(map[string]commands.Command)}
}

func (cli *CLI) RegisterDefaultCommand(command commands.Command) {
	cli.defaultCommand = &command
}

func (cli *CLI) Register(command commands.Command) {
	if command.Handler == nil {
		panic("command handler is nil")
	}

	if command.Args == nil {
		command.Args = []string{}
	}

	cli.commands[command.Name] = command
}

func (cli *CLI) RunDefaultCommand(args []string) error {
	command := cli.defaultCommand
	if command == nil {
		return fmt.Errorf("default command is not set")
	}

	if len(args) != len(command.Args) {
		return fmt.Errorf("wrong number of arguments: expected %d arguments, got %d", len(command.Args), len(args))
	}

	return command.Handler(cli.state, args...)
}

func (cli *CLI) RunCommand(name string, args []string) error {
	command, exists := cli.commands[name]
	if !exists {
		return fmt.Errorf("command not found: %s", name)
	}

	if len(args) != len(command.Args) {
		return fmt.Errorf("wrong number of arguments: expected %d arguments, got %d", len(command.Args), len(args))
	}

	return command.Handler(cli.state, args...)
}

func (cli *CLI) GetCommandList() []commands.Command {
	return slices.Collect(maps.Values(cli.commands))
}

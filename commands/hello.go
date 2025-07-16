package commands

import (
	"fmt"
)

// HelloCommand implements the Command interface
type HelloCommand struct {
	BaseCommand
}

func init() {
	// Initialize the command
	cmd := &HelloCommand{
		BaseCommand: BaseCommand{
			name:        "hello",
			description: "Prints a greeting message",
			execute: func() error {
				fmt.Println("Hello, welcome to the Minimal CLI!")
				return nil
			},
			subcommands: []Command{
				&SubHelloCommand{},
			},
		},
	}
	registeredCommands = append(registeredCommands, cmd)
}

// SubHelloCommand implements a subcommand
type SubHelloCommand struct {
	BaseCommand
}

func init() {
	cmd := &SubHelloCommand{
		BaseCommand: BaseCommand{
			name:        "greet",
			description: "Prints a personalized greeting",
			execute: func() error {
				fmt.Println("Hi there! Nice to see you!")
				return nil
			},
			subcommands: nil,
		},
	}
	registeredCommands = append(registeredCommands, cmd)
}

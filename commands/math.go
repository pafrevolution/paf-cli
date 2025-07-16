package commands

import (
	"fmt"
)

// MathCommand implements the Command interface
type MathCommand struct {
	BaseCommand
}

func init() {
	// Initialize the command
	cmd := &MathCommand{
		BaseCommand: BaseCommand{
			name:        "math",
			description: "Performs a simple calculation",
			execute: func() error {
				fmt.Println("2 + 2 = 4")
				return nil
			},
			subcommands: []Command{
				&AddCommand{},
			},
		},
	}
	registeredCommands = append(registeredCommands, cmd)
}

// AddCommand implements a subcommand
type AddCommand struct {
	BaseCommand
}

func init() {
	cmd := &AddCommand{
		BaseCommand: BaseCommand{
			name:        "add",
			description: "Adds two numbers",
			execute: func() error {
				fmt.Println("Enter two numbers to add (e.g., 3 4):")
				var a, b float64
				fmt.Scan(&a, &b)
				fmt.Printf("%.2f + %.2f = %.2f\n", a, b, a+b)
				return nil
			},
			subcommands: nil,
		},
	}
	registeredCommands = append(registeredCommands, cmd)
}
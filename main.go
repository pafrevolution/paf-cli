package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/pafrevolution/paf-cli/commands"
)

// Command defines the interface for pluggable commands
type Command interface {
	Name() string
	Description() string
	Execute() error
	Subcommands() []Command
}

// BaseCommand provides a default implementation for commands
type BaseCommand struct {
	name        string
	description string
	execute     func() error
	subcommands []Command
}

func (c BaseCommand) Name() string           { return c.name }
func (c BaseCommand) Description() string    { return c.description }
func (c BaseCommand) Execute() error         { return c.execute() }
func (c BaseCommand) Subcommands() []Command { return c.subcommands }

// MenuModel for the Bubble Tea TUI
type MenuModel struct {
	choices     []Command
	cursor      int
	selected    map[int]struct{}
	currentMenu Command
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.choices[m.cursor].Subcommands()) > 0 {
				// Navigate to submenu
				m.selected[m.cursor] = struct{}{}
				m.currentMenu = m.choices[m.cursor]
				m.choices = m.choices[m.cursor].Subcommands()
				m.cursor = 0
			} else {
				// Execute command
				if err := m.choices[m.cursor].Execute(); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
				return m, tea.Quit
			}
		case "b":
			// Go back to parent menu if exists
			if m.currentMenu != nil {
				m.choices = findParentCommands(m.currentMenu, registeredCommands)
				m.cursor = 0
			}
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString("\033[1;34m=== Minimal CLI ===\033[0m\n\n")
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "\033[1;32m>>\033[0m"
		}
		s.WriteString(fmt.Sprintf("%s \033[1;36m%s\033[0m: %s\n", cursor, choice.Name(), choice.Description()))
	}
	s.WriteString("\n\033[1;33mUse ↑/↓ to navigate, Enter to select, b to go back, q to quit\033[0m\n")
	return s.String()
}

// Registered commands
var registeredCommands []Command

// Find parent commands for navigation
func findParentCommands(current Command, commands []Command) []Command {
	for _, cmd := range commands {
		for _, sub := range cmd.Subcommands() {
			if sub.Name() == current.Name() {
				return commands
			}
			if subcommands := findParentCommands(current, sub.Subcommands()); subcommands != nil {
				return subcommands
			}
		}
	}
	return registeredCommands
}

// Register commands dynamically
func registerCommands() {
	// Register all commands from the commands package
	for _, cmd := range commands.RegisteredCommands {
		registeredCommands = append(registeredCommands, cmd)
	}
}

func main() {
	// Initialize commands
	registerCommands()

	// Cobra root command
	rootCmd := &cobra.Command{
		Use:   "minicli",
		Short: "A minimalistic CLI with dynamic commands",
		Run: func(cmd *cobra.Command, args []string) {
			// Start TUI
			m := MenuModel{
				choices:     registeredCommands,
				selected:    make(map[int]struct{}),
				currentMenu: nil,
			}
			p := tea.NewProgram(m)
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Dynamically load commands from files
	err := filepath.Walk("./commands", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			// Placeholder for dynamic command loading
			// In a real implementation, use plugin or source code parsing
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading commands: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

package cmd

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var greetCmd = &cobra.Command{
    Use:   "greet",
    Short: "Greet a person",
    Run: func(cmd *cobra.Command, args []string) {
        name := "World"
        if len(args) > 0 {
            name = args[0]
        }
        fmt.Printf("Hello, %s!\n", name)
    },
}

func init() {
    // Add the greet command to the root command
    rootCmd.AddCommand(greetCmd)
}
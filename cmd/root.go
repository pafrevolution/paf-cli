package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "paf",
		Short: "A simple CLI tool",
		Long: `A flexible CLI application built with Go`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to PAF CLI!")
		},
	}
	verbose bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func Execute() error {
	return rootCmd.Execute()
}
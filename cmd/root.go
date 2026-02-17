package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cire",
	Short: "Cire - Convenient Wire generator",
	Long: `Cire is a CLI tool that generates wire.go from struct dependencies in Go projects.
It analyzes struct dependencies and generates Wire injection code automatically.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

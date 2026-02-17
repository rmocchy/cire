package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "convinient_wire",
	Short: "Wire dependency analyzer - analyzes struct dependencies",
	Long: `Convinient Wire is a CLI tool that analyzes struct dependencies in Go projects.
It uses the internal/analyze package to build a dependency tree and outputs the result in YAML format.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

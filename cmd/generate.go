package cmd

import (
	"github.com/rmocchy/cire/internal/app"
	"github.com/spf13/cobra"
)

var (
	filePath string
	genJson  bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate wire.go from struct dependencies",
	Long: `Analyze structs defined in a file with //go:build cire tag and generate wire.go file.
The target file must have the build tag "//go:build cire" and contain struct definitions.`,
	Example: `  cire generate --file ./cire.go
  cire generate -f ./cire.go --yaml`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&filePath, "file", "f", "", "Go file path with //go:build cire tag containing struct definitions (required)")
	generateCmd.Flags().BoolVarP(&genJson, "json", "j", false, "Generate YAML file in the same directory as the input file")

	generateCmd.MarkFlagRequired("file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	input := app.GenerateInput{
		FilePath: filePath,
		GenJson:  genJson,
	}
	return app.RunGenerate(&input)
}

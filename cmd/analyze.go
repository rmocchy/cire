package cmd

import (
	"fmt"

	pipe "github.com/rmocchy/convinient_wire/internal/analyze"
	"github.com/rmocchy/convinient_wire/internal/core"
	"github.com/rmocchy/convinient_wire/internal/load"
	"github.com/spf13/cobra"
)

var (
	packagePath string
	structName  string
	outputFile  string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze struct dependencies and output to YAML",
	Long: `Analyze a specific struct's dependencies and output the dependency tree in YAML format.
You can specify the package path and struct name to analyze.`,
	Example: `  convinient_wire analyze --package github.com/example/myapp/service --struct UserService
  convinient_wire analyze --struct UserService --output dependencies.yaml`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&packagePath, "package", "p", "", "Package path where the struct is defined")
	analyzeCmd.Flags().StringVarP(&structName, "struct", "s", "", "Struct name to analyze (required)")
	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")

	analyzeCmd.MarkFlagRequired("struct")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// パッケージのロード
	loader, err := load.LoadPackages(packagePath)
	if err != nil {
		return err
	}

	// アナライザの作成
	analyzer, err := pipe.NewWireAnalyzer(loader.FunctionCache, loader.StructCache)
	if err != nil {
		return fmt.Errorf("failed to create analyzer: %w", err)
	}

	// 構造体の解析
	var pkgPath core.PackagePath
	if packagePath != "" {
		pkgPath = core.NewPackagePath(packagePath)
	}

	result, err := analyzer.AnalyzeStruct(structName, pkgPath)
	if err != nil {
		return fmt.Errorf("failed to analyze struct: %w", err)
	}

	// YAML形式で出力
	if err := pipe.OutputToYAML(result, outputFile); err != nil {
		return err
	}

	return nil
}

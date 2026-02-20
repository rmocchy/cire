package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmocchy/cire/internal/analyze"
	pipe "github.com/rmocchy/cire/internal/analyze"
	"github.com/rmocchy/cire/internal/load"
	"github.com/rmocchy/cire/internal/yaml"
	"github.com/spf13/cobra"
)

var (
	filePath     string
	generateYAML bool
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
	generateCmd.Flags().BoolVarP(&generateYAML, "yaml", "y", false, "Generate YAML file in the same directory as the input file")

	generateCmd.MarkFlagRequired("file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	targetStructs, err := load.FindAnnotatedStructs(filePath)
	if err != nil {
		return err
	}

	// パッケージパスを解決
	pkgPath, err := load.ResolvePackagePath(filePath)
	if err != nil {
		return err
	}

	// パッケージのロード
	loader, err := load.LoadPackagesFromFile(filePath)
	if err != nil {
		return err
	}

	// 各アノテーション付き構造体を解析
	results := make([]*analyze.StructNode, 0, len(targetStructs))
	for _, structName := range targetStructs {
		// アナライザの作成
		analyzer, err := pipe.NewWireAnalyzer(loader.FunctionCache, loader.StructCache)
		if err != nil {
			return fmt.Errorf("failed to create analyzer: %w", err)
		}

		// 構造体の解析
		result, err := analyzer.AnalyzeStruct(structName, *pkgPath)
		if err != nil {
			return fmt.Errorf("failed to analyze struct %s: %w", structName, err)
		}

		results = append(results, result)
	}

	// YAMLフラグが指定されている場合のみYAML生成
	if generateYAML {
		// 入力ファイルと同じディレクトリにYAMLファイルを出力
		dir := filepath.Dir(filePath)
		outputPath := filepath.Join(dir, "cire.yaml")

		if err := yaml.OutputMultipleToYAML(results, outputPath); err != nil {
			return err
		}

		fmt.Printf("YAML file generated: %s\n", outputPath)
	}

	// wire.go ファイルの生成
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write wire.go: %w", err)
	}

	fmt.Printf("Wire file generated: %s\n", outputPath)
	return nil

	return nil
}

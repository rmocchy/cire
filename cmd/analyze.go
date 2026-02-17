package cmd

import (
	"fmt"
	"path/filepath"

	pipe "github.com/rmocchy/convinient_wire/internal/analyze"
	"github.com/rmocchy/convinient_wire/internal/load"
	wiregenerate "github.com/rmocchy/convinient_wire/internal/wire_generate"
	"github.com/spf13/cobra"
)

var (
	filePath     string
	generateYAML bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze struct dependencies and output to YAML",
	Long: `Analyze structs defined in a file with //go:build cire tag and output the dependency tree in YAML format.
The target file must have the build tag "//go:build cire" and contain struct definitions.`,
	Example: `  convinient_wire analyze --file ./cire_structs.go --yaml
  convinient_wire analyze -f ./cire_structs.go -y`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&filePath, "file", "f", "", "Go file path with //go:build cire tag containing struct definitions (required)")
	analyzeCmd.Flags().BoolVarP(&generateYAML, "yaml", "y", false, "Generate YAML file in the same directory as the input file")

	analyzeCmd.MarkFlagRequired("file")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
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
	results := make([]*pipe.StructNode, 0, len(targetStructs))
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

		if err := pipe.OutputMultipleToYAML(results, outputPath); err != nil {
			return err
		}

		fmt.Printf("YAML file generated: %s\n", outputPath)
	}

	// wire.go ファイルの生成
	if err := wiregenerate.GenerateWireFile(results, filePath); err != nil {
		return fmt.Errorf("failed to generate wire.go: %w", err)
	}

	return nil
}

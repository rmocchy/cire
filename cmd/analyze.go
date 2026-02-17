package cmd

import (
	"fmt"

	pipe "github.com/rmocchy/convinient_wire/internal/analyze"
	"github.com/rmocchy/convinient_wire/internal/core"
	"github.com/rmocchy/convinient_wire/internal/load"
	"github.com/spf13/cobra"
)

var (
	filePath   string
	outputFile string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze struct dependencies and output to YAML",
	Long: `Analyze structs with @cire annotation and output the dependency tree in YAML format.
The target struct must have a comment with @cire annotation.`,
	Example: `  convinient_wire analyze --file ./handler/user_handler.go --output dependencies.yaml
  convinient_wire analyze -f ./service/user_service.go`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&filePath, "file", "f", "", "Go file path containing the struct with @cire annotation (required)")
	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")

	analyzeCmd.MarkFlagRequired("file")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	// ファイルからアノテーション付き構造体を検出
	annotatedStructs, err := load.FindAnnotatedStructs(filePath)
	if err != nil {
		return err
	}

	// パッケージパスを解決
	pkgPath, err := load.ResolvePackagePath(filePath)
	if err != nil {
		return err
	}

	// 最初の構造体のパッケージパスを設定
	for i := range annotatedStructs {
		annotatedStructs[i].PackagePath = pkgPath
	}

	// パッケージのロード
	loader, err := load.LoadPackagesFromFile(filePath)
	if err != nil {
		return err
	}

	// 各アノテーション付き構造体を解析
	for _, annotated := range annotatedStructs {
		// アナライザの作成
		analyzer, err := pipe.NewWireAnalyzer(loader.FunctionCache, loader.StructCache)
		if err != nil {
			return fmt.Errorf("failed to create analyzer: %w", err)
		}

		// 構造体の解析
		corePkgPath := core.NewPackagePath(annotated.PackagePath)
		result, err := analyzer.AnalyzeStruct(annotated.Name, corePkgPath)
		if err != nil {
			return fmt.Errorf("failed to analyze struct %s: %w", annotated.Name, err)
		}

		// YAML形式で出力
		if err := pipe.OutputToYAML(result, outputFile); err != nil {
			return err
		}
	}

	return nil
}

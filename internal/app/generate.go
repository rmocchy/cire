package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmocchy/cire/internal/analyze"
	"github.com/rmocchy/cire/internal/file"
	"github.com/rmocchy/cire/internal/generate"
)

type GenerateInput struct {
	FilePath string
	GenJson  bool
}

func RunGenerate(input *GenerateInput) error {
	dir := filepath.Dir(input.FilePath)
	outputPath := filepath.Join(dir, "wire.go")

	pkgs, err := file.LoadAllPkgsFromPath(input.FilePath)
	if err != nil {
		return err
	}

	structs, err := file.LoadNamedStructs(input.FilePath, pkgs)
	if err != nil {
		return err
	}

	// キャッシュの準備
	fnCache := analyze.NewFunctionCache(pkgs)
	anCache := analyze.NewAnalysisCache()
	analyzer := analyze.NewAnalyze(fnCache, anCache)

	// コード生成の準備
	config := &generate.GenerateConfig{}
	usePkgName, err := file.ExtractPackageName(input.FilePath)
	if err != nil {
		return err
	}
	config.SetPackageName(*usePkgName)

	// 構造体ごとに解析実行
	for _, s := range structs {
		trees, err := analyzer.ExecuteFromStruct(s)
		if err != nil {
			return err
		}
		converter := analyze.NewConvertTreeToUniqueList()
		for _, tree := range trees {
			converter.Execute(tree)
		}

		// 依存関係から生成可能かどうかをチェック
		verr := analyze.IsDepTreeSatisfiable(converter.List())

		if input.GenJson || verr != nil {
			jsonConfig := &analyze.JsonConfig{
				Dir:        dir,
				StructName: s.Obj().Name(),
				Data:       trees,
			}
			if err := analyze.WriteOnJsonFile(jsonConfig); err != nil {
				return err
			}
			if verr != nil {
				return fmt.Errorf("dependency tree is not satisfiable for struct %s: %w", s.Obj().Name(), verr)
			}
		}

		providers := make([]generate.Provider, 0, len(converter.List()))
		for _, node := range converter.List() {
			providers = append(providers, generate.Provider{
				PkgPath: node.PkgPath,
				Name:    fmt.Sprintf("%s.%s", file.PkgNameFromPath(node.PkgPath), node.Name),
			})
		}
		config.AddStructSet(s.Obj().Name(), providers)
	}

	// コード生成
	formatted, err := config.Generate()
	if err != nil {
		return err
	}

	// 結果の出力
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write wire.go: %w", err)
	}

	fmt.Printf("Wire file generated: %s\n", outputPath)
	return nil
}

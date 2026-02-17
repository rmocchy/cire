package wiregenerate

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	pipe "github.com/rmocchy/convinient_wire/internal/analyze"
)

// WireData は wire.go テンプレートに渡すデータ
type WireData struct {
	PackageName  string
	Imports      []string
	ProviderSets []ProviderSetData
}

// ProviderSetData は各 Provider セットのデータ
type ProviderSetData struct {
	StructName string
	Providers  []string
}

// GenerateWireFile は result から initFunction を取得して wire.go を生成する
func GenerateWireFile(results []*pipe.StructNode, inputFilePath string) error {
	// 入力ファイルと同じディレクトリに wire.go を出力
	dir := filepath.Dir(inputFilePath)
	outputPath := filepath.Join(dir, "wire.go")

	// パッケージ名を取得（ディレクトリ名を使用）
	packageName := filepath.Base(dir)

	// インポートパスとinitFunctionを収集
	importMap := make(map[string]bool)
	providerSets := []ProviderSetData{}

	for _, result := range results {
		providers := []string{}

		// Dependencies から initFunction を取得
		for _, dep := range result.Dependencies {
			if structNode, ok := dep.(*pipe.StructNode); ok {
				for _, initFunc := range structNode.InitFunctions {
					// パッケージパスとインポートを収集
					pkg := initFunc.Pkg()
					if pkg != nil {
						pkgPath := pkg.Path()
						importMap[pkgPath] = true

						// 関数名をプロバイダーリストに追加
						funcName := initFunc.Name()
						// パッケージ名を含める場合
						pkgName := pkg.Name()
						providers = append(providers, pkgName+"."+funcName)
					}
				}
			}
		}

		// 自分自身の InitFunctions も追加
		for _, initFunc := range result.InitFunctions {
			pkg := initFunc.Pkg()
			if pkg != nil {
				pkgPath := pkg.Path()
				importMap[pkgPath] = true

				funcName := initFunc.Name()
				pkgName := pkg.Name()
				providers = append(providers, pkgName+"."+funcName)
			}
		}

		if len(providers) > 0 {
			providerSets = append(providerSets, ProviderSetData{
				StructName: result.StructName,
				Providers:  providers,
			})
		}
	}

	// インポートリストを作成
	imports := []string{}
	for imp := range importMap {
		imports = append(imports, imp)
	}

	// テンプレートデータを準備
	data := WireData{
		PackageName:  packageName,
		Imports:      imports,
		ProviderSets: providerSets,
	}

	// テンプレートを定義
	tmpl := template.Must(template.New("wire").Parse(wireTemplate))

	// テンプレートを実行
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// go/format でフォーマット
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// フォーマットに失敗した場合は元のコードを出力
		formatted = buf.Bytes()
	}

	// ファイルに書き込み
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write wire.go: %w", err)
	}

	fmt.Printf("Wire file generated: %s\n", outputPath)
	return nil
}

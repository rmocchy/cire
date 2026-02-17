package wiregenerate

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
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

// GenerateWireFile は解析結果から wire.go を生成する
func GenerateWireFile(rootStructs []*pipe.StructNode, inputFilePath string) error {
	dir := filepath.Dir(inputFilePath)
	outputPath := filepath.Join(dir, "wire.go")

	// 入力ファイルからパッケージ名を取得
	packageName, err := extractPackageName(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to extract package name: %w", err)
	}

	// インポート収集用のmap
	importMap := make(map[string]bool)

	// プロバイダーセットを収集
	providerSets := collectProviderSets(rootStructs, importMap)

	// インポートリストを作成
	imports := make([]string, 0, len(importMap))
	for imp := range importMap {
		imports = append(imports, imp)
	}

	// テンプレートを実行してファイルに書き込み
	return writeWireFile(outputPath, WireData{
		PackageName:  packageName,
		Imports:      imports,
		ProviderSets: providerSets,
	})
}

// writeWireFile はテンプレートを実行して wire.go を書き込む
func writeWireFile(outputPath string, data WireData) error {
	tmpl := template.Must(template.New("wire").Parse(wireTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}

	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write wire.go: %w", err)
	}

	fmt.Printf("Wire file generated: %s\n", outputPath)
	return nil
}

// extractPackageName は指定されたGoファイルからパッケージ名を取得する
func extractPackageName(filePath string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}
	return f.Name.Name, nil
}

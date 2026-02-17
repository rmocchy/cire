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
	Structs      []StructDef
}

// ProviderSetData は各 Provider セットのデータ
type ProviderSetData struct {
	StructName string
	Providers  []string
}

// StructDef は構造体定義のデータ
type StructDef struct {
	Name   string
	Fields []StructFieldDef
}

// StructFieldDef は構造体フィールドのデータ
type StructFieldDef struct {
	Name    string
	Type    string
	Pointer bool
}

// GenerateWireFile は解析結果から wire.go を生成する
func GenerateWireFile(results []*pipe.StructNode, inputFilePath string) error {
	dir := filepath.Dir(inputFilePath)
	outputPath := filepath.Join(dir, "wire.go")
	packageName := filepath.Base(dir)

	// インポート収集用のmap
	importMap := make(map[string]bool)

	// 構造体定義とプロバイダーセットを収集
	structs := collectStructDefs(results, importMap)
	providerSets := collectProviderSets(results, importMap)

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
		Structs:      structs,
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

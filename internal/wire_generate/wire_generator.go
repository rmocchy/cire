package wiregenerate

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
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
	structs := []StructDef{}

	for _, result := range results {
		// 重複排除用のmap
		providerMap := make(map[string]bool)
		providers := []string{}

		// Fieldsから再帰的にinitFunctionを収集
		collectInitFunctions(result.Fields, importMap, providerMap, &providers)

		// 自分自身の InitFunctions も追加
		for _, initFunc := range result.InitFunctions {
			addInitFunction(initFunc, importMap, providerMap, &providers)
		}

		if len(providers) > 0 {
			providerSets = append(providerSets, ProviderSetData{
				StructName: result.StructName,
				Providers:  providers,
			})
		}

		// 構造体定義を収集
		structDef := StructDef{
			Name:   result.StructName,
			Fields: []StructFieldDef{},
		}
		for _, field := range result.Fields {
			fieldDef := convertFieldToStructFieldDef(field, importMap)
			if fieldDef != nil {
				structDef.Fields = append(structDef.Fields, *fieldDef)
			}
		}
		structs = append(structs, structDef)
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
		Structs:      structs,
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

// convertFieldToStructFieldDef はFieldNodeをStructFieldDefに変換する
func convertFieldToStructFieldDef(field pipe.FieldNode, importMap map[string]bool) *StructFieldDef {
	switch f := field.(type) {
	case *pipe.StructNode:
		// パッケージをインポートに追加
		if f.PackagePath != "" {
			importMap[f.PackagePath] = true
		}
		// パッケージ名を取得
		pkgName := filepath.Base(f.PackagePath)
		typeName := pkgName + "." + f.StructName
		return &StructFieldDef{
			Name:    f.FieldName,
			Type:    typeName,
			Pointer: true, // 構造体は通常ポインタで持つ
		}
	case *pipe.InterfaceNode:
		// インターフェースの場合
		if f.PackagePath != "" {
			importMap[f.PackagePath] = true
		}
		pkgName := filepath.Base(f.PackagePath)
		typeName := pkgName + "." + f.TypeName
		return &StructFieldDef{
			Name:    f.FieldName,
			Type:    typeName,
			Pointer: false, // インターフェースはポインタなし
		}
	case *pipe.BuiltinNode:
		return &StructFieldDef{
			Name:    f.FieldName,
			Type:    f.TypeName,
			Pointer: false,
		}
	}
	return nil
}

// collectInitFunctions はFieldNodeから再帰的にinitFunctionを収集する
func collectInitFunctions(fields []pipe.FieldNode, importMap map[string]bool, providerMap map[string]bool, providers *[]string) {
	for _, field := range fields {
		switch f := field.(type) {
		case *pipe.StructNode:
			// まず依存関係を再帰的に収集
			collectInitFunctions(f.Dependencies, importMap, providerMap, providers)
			// 自身のInitFunctionsを追加
			for _, initFunc := range f.InitFunctions {
				addInitFunction(initFunc, importMap, providerMap, providers)
			}
		case *pipe.InterfaceNode:
			// まず依存関係を再帰的に収集
			collectInitFunctions(f.Dependencies, importMap, providerMap, providers)
			// 自身のInitFunctionsを追加
			for _, initFunc := range f.InitFunctions {
				addInitFunction(initFunc, importMap, providerMap, providers)
			}
		}
	}
}

// addInitFunction はinitFunctionをprovidersに追加する（重複チェック付き）
func addInitFunction(initFunc *types.Func, importMap map[string]bool, providerMap map[string]bool, providers *[]string) {
	pkg := initFunc.Pkg()
	if pkg == nil {
		return
	}

	pkgPath := pkg.Path()
	funcName := initFunc.Name()
	pkgName := pkg.Name()
	fullName := pkgName + "." + funcName

	// 重複チェック
	if providerMap[fullName] {
		return
	}

	importMap[pkgPath] = true
	providerMap[fullName] = true
	*providers = append(*providers, fullName)
}

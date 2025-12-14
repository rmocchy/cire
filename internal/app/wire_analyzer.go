package app

import (
	"fmt"

	file "github.com/rmocchy/convinient_wire/internal/files"
	"github.com/rmocchy/convinient_wire/internal/packages"
	gopkgs "golang.org/x/tools/go/packages"
)

// WireAnalyzer はwire.goの解析を行う
type WireAnalyzer struct {
	workDir       string
	searchPattern string
	analyzed      map[string]*StructNode // 解析済みの構造体をキャッシュ（無限ループ防止）
}

// NewWireAnalyzer は新しいWireAnalyzerを作成する
func NewWireAnalyzer(workDir, searchPattern string) *WireAnalyzer {
	return &WireAnalyzer{
		workDir:       workDir,
		searchPattern: searchPattern,
		analyzed:      make(map[string]*StructNode),
	}
}

// AnalyzeWireFile はwire.goファイルを解析する
func (wa *WireAnalyzer) AnalyzeWireFile(wireFilePath string) ([]*StructNode, error) {
	// wire.goから構造体を取得
	functions, err := file.ParseWireFileStructs(wireFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wire file: %w", err)
	}

	var results []*StructNode

	// 各関数の返り値構造体を解析
	for _, funcInfo := range functions {
		for _, structInfo := range funcInfo.ReturnTypes {
			// 構造体を再帰的に解析
			structNode, err := wa.analyzeStruct("", structInfo.Name)
			if err != nil {
				// エラーがあっても他の構造体の解析を続ける
				results = append(results, &StructNode{
					StructName: structInfo.Name,
					Skipped:    true,
					SkipReason: fmt.Sprintf("failed to analyze: %v", err),
				})
				continue
			}
			results = append(results, structNode)
		}
	}

	return results, nil
}

// analyzeStruct は構造体を再帰的に解析する
func (wa *WireAnalyzer) analyzeStruct(packagePath, structName string) (*StructNode, error) {
	// キャッシュキーを生成
	cacheKey := packagePath + "." + structName
	if packagePath == "" {
		cacheKey = structName
	}

	// 既に解析済みの場合はキャッシュから返す
	if cached, ok := wa.analyzed[cacheKey]; ok {
		return cached, nil
	}

	// 構造体のフィールド情報を取得
	fieldsInfo, err := packages.ExtractStructFields(wa.workDir, packagePath, structName)
	if err != nil {
		return nil, fmt.Errorf("failed to extract struct fields for %s: %w", structName, err)
	}

	result := &StructNode{
		StructName:    structName,
		PackagePath:   packagePath,
		InitFunctions: make([]InitFunctionInfo, 0),
		Fields:        make([]FieldNode, 0, len(fieldsInfo.Fields)),
	}

	// キャッシュに登録（無限ループ防止のため、フィールド解析前に登録）
	wa.analyzed[cacheKey] = result

	// 初期化関数を探す
	initFuncs, err := wa.findInitFunctions(packagePath, structName)
	if err == nil {
		result.InitFunctions = initFuncs
	}

	// 各フィールドを解析
	for _, field := range fieldsInfo.Fields {
		fieldNode := wa.analyzeField(field)
		if fieldNode != nil {
			result.Fields = append(result.Fields, fieldNode)
		}
	}

	return result, nil
}

// findInitFunctions は構造体を返す初期化関数を探す
func (wa *WireAnalyzer) findInitFunctions(packagePath, structName string) ([]InitFunctionInfo, error) {
	// パッケージを読み込む
	cfg := &gopkgs.Config{
		Mode: gopkgs.NeedName | gopkgs.NeedFiles | gopkgs.NeedImports |
			gopkgs.NeedDeps | gopkgs.NeedTypes | gopkgs.NeedSyntax | gopkgs.NeedTypesInfo,
		Dir: wa.workDir,
	}

	pkgs, err := gopkgs.Load(cfg, wa.searchPattern)
	if err != nil {
		return nil, err
	}

	// 構造体を返す関数を探す
	functions := packages.FindFunctionsReturningStruct(structName, packagePath, pkgs)

	// InitFunctionInfoに変換
	initFuncs := make([]InitFunctionInfo, 0, len(functions))
	for _, fn := range functions {
		initFuncs = append(initFuncs, InitFunctionInfo{
			Name:        fn.Name,
			PackagePath: fn.PackagePath,
		})
	}

	return initFuncs, nil
}

// analyzeField はフィールドを解析する
func (wa *WireAnalyzer) analyzeField(field packages.FieldInfo) FieldNode {
	// インターフェース型の場合
	if field.IsInterface {
		resolvedStruct, skipReason := wa.resolveInterface(field)
		return &InterfaceNode{
			FieldName:      field.Name,
			TypeName:       field.TypeName,
			PackagePath:    field.PackagePath,
			ResolvedStruct: resolvedStruct,
			Skipped:        skipReason != "",
			SkipReason:     skipReason,
		}
	}

	// 構造体型の場合
	if field.TypeName != "" && field.PackagePath != "" && !isBuiltinType(field.TypeName) {
		resolvedStruct, err := wa.analyzeStruct(field.PackagePath, field.TypeName)
		if err != nil {
			// エラーの場合はnilを返す（スキップ）
			return nil
		}
		// FieldNameをセットして返す
		resolvedStruct.FieldName = field.Name
		return resolvedStruct
	}

	// 基本型などはスキップ（nilを返す）
	return nil
}

// resolveInterface はインターフェースから具体的な構造体を解決する
func (wa *WireAnalyzer) resolveInterface(field packages.FieldInfo) (*StructNode, string) {
	// インターフェースを参照する関数を検索
	refs, err := packages.FindInterfaceReferences(
		wa.workDir,
		field.TypeName,
		field.PackagePath,
		wa.searchPattern,
	)
	if err != nil {
		return nil, fmt.Sprintf("failed to find interface references: %v", err)
	}

	// 参照が見つからない場合
	if len(refs) == 0 {
		return nil, "no implementing types found"
	}

	// 複数の実装がある場合はスキップ
	if len(refs) > 1 {
		return nil, fmt.Sprintf("multiple implementing types found (%d)", len(refs))
	}

	// 実装型を再帰的に解析
	ref := refs[0]
	resolvedStruct, err := wa.analyzeStruct(ref.ImplementingPkgPath, ref.ImplementingType)
	if err != nil {
		return nil, fmt.Sprintf("failed to analyze implementing type: %v", err)
	}

	return resolvedStruct, ""
}

// isBuiltinType はビルトイン型かどうかを判定する
func isBuiltinType(typeName string) bool {
	builtinTypes := map[string]bool{
		"string":     true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"float32":    true,
		"float64":    true,
		"bool":       true,
		"byte":       true,
		"rune":       true,
		"error":      true,
		"complex64":  true,
		"complex128": true,
	}
	return builtinTypes[typeName]
}

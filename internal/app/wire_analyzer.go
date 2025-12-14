package app

import (
	"fmt"

	file "github.com/rmocchy/convinient_wire/internal/files"
	"github.com/rmocchy/convinient_wire/internal/packages"
)

// WireAnalyzer はwire.goの解析を行う
type WireAnalyzer struct {
	workDir       string
	searchPattern string
	analyzed      map[string]*StructAnalysisResult // 解析済みの構造体をキャッシュ（無限ループ防止）
}

// NewWireAnalyzer は新しいWireAnalyzerを作成する
func NewWireAnalyzer(workDir, searchPattern string) *WireAnalyzer {
	return &WireAnalyzer{
		workDir:       workDir,
		searchPattern: searchPattern,
		analyzed:      make(map[string]*StructAnalysisResult),
	}
}

// AnalyzeWireFile はwire.goファイルを解析する
func (wa *WireAnalyzer) AnalyzeWireFile(wireFilePath string) ([]StructAnalysisResult, error) {
	// wire.goから構造体を取得
	functions, err := file.ParseWireFileStructs(wireFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wire file: %w", err)
	}

	var results []StructAnalysisResult

	// 各関数の返り値構造体を解析
	for _, funcInfo := range functions {
		for _, structInfo := range funcInfo.ReturnTypes {
			// 構造体を再帰的に解析
			result, err := wa.analyzeStruct("", structInfo.Name)
			if err != nil {
				// エラーがあっても他の構造体の解析を続ける
				results = append(results, StructAnalysisResult{
					StructName: structInfo.Name,
					Skipped:    true,
					SkipReason: fmt.Sprintf("failed to analyze: %v", err),
				})
				continue
			}
			results = append(results, *result)
		}
	}

	return results, nil
}

// analyzeStruct は構造体を再帰的に解析する
func (wa *WireAnalyzer) analyzeStruct(packagePath, structName string) (*StructAnalysisResult, error) {
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

	result := &StructAnalysisResult{
		StructName:  structName,
		PackagePath: packagePath,
		Fields:      make([]FieldAnalysisResult, 0, len(fieldsInfo.Fields)),
	}

	// キャッシュに登録（無限ループ防止のため、フィールド解析前に登録）
	wa.analyzed[cacheKey] = result

	// 各フィールドを解析
	for _, field := range fieldsInfo.Fields {
		fieldResult := wa.analyzeField(field)
		result.Fields = append(result.Fields, fieldResult)
	}

	return result, nil
}

// analyzeField はフィールドを解析する
func (wa *WireAnalyzer) analyzeField(field packages.FieldInfo) FieldAnalysisResult {
	result := FieldAnalysisResult{
		Name:        field.Name,
		TypeName:    field.TypeName,
		PackagePath: field.PackagePath,
		IsPointer:   field.IsPointer,
		IsInterface: field.IsInterface,
	}

	// インターフェース型の場合、具体的な構造体に解決を試みる
	if field.IsInterface {
		resolvedStruct, skipReason := wa.resolveInterface(field)
		if resolvedStruct != nil {
			result.ResolvedStruct = resolvedStruct
		} else if skipReason != "" {
			result.InterfaceSkipped = true
			result.InterfaceSkipReason = skipReason
		}
	} else if field.TypeName != "" && field.PackagePath != "" {
		// 通常の構造体型の場合、再帰的に解析
		// ただし、基本型やビルトイン型は除外
		if !isBuiltinType(field.TypeName) {
			resolvedStruct, err := wa.analyzeStruct(field.PackagePath, field.TypeName)
			if err != nil {
				// エラーの場合は解析をスキップ
				result.InterfaceSkipped = true
				result.InterfaceSkipReason = fmt.Sprintf("failed to analyze: %v", err)
			} else {
				result.ResolvedStruct = resolvedStruct
			}
		}
	}

	return result
}

// resolveInterface はインターフェースから具体的な構造体を解決する
func (wa *WireAnalyzer) resolveInterface(field packages.FieldInfo) (*StructAnalysisResult, string) {
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

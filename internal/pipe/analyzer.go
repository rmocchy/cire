package pipe

import (
	"fmt"
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
	"golang.org/x/tools/go/packages"
)

// WireAnalyzer はwire.goの解析を行う（internal/coreを使用）
type WireAnalyzer struct {
	workDir       string
	searchPattern string
	analyzed      map[string]*StructNode // 解析済みの構造体をキャッシュ（無限ループ防止）
	pkgs          []*packages.Package    // ロード済みのパッケージ
}

// NewWireAnalyzer は新しいWireAnalyzerを作成する
func NewWireAnalyzer(workDir, searchPattern string) (*WireAnalyzer, error) {
	// パッケージを読み込む
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports |
			packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir: workDir,
	}

	pkgs, err := packages.Load(cfg, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	return &WireAnalyzer{
		workDir:       workDir,
		searchPattern: searchPattern,
		analyzed:      make(map[string]*StructNode),
		pkgs:          pkgs,
	}, nil
}

// AnalyzeStruct は構造体を解析する（エントリーポイント）
// packagePath: 構造体のパッケージパス（空文字列の場合は全パッケージから検索）
// structName: 構造体名
func (wa *WireAnalyzer) AnalyzeStruct(structName string, packagePath core.PackagePath) (*StructNode, error) {
	// 対象の構造体型を検索（全パッケージから検索）
	structType, ok := core.FindStruct(structName, packagePath, wa.pkgs)
	if !ok {
		return nil, fmt.Errorf("struct %s not found in any package", structName)
	}

	return wa.analyzeStructType(structName, packagePath, structType)
}

// analyzeStructType は構造体型を解析する
func (wa *WireAnalyzer) analyzeStructType(structName string, packagePath core.PackagePath, structType *types.Struct) (*StructNode, error) {
	cacheKey := packagePath.String() + "." + structName

	// 既に解析済みの場合はキャッシュから返す
	if cached, ok := wa.analyzed[cacheKey]; ok {
		return cached, nil
	}

	result := &StructNode{
		StructName:    structName,
		PackagePath:   packagePath.String(),
		InitFunctions: make([]InitFunctionInfo, 0),
		Fields:        make([]FieldNode, 0),
	}

	// キャッシュに登録（無限ループ防止のため、フィールド解析前に登録）
	wa.analyzed[cacheKey] = result

	// 初期化関数を探す
	functions := core.FindFunctionsReturningStruct(structType, wa.pkgs)
	for _, fn := range functions {
		result.InitFunctions = append(result.InitFunctions, InitFunctionInfo{
			Name:        fn.Name,
			PackagePath: fn.PackagePath,
		})
	}

	// 各フィールドを解析
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldNode := wa.analyzeField(field)
		if fieldNode != nil {
			result.Fields = append(result.Fields, fieldNode)
		}
	}

	return result, nil
}

// analyzeField はフィールドを解析する
func (wa *WireAnalyzer) analyzeField(field *types.Var) FieldNode {
	fieldType := core.Deref(field.Type())
	fieldName := field.Name()
	packagePath := string(core.GetPackagePath(field.Type()))

	// インターフェース型の場合
	if interfaceName, interfaceType, ok := core.ConvertToNamedInterface(fieldType); ok {
		resolvedStruct, skipReason := wa.resolveInterface(interfaceType)
		return &InterfaceNode{
			FieldName:      fieldName,
			TypeName:       interfaceName,
			PackagePath:    packagePath,
			ResolvedStruct: resolvedStruct,
			Skipped:        skipReason != "",
			SkipReason:     skipReason,
		}
	}

	// 構造体型の場合
	if structName, structType, ok := core.ConvertToNamedStruct(fieldType); ok {
		if isBuiltinType(structName) {
			return nil
		}
		resolvedStruct, err := wa.analyzeStructType(packagePath, structName, structType)
		if err != nil {
			return nil
		}
		resolvedStruct.FieldName = fieldName
		return resolvedStruct
	}

	return nil
}

// resolveInterface はインターフェースから具体的な構造体を解決する
func (wa *WireAnalyzer) resolveInterface(interfaceType *types.Interface) (*StructNode, string) {
	functions := core.FindFunctionsReturningInterface(interfaceType, wa.pkgs)

	if len(functions) == 0 {
		return nil, "no implementing functions found"
	}

	if len(functions) > 1 {
		return nil, fmt.Sprintf("multiple implementing functions found (%d)", len(functions))
	}

	fn := functions[0]

	// 関数の返り値の型を取得
	implType, ok := core.GetFunctionReturnType(fn.Name, fn.PackagePath, wa.pkgs)
	if !ok {
		return nil, fmt.Sprintf("failed to find return type for %s", fn.Name)
	}

	implType = core.Deref(implType)
	structName, structType, ok := core.ConvertToNamedStruct(implType)
	if !ok {
		return nil, "implementing type is not a struct"
	}

	resolvedStruct, err := wa.analyzeStructType(fn.PackagePath, structName, structType)
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

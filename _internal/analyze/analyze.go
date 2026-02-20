package analyze

import (
	"fmt"
	"go/types"

	"github.com/rmocchy/cire/internal/core"
)

// wireAnalyzer はwire.goの解析を行う（internal/coreを使用）
type wireAnalyzer struct {
	analyzed      map[string]*StructNode // 解析済みの構造体をキャッシュ（無限ループ防止）
	functionCache core.FunctionCache     // 関数キャッシュ
	structCache   core.StructCache       // 構造体キャッシュ
	fieldAnalyzer FieldAnalyzer          // フィールド解析器
}

// NewWireAnalyzer は新しいStructAnalyzerを作成する
func NewWireAnalyzer(
	functionCache core.FunctionCache,
	structCache core.StructCache,
) (StructAnalyzer, error) {
	wa := &wireAnalyzer{
		analyzed:      make(map[string]*StructNode),
		functionCache: functionCache,
		structCache:   structCache,
	}

	// 循環依存を避けるため、FieldAnalyzerの初期化はWireAnalyzerの初期化後に行う
	wa.fieldAnalyzer = NewFieldAnalyzer(functionCache, wa)

	return wa, nil
}

// AnalyzeStruct は構造体を解析する（エントリーポイント）
// packagePath: 構造体のパッケージパス（空文字列の場合は全パッケージから検索）
// structName: 構造体名
func (wa *wireAnalyzer) AnalyzeStruct(structName string, packagePath core.PackagePath) (*StructNode, error) {
	// 対象の構造体型を検索
	structType, ok := wa.structCache.Get(structName, packagePath)
	if !ok {
		return nil, fmt.Errorf("struct %s not found", structName)
	}

	return wa.AnalyzeNamedStructType(structName, packagePath, structType)
}

// AnalyzeNamedStructType は構造体型を解析する
func (wa *wireAnalyzer) AnalyzeNamedStructType(structName string, packagePath core.PackagePath, structType *types.Struct) (*StructNode, error) {
	cacheKey := packagePath.String() + "." + structName

	// 既に解析済みの場合はキャッシュから返す
	if cached, ok := wa.analyzed[cacheKey]; ok {
		return cached, nil
	}

	// 初期化関数を探す
	fns := wa.functionCache.BulkGetByStructResult(structType)

	result := &StructNode{
		StructName:    structName,
		PackagePath:   packagePath.String(),
		InitFunctions: fns,
		Fields:        make([]FieldNode, 0),
		Dependencies:  make([]FieldNode, 0),
	}

	// キャッシュに登録（無限ループ防止のため、フィールド解析前に登録）
	wa.analyzed[cacheKey] = result

	// 初期化関数の引数を解析してDependenciesに追加
	deps := wa.analyzeDependencies(fns)
	result.Dependencies = append(result.Dependencies, deps...)

	// 各フィールドを解析
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if node := wa.fieldAnalyzer.AnalyzeTypeToFieldNode(field.Name(), field.Type()); node != nil {
			result.Fields = append(result.Fields, node)
		}
	}

	return result, nil
}

// analyzeDependencies は初期化関数の依存関係を解析する
func (wa *wireAnalyzer) analyzeDependencies(fns []*types.Func) []FieldNode {
	deps := make([]FieldNode, 0)

	for _, fn := range fns {
		sig := fn.Signature()
		params := sig.Params()
		if params == nil {
			continue
		}

		for i := 0; i < params.Len(); i++ {
			param := params.At(i)
			if node := wa.fieldAnalyzer.AnalyzeTypeToFieldNode(param.Name(), param.Type()); node != nil {
				deps = append(deps, node)
			}
		}
	}

	return deps
}

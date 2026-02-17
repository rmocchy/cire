package pipe

import (
	"fmt"
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
)

// FieldAnalyzer はフィールドの型解析を行うインターフェース
type FieldAnalyzer interface {
	// AnalyzeTypeToFieldNode は型を解析してFieldNodeを返す
	AnalyzeTypeToFieldNode(fieldName string, fieldType types.Type) FieldNode
}

// fieldAnalyzer はFieldAnalyzerの実装
type fieldAnalyzer struct {
	functionCache  core.FunctionCache
	structAnalyzer StructAnalyzer
}

// NewFieldAnalyzer は新しいFieldAnalyzerを作成する
func NewFieldAnalyzer(
	functionCache core.FunctionCache,
	structAnalyzer StructAnalyzer,
) FieldAnalyzer {
	return &fieldAnalyzer{
		functionCache:  functionCache,
		structAnalyzer: structAnalyzer,
	}
}

// AnalyzeTypeToFieldNode は型を解析してFieldNodeを返す
func (fa *fieldAnalyzer) AnalyzeTypeToFieldNode(fieldName string, fieldType types.Type) FieldNode {
	// ポインタを剥がす
	fieldType = core.Deref(fieldType)

	// ビルトイン型の場合はBuiltinNodeを返す
	if isBuiltinType(fieldType) {
		return &BuiltinNode{
			FieldName: fieldName,
			TypeName:  fieldType.String(),
		}
	}

	// Named型でない場合はnilを返す（スキップ）
	namedType, isNamed := fieldType.(*types.Named)
	if !isNamed {
		return nil
	}

	pkgPath := core.NewPackagePath(namedType.Obj().Pkg().Path())
	typeName := namedType.Obj().Name()
	underlying := namedType.Underlying()

	// 構造体型の場合
	if structType, isStruct := underlying.(*types.Struct); isStruct {
		return fa.analyzeStructField(fieldName, typeName, pkgPath, structType)
	}

	// インターフェース型の場合
	if interfaceType, isInterface := underlying.(*types.Interface); isInterface {
		return fa.analyzeInterfaceField(fieldName, typeName, pkgPath, interfaceType)
	}

	return nil
}

// analyzeStructField は構造体フィールドを解析してStructNodeを返す
func (fa *fieldAnalyzer) analyzeStructField(fieldName, typeName string, pkgPath core.PackagePath, structType *types.Struct) FieldNode {
	initFuncs := fa.functionCache.BulkGetByStructResult(structType)

	childNode, err := fa.structAnalyzer.AnalyzeNamedStructType(typeName, pkgPath, structType)
	if err != nil {
		// エラーの場合はスキップ情報を含むノードを返す
		return &StructNode{
			FieldName:     fieldName,
			StructName:    typeName,
			PackagePath:   pkgPath.String(),
			InitFunctions: initFuncs,
			Skipped:       true,
			SkipReason:    fmt.Sprintf("failed to analyze struct: %v", err),
		}
	}

	return &StructNode{
		FieldName:     fieldName,
		StructName:    typeName,
		PackagePath:   pkgPath.String(),
		InitFunctions: childNode.InitFunctions,
		Dependencies:  childNode.Dependencies,
		Fields:        childNode.Fields,
	}
}

// analyzeInterfaceField はインターフェースフィールドを解析してInterfaceNodeを返す
func (fa *fieldAnalyzer) analyzeInterfaceField(fieldName, typeName string, pkgPath core.PackagePath, interfaceType *types.Interface) FieldNode {
	initFns := fa.functionCache.BulkGetByInterfaceResult(interfaceType)

	// インターフェースの依存関係を解析
	deps := fa.analyzeInitFunctionParams(initFns)

	return &InterfaceNode{
		FieldName:     fieldName,
		TypeName:      typeName,
		PackagePath:   pkgPath.String(),
		InitFunctions: initFns,
		Dependencies:  deps,
	}
}

// analyzeInitFunctionParams は初期化関数の引数を解析してFieldNodeのリストを返す
func (fa *fieldAnalyzer) analyzeInitFunctionParams(fns []*types.Func) []FieldNode {
	deps := make([]FieldNode, 0)

	for _, fn := range fns {
		sig := fn.Signature()
		params := sig.Params()
		if params == nil {
			continue
		}

		for i := 0; i < params.Len(); i++ {
			param := params.At(i)
			if node := fa.AnalyzeTypeToFieldNode(param.Name(), param.Type()); node != nil {
				deps = append(deps, node)
			}
		}
	}

	return deps
}

// isBuiltinType はビルトイン型かどうかを判定する
func isBuiltinType(t types.Type) bool {
	_, ok := t.(*types.Basic)
	return ok
}

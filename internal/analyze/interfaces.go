package pipe

import (
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
)

// StructAnalyzer は構造体の解析を行うインターフェース
type StructAnalyzer interface {
	// AnalyzeStruct は構造体を解析する（エントリーポイント）
	AnalyzeStruct(structName string, packagePath core.PackagePath) (*StructNode, error)
	// AnalyzeNamedStructType は名前付き構造体型を解析する
	AnalyzeNamedStructType(structName string, packagePath core.PackagePath, structType *types.Struct) (*StructNode, error)
}

// isBuiltinType はビルトイン型かどうかを判定する
func isBuiltinType(t types.Type) bool {
	_, ok := t.(*types.Basic)
	return ok
}

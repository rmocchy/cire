package pipe

import (
	"fmt"
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
)

// WireAnalyzer はwire.goの解析を行う（internal/coreを使用）
type WireAnalyzer struct {
	analyzed      map[string]*StructNode // 解析済みの構造体をキャッシュ（無限ループ防止）
	functionCache core.FunctionCache     // 関数キャッシュ
	structCache   core.StructCache       // 構造体キャッシュ
}

// NewWireAnalyzer は新しいWireAnalyzerを作成する
func NewWireAnalyzer(
	functionCache core.FunctionCache,
	structCache core.StructCache,
) (*WireAnalyzer, error) {
	return &WireAnalyzer{
		analyzed:      make(map[string]*StructNode),
		functionCache: functionCache,
		structCache:   structCache,
	}, nil
}

// AnalyzeStruct は構造体を解析する（エントリーポイント）
// packagePath: 構造体のパッケージパス（空文字列の場合は全パッケージから検索）
// structName: 構造体名
func (wa *WireAnalyzer) AnalyzeStruct(structName string, packagePath core.PackagePath) (*StructNode, error) {
	// 対象の構造体型を検索
	structType, ok := wa.structCache.Get(structName, packagePath)
	if !ok {
		return nil, fmt.Errorf("struct %s not found", structName)
	}

	return wa.analyzeNamedStructType(structName, packagePath, structType)
}

// analyzeStructType は構造体型を解析する
func (wa *WireAnalyzer) analyzeNamedStructType(structName string, packagePath core.PackagePath, structType *types.Struct) (*StructNode, error) {
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
	}

	// キャッシュに登録（無限ループ防止のため、フィールド解析前に登録）
	wa.analyzed[cacheKey] = result

	fieldNodes := make([]FieldNode, 0, structType.NumFields())
	// 各フィールドを解析
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldType := core.Deref(field.Type())

		// ビルトイン型の場合はBuiltinNodeを作成して追加
		if isBuiltinType(fieldType) {
			fieldNodes = append(fieldNodes, &BuiltinNode{
				FieldName: field.Name(),
				TypeName:  fieldType.String(),
			})
			continue
		}

		// 無名宣言の場合は無視
		namedField, isNamed := fieldType.(*types.Named)
		if !isNamed {
			continue
		}

		defPkgPath := core.NewPackagePath(namedField.Obj().Pkg().Path())
		structField, isStruct := namedField.Underlying().(*types.Struct)
		interfaceField, isInterface := namedField.Underlying().(*types.Interface)

		// 構造体型の場合は再帰的に解析
		if isStruct {
			initFuncs := wa.functionCache.BulkGetByStructResult(structField)
			if len(initFuncs) > 0 {
				result.InitFunctions = append(result.InitFunctions, initFuncs...)
			}

			childNode, err := wa.analyzeNamedStructType(namedField.Obj().Name(), defPkgPath, structField)
			if err != nil {
				// エラーがあっても他のフィールドの解析を続ける
				fieldNodes = append(fieldNodes, &StructNode{
					FieldName:     field.Name(),
					StructName:    namedField.Obj().Name(),
					PackagePath:   defPkgPath.String(),
					InitFunctions: initFuncs,
					Skipped:       true,
					SkipReason:    fmt.Sprintf("failed to analyze struct field: %v", err),
				})
				continue
			}
			fieldNodes = append(fieldNodes, &StructNode{
				FieldName:     field.Name(),
				StructName:    namedField.Obj().Name(),
				PackagePath:   defPkgPath.String(),
				InitFunctions: initFuncs,
				Fields:        childNode.Fields,
			})
			continue
		}

		// インターフェース型の場合はInterfaceNodeを作成して追加
		if isInterface {
			initFns := wa.functionCache.BulkGetByInterfaceResult(interfaceField)

			fieldNodes = append(fieldNodes, &InterfaceNode{
				FieldName:     field.Name(),
				TypeName:      namedField.Obj().Name(),
				PackagePath:   defPkgPath.String(),
				InitFunctions: initFns,
			})
			continue
		}
	}

	return &StructNode{
		StructName:    structName,
		PackagePath:   packagePath.String(),
		InitFunctions: wa.functionCache.BulkGetByStructResult(structType),
		Fields:        fieldNodes,
	}, nil
}

// isBuiltinType はビルトイン型かどうかを判定する
func isBuiltinType(typeName types.Type) bool {
	_, ok := typeName.(*types.Basic)
	return ok
}

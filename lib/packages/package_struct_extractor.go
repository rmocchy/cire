package packages

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// ExtractStructFields は作業ディレクトリを指定してpackagePathと構造体名から構造体のフィールド情報を取得する
// workDir: パッケージ解決の基準となる作業ディレクトリ（空文字列の場合はカレントディレクトリ）
// packagePath: パッケージパス（モジュールパスまたは相対パス）
// structName: 取得する構造体の名前
func ExtractStructFields(workDir, packagePath, structName string) (*StructFieldsInfo, error) {
	// パッケージをロード
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  workDir,
	}

	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found for path: %s", packagePath)
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("package has errors: %v", pkg.Errors)
	}

	// 構造体の型を検索
	scope := pkg.Types.Scope()
	obj := scope.Lookup(structName)
	if obj == nil {
		return nil, fmt.Errorf("struct %s not found in package %s", structName, packagePath)
	}

	typeName, ok := obj.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("%s is not a type", structName)
	}

	named, ok := typeName.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("%s is not a named type", structName)
	}

	// Underlying型を取得してaliasを展開
	under := types.Unalias(named.Underlying())
	structType, ok := under.(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("%s is not a struct type", structName)
	}

	// フィールド情報を抽出
	fields := extractFields(structType)

	return &StructFieldsInfo{
		StructName: structName,
		Fields:     fields,
	}, nil
}

// extractFields は構造体のフィールド情報を抽出する
func extractFields(structType *types.Struct) []FieldInfo {
	var fields []FieldInfo

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldType := field.Type()

		fieldInfo := parseFieldType(field.Name(), fieldType)
		fields = append(fields, fieldInfo)
	}

	return fields
}

// parseFieldType はフィールドの型情報を解析してFieldInfoを作成する
func parseFieldType(fieldName string, fieldType types.Type) FieldInfo {
	info := FieldInfo{
		Name: fieldName,
	}

	// ポインタ型の場合は剥がす
	originalType := fieldType
	fieldType = derefType(fieldType)
	info.IsPointer = originalType != fieldType

	// Aliasを展開
	fieldType = types.Unalias(fieldType)

	// Named型の場合、パッケージパスと型名を取得
	if named, ok := fieldType.(*types.Named); ok {
		obj := named.Obj()
		info.TypeName = obj.Name()

		// パッケージ情報を取得
		if pkg := obj.Pkg(); pkg != nil {
			info.PackagePath = pkg.Path()
		}

		// Named型の基底型がインターフェースかどうかチェック
		underlying := types.Unalias(named.Underlying())
		if _, ok := underlying.(*types.Interface); ok {
			info.IsInterface = true
		}
	} else {
		// Named型でない場合（基本型など）は型の文字列表現を使用
		info.TypeName = fieldType.String()

		// 直接インターフェース型かどうかチェック
		if _, ok := fieldType.(*types.Interface); ok {
			info.IsInterface = true
		}
	}

	return info
}

// derefType はポインタ型を再帰的に剥がす
func derefType(t types.Type) types.Type {
	for {
		ptr, ok := t.(*types.Pointer)
		if !ok {
			return t
		}
		t = ptr.Elem()
	}
}

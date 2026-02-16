package core

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type StructInfo struct {
	StructName string
	Fields     []Field
}

type Field struct {
	Name string
	Type types.Type
}

// NewStructInfo は*types.Structから*StructInfoを生成する
func NewStructInfo(structType *types.Struct) (*StructInfo, error) {
	if structType == nil {
		return nil, fmt.Errorf("invalid *types.Struct argument: nil")
	}

	fields := make([]Field, 0, structType.NumFields())
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fields = append(fields, Field{
			Name: field.Name(),
			Type: field.Type(),
		})
	}

	return &StructInfo{
		StructName: structType.String(),
		Fields:     fields,
	}, nil
}

func FindStruct(name string, packagePath PackagePath, pkgs []*packages.Package) (*types.Struct, bool) {
	for _, pkg := range pkgs {
		if pkg.PkgPath != packagePath.String() {
			continue
		}

		obj := pkg.Types.Scope().Lookup(name)
		if obj == nil {
			continue
		}

		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		named, ok := typeName.Type().(*types.Named)
		if !ok {
			continue
		}

		structType, ok := named.Underlying().(*types.Struct)
		if !ok {
			continue
		}

		return structType, true
	}

	return nil, false
}

// matchesStruct は型が指定された構造体と一致するかチェック
func matchesStruct(t types.Type, targetStruct *types.Struct) bool {
	// ポインタを剥がす
	t = derefType(t)

	// エイリアスを解決
	t = types.Unalias(t)

	// Named型かチェック
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	// 基底型が構造体であることをチェック（インターフェースを除外）
	underlying := named.Underlying()
	structType, ok := underlying.(*types.Struct)
	if !ok {
		return false
	}

	// 構造体が同一かをチェック
	return types.Identical(structType, targetStruct)
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

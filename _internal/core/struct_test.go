package core

import (
	"go/types"
	"testing"
)

func TestNewStructInfo(t *testing.T) {
	t.Run("正常なStruct型の場合", func(t *testing.T) {
		// テスト用のStruct型を作成
		pkg := types.NewPackage("test/pkg", "pkg")
		fields := []*types.Var{
			types.NewField(0, pkg, "Name", types.Typ[types.String], false),
			types.NewField(0, pkg, "Age", types.Typ[types.Int], false),
		}
		structType := types.NewStruct(fields, nil)

		info, err := NewStructInfo(structType)
		if err != nil {
			t.Fatalf("エラーが発生しました: %v", err)
		}

		if info == nil {
			t.Fatal("StructInfoがnilです")
		}

		if len(info.Fields) != 2 {
			t.Errorf("フィールド数が期待値と異なります。期待: 2, 実際: %d", len(info.Fields))
		}

		if info.Fields[0].Name != "Name" {
			t.Errorf("フィールド名が期待値と異なります。期待: Name, 実際: %s", info.Fields[0].Name)
		}

		if info.Fields[1].Name != "Age" {
			t.Errorf("フィールド名が期待値と異なります。期待: Age, 実際: %s", info.Fields[1].Name)
		}
	})

	t.Run("フィールドがないStruct型の場合", func(t *testing.T) {
		structType := types.NewStruct(nil, nil)

		info, err := NewStructInfo(structType)
		if err != nil {
			t.Fatalf("エラーが発生しました: %v", err)
		}

		if info == nil {
			t.Fatal("StructInfoがnilです")
		}

		if len(info.Fields) != 0 {
			t.Errorf("フィールド数が期待値と異なります。期待: 0, 実際: %d", len(info.Fields))
		}
	})

	t.Run("nilを渡した場合", func(t *testing.T) {
		info, err := NewStructInfo(nil)
		if err == nil {
			t.Error("エラーが発生しませんでした")
		}

		if info != nil {
			t.Error("StructInfoがnilではありません")
		}
	})
}

func TestDerefType(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	structType := types.NewStruct(nil, nil)
	namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)

	t.Run("ポインタ型を剥がす", func(t *testing.T) {
		ptrType := types.NewPointer(namedStruct)
		result := derefType(ptrType)

		if result != namedStruct {
			t.Error("ポインタが剥がされていません")
		}
	})

	t.Run("多重ポインタ型を剥がす", func(t *testing.T) {
		ptrType := types.NewPointer(types.NewPointer(types.NewPointer(namedStruct)))
		result := derefType(ptrType)

		if result != namedStruct {
			t.Error("多重ポインタが剥がされていません")
		}
	})

	t.Run("非ポインタ型はそのまま返す", func(t *testing.T) {
		result := derefType(namedStruct)

		if result != namedStruct {
			t.Error("元の型と異なります")
		}
	})
}

func TestMatchesStruct(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")

	// フィールドを持つ構造体1
	fields1 := []*types.Var{
		types.NewField(0, pkg, "Field1", types.Typ[types.String], false),
	}
	structType1 := types.NewStruct(fields1, nil)
	namedStruct1 := types.NewNamed(types.NewTypeName(0, pkg, "Struct1", nil), structType1, nil)

	// フィールドを持つ構造体2
	fields2 := []*types.Var{
		types.NewField(0, pkg, "Field2", types.Typ[types.Int], false),
	}
	structType2 := types.NewStruct(fields2, nil)

	t.Run("同じ構造体の場合", func(t *testing.T) {
		if !matchesStruct(namedStruct1, structType1) {
			t.Error("同じ構造体が一致しませんでした")
		}
	})

	t.Run("ポインタ型の構造体の場合", func(t *testing.T) {
		ptrType := types.NewPointer(namedStruct1)
		if !matchesStruct(ptrType, structType1) {
			t.Error("ポインタ型の構造体が一致しませんでした")
		}
	})

	t.Run("異なる構造体の場合", func(t *testing.T) {
		if matchesStruct(namedStruct1, structType2) {
			t.Error("異なる構造体が一致してしまいました")
		}
	})

	t.Run("インターフェース型の場合", func(t *testing.T) {
		ifaceType := types.NewInterfaceType(nil, nil)
		namedIface := types.NewNamed(types.NewTypeName(0, pkg, "MyInterface", nil), ifaceType, nil)

		if matchesStruct(namedIface, structType1) {
			t.Error("インターフェース型が構造体として一致してしまいました")
		}
	})

	t.Run("基本型の場合", func(t *testing.T) {
		if matchesStruct(types.Typ[types.Int], structType1) {
			t.Error("基本型が構造体として一致してしまいました")
		}
	})
}

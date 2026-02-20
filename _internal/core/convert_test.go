package core

import (
	"go/types"
	"testing"
)

func TestConvertToPointer(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	structType := types.NewStruct(nil, nil)
	namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)
	ptrType := types.NewPointer(namedStruct)

	t.Run("ポインタ型の場合", func(t *testing.T) {
		ptr, ok := ConvertToPointer(ptrType)
		if !ok {
			t.Fatal("ポインタ型の取得に失敗しました")
		}
		if ptr.Elem() != namedStruct {
			t.Error("期待した要素型と異なります")
		}
	})

	t.Run("ポインタ型でない場合", func(t *testing.T) {
		_, ok := ConvertToPointer(namedStruct)
		if ok {
			t.Error("非ポインタ型がポインタとして取得できてしまいました")
		}
	})
}

func TestConvertToNamed(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	structType := types.NewStruct(nil, nil)
	namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)

	t.Run("Named型の場合", func(t *testing.T) {
		named, ok := ConvertToNamed(namedStruct)
		if !ok {
			t.Fatal("Named型の取得に失敗しました")
		}
		if named.Obj().Name() != "MyStruct" {
			t.Error("期待した型名と異なります")
		}
	})

	t.Run("Named型でない場合", func(t *testing.T) {
		_, ok := ConvertToNamed(types.Typ[types.Int])
		if ok {
			t.Error("基本型がNamed型として取得できてしまいました")
		}
	})
}

func TestConvertToNamedStruct(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	structType := types.NewStruct(nil, nil)
	namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)

	t.Run("Named型のStructの場合", func(t *testing.T) {
		name, st, ok := ConvertToNamedStruct(namedStruct)
		if !ok {
			t.Fatal("Named型のStruct取得に失敗しました")
		}
		if name != "MyStruct" {
			t.Errorf("expected name 'MyStruct', got '%s'", name)
		}
		if st != structType {
			t.Error("期待した構造体と異なります")
		}
	})

	t.Run("Named型でない場合", func(t *testing.T) {
		_, _, ok := ConvertToNamedStruct(structType)
		if ok {
			t.Error("Named型でないのに取得できてしまいました")
		}
	})
}

func TestConvertToNamedInterface(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	interfaceType := types.NewInterfaceType(nil, nil)
	namedInterface := types.NewNamed(types.NewTypeName(0, pkg, "MyInterface", nil), interfaceType, nil)

	t.Run("Named型のInterfaceの場合", func(t *testing.T) {
		name, iface, ok := ConvertToNamedInterface(namedInterface)
		if !ok {
			t.Fatal("Named型のInterface取得に失敗しました")
		}
		if name != "MyInterface" {
			t.Errorf("expected name 'MyInterface', got '%s'", name)
		}
		if iface != interfaceType {
			t.Error("期待したインターフェースと異なります")
		}
	})

	t.Run("Named型でない場合", func(t *testing.T) {
		_, _, ok := ConvertToNamedInterface(interfaceType)
		if ok {
			t.Error("Named型でないのに取得できてしまいました")
		}
	})
}

func TestDeref(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	structType := types.NewStruct(nil, nil)
	namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)
	ptrType := types.NewPointer(namedStruct)

	t.Run("ポインタ型の場合", func(t *testing.T) {
		elem := Deref(ptrType)
		if elem != namedStruct {
			t.Error("期待した要素型と異なります")
		}
	})

	t.Run("ポインタ型でない場合", func(t *testing.T) {
		elem := Deref(namedStruct)
		if elem != namedStruct {
			t.Error("元の型がそのまま返されませんでした")
		}
	})
}

func TestUsagePattern(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")
	structType := types.NewStruct(nil, nil)
	namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)
	ptrType := types.NewPointer(namedStruct)

	t.Run("実践的な使用例: ポインタを剥がしてNamed Structを取得", func(t *testing.T) {
		// Goの型アサーションのように使える
		elem := Deref(ptrType)
		name, st, ok := ConvertToNamedStruct(elem)
		if !ok {
			t.Fatal("Named Structの取得に失敗しました")
		}
		if name != "MyStruct" {
			t.Errorf("expected name 'MyStruct', got '%s'", name)
		}
		if st != structType {
			t.Error("期待した構造体と異なります")
		}
	})

	t.Run("実践的な使用例: 型によって処理を分岐", func(t *testing.T) {
		elem := Deref(ptrType)

		// structかinterfaceかを判定
		if name, st, ok := ConvertToNamedStruct(elem); ok {
			// 構造体の場合の処理
			if name != "MyStruct" || st == nil {
				t.Error("構造体の処理が正しくありません")
			}
		} else if _, _, ok := ConvertToNamedInterface(elem); ok {
			// インターフェースの場合の処理
			t.Error("構造体がインターフェースと判定されました")
		}
	})
}

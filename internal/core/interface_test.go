package core

import (
	"go/types"
	"testing"
)

func TestMatchesInterface(t *testing.T) {
	pkg := types.NewPackage("test/pkg", "pkg")

	// インターフェース1を作成
	methods1 := []*types.Func{
		types.NewFunc(0, pkg, "Method1", types.NewSignatureType(nil, nil, nil, nil, nil, false)),
	}
	ifaceType1 := types.NewInterfaceType(methods1, nil)
	ifaceType1.Complete()
	namedIface1 := types.NewNamed(types.NewTypeName(0, pkg, "Interface1", nil), ifaceType1, nil)

	// インターフェース2を作成
	methods2 := []*types.Func{
		types.NewFunc(0, pkg, "Method2", types.NewSignatureType(nil, nil, nil, nil, nil, false)),
	}
	ifaceType2 := types.NewInterfaceType(methods2, nil)
	ifaceType2.Complete()

	t.Run("同じインターフェースの場合", func(t *testing.T) {
		if !matchesInterface(namedIface1, ifaceType1) {
			t.Error("同じインターフェースが一致しませんでした")
		}
	})

	t.Run("ポインタ型のインターフェースの場合", func(t *testing.T) {
		ptrType := types.NewPointer(namedIface1)
		if !matchesInterface(ptrType, ifaceType1) {
			t.Error("ポインタ型のインターフェースが一致しませんでした")
		}
	})

	t.Run("異なるインターフェースの場合", func(t *testing.T) {
		if matchesInterface(namedIface1, ifaceType2) {
			t.Error("異なるインターフェースが一致してしまいました")
		}
	})

	t.Run("構造体型の場合", func(t *testing.T) {
		structType := types.NewStruct(nil, nil)
		namedStruct := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)

		if matchesInterface(namedStruct, ifaceType1) {
			t.Error("構造体型がインターフェースとして一致してしまいました")
		}
	})

	t.Run("基本型の場合", func(t *testing.T) {
		if matchesInterface(types.Typ[types.Int], ifaceType1) {
			t.Error("基本型がインターフェースとして一致してしまいました")
		}
	})

	t.Run("Named型でない場合", func(t *testing.T) {
		// Named型でないインターフェースを直接渡す
		if matchesInterface(ifaceType1, ifaceType1) {
			t.Error("Named型でないインターフェースが一致してしまいました")
		}
	})
}

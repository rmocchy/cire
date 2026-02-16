package core

import (
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestFindFunctionsReturningStruct(t *testing.T) {
	// テスト用のパッケージを作成
	pkg := types.NewPackage("test/pkg", "pkg")
	scope := pkg.Scope()

	// テスト用の構造体型を作成
	targetFields := []*types.Var{
		types.NewField(0, pkg, "TargetField", types.Typ[types.String], false),
	}
	targetStructType := types.NewStruct(targetFields, nil)
	targetNamed := types.NewNamed(types.NewTypeName(0, pkg, "TargetStruct", nil), targetStructType, nil)

	// 別の構造体型を作成
	otherFields := []*types.Var{
		types.NewField(0, pkg, "OtherField", types.Typ[types.Int], false),
	}
	otherStructType := types.NewStruct(otherFields, nil)
	otherNamed := types.NewNamed(types.NewTypeName(0, pkg, "OtherStruct", nil), otherStructType, nil)

	// TargetStructを返す関数を作成
	funcReturnsTarget := types.NewFunc(0, pkg, "GetTargetStruct", types.NewSignatureType(
		nil, nil, nil,
		nil,
		types.NewTuple(types.NewVar(0, pkg, "", targetNamed)),
		false,
	))
	scope.Insert(funcReturnsTarget)

	// TargetStructのポインタを返す関数を作成
	funcReturnsPtrTarget := types.NewFunc(0, pkg, "GetTargetStructPtr", types.NewSignatureType(
		nil, nil, nil,
		nil,
		types.NewTuple(types.NewVar(0, pkg, "", types.NewPointer(targetNamed))),
		false,
	))
	scope.Insert(funcReturnsPtrTarget)

	// OtherStructを返す関数を作成
	funcReturnsOther := types.NewFunc(0, pkg, "GetOtherStruct", types.NewSignatureType(
		nil, nil, nil,
		nil,
		types.NewTuple(types.NewVar(0, pkg, "", otherNamed)),
		false,
	))
	scope.Insert(funcReturnsOther)

	// 返り値がない関数を作成
	funcNoReturn := types.NewFunc(0, pkg, "NoReturn", types.NewSignatureType(
		nil, nil, nil,
		nil,
		nil,
		false,
	))
	scope.Insert(funcNoReturn)

	// packages.Packageを作成
	testPkg := &packages.Package{
		PkgPath: "test/pkg",
		Types:   pkg,
	}
	pkgs := []*packages.Package{testPkg}

	t.Run("指定された構造体を返す関数を見つける", func(t *testing.T) {
		functions := FindFunctionsReturningStruct(targetStructType, pkgs)

		if len(functions) != 2 {
			t.Fatalf("見つかった関数の数が期待値と異なります。期待: 2, 実際: %d", len(functions))
		}

		foundNames := make(map[string]bool)
		for _, fn := range functions {
			foundNames[fn.Name] = true
		}

		if !foundNames["GetTargetStruct"] {
			t.Error("GetTargetStructが見つかりませんでした")
		}
		if !foundNames["GetTargetStructPtr"] {
			t.Error("GetTargetStructPtrが見つかりませんでした")
		}
	})

	t.Run("nilを渡した場合", func(t *testing.T) {
		functions := FindFunctionsReturningStruct(nil, pkgs)

		if functions != nil {
			t.Error("nilを渡した場合はnilが返されるべきです")
		}
	})

	t.Run("空のパッケージリストの場合", func(t *testing.T) {
		functions := FindFunctionsReturningStruct(targetStructType, []*packages.Package{})

		if len(functions) != 0 {
			t.Errorf("空のパッケージリストの場合は空のスライスが返されるべきです。実際: %d", len(functions))
		}
	})
}

func TestFindFunctionsReturningInterface(t *testing.T) {
	// テスト用のパッケージを作成
	pkg := types.NewPackage("test/pkg", "pkg")
	scope := pkg.Scope()

	// テスト用のインターフェース型を作成
	methods := []*types.Func{
		types.NewFunc(0, pkg, "DoSomething", types.NewSignatureType(nil, nil, nil, nil, nil, false)),
	}
	targetIfaceType := types.NewInterfaceType(methods, nil)
	targetIfaceType.Complete()
	targetNamed := types.NewNamed(types.NewTypeName(0, pkg, "TargetInterface", nil), targetIfaceType, nil)

	// 別のインターフェース型を作成
	otherMethods := []*types.Func{
		types.NewFunc(0, pkg, "DoOther", types.NewSignatureType(nil, nil, nil, nil, nil, false)),
	}
	otherIfaceType := types.NewInterfaceType(otherMethods, nil)
	otherIfaceType.Complete()
	otherNamed := types.NewNamed(types.NewTypeName(0, pkg, "OtherInterface", nil), otherIfaceType, nil)

	// TargetInterfaceを返す関数を作成
	funcReturnsTarget := types.NewFunc(0, pkg, "GetTargetInterface", types.NewSignatureType(
		nil, nil, nil,
		nil,
		types.NewTuple(types.NewVar(0, pkg, "", targetNamed)),
		false,
	))
	scope.Insert(funcReturnsTarget)

	// TargetInterfaceのポインタを返す関数を作成
	funcReturnsPtrTarget := types.NewFunc(0, pkg, "GetTargetInterfacePtr", types.NewSignatureType(
		nil, nil, nil,
		nil,
		types.NewTuple(types.NewVar(0, pkg, "", types.NewPointer(targetNamed))),
		false,
	))
	scope.Insert(funcReturnsPtrTarget)

	// OtherInterfaceを返す関数を作成
	funcReturnsOther := types.NewFunc(0, pkg, "GetOtherInterface", types.NewSignatureType(
		nil, nil, nil,
		nil,
		types.NewTuple(types.NewVar(0, pkg, "", otherNamed)),
		false,
	))
	scope.Insert(funcReturnsOther)

	// packages.Packageを作成
	testPkg := &packages.Package{
		PkgPath: "test/pkg",
		Types:   pkg,
	}
	pkgs := []*packages.Package{testPkg}

	t.Run("指定されたインターフェースを返す関数を見つける", func(t *testing.T) {
		functions := FindFunctionsReturningInterface(targetIfaceType, pkgs)

		if len(functions) != 2 {
			t.Fatalf("見つかった関数の数が期待値と異なります。期待: 2, 実際: %d", len(functions))
		}

		foundNames := make(map[string]bool)
		for _, fn := range functions {
			foundNames[fn.Name] = true
		}

		if !foundNames["GetTargetInterface"] {
			t.Error("GetTargetInterfaceが見つかりませんでした")
		}
		if !foundNames["GetTargetInterfacePtr"] {
			t.Error("GetTargetInterfacePtrが見つかりませんでした")
		}
	})

	t.Run("nilを渡した場合", func(t *testing.T) {
		functions := FindFunctionsReturningInterface(nil, pkgs)

		if functions != nil {
			t.Error("nilを渡した場合はnilが返されるべきです")
		}
	})

	t.Run("空のパッケージリストの場合", func(t *testing.T) {
		functions := FindFunctionsReturningInterface(targetIfaceType, []*packages.Package{})

		if len(functions) != 0 {
			t.Errorf("空のパッケージリストの場合は空のスライスが返されるべきです。実際: %d", len(functions))
		}
	})
}

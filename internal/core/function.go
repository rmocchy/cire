package core

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type FunctionInfo struct {
	Name        string
	PackagePath string
}

// FindFunctionsReturningStruct は指定された構造体を返り値に持つ関数を探す
func FindFunctionsReturningStruct(targetStruct *types.Struct, pkgs []*packages.Package) []FunctionInfo {
	if targetStruct == nil {
		return nil
	}

	functions := make([]FunctionInfo, 0)

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)

			// 関数かどうかをチェック
			fn, ok := obj.(*types.Func)
			if !ok {
				continue
			}

			// 関数のシグネチャを取得
			sig, ok := fn.Type().(*types.Signature)
			if !ok {
				continue
			}

			// 返り値をチェック
			results := sig.Results()
			if results == nil {
				continue
			}

			// 各返り値が指定された構造体かどうかをチェック
			for i := 0; i < results.Len(); i++ {
				result := results.At(i)
				if matchesStruct(result.Type(), targetStruct) {
					functions = append(functions, FunctionInfo{
						Name:        fn.Name(),
						PackagePath: pkg.PkgPath,
					})
					break // 同じ関数を複数回追加しないように
				}
			}
		}
	}

	return functions
}

// FindFunctionsReturningInterface は指定されたinterfaceを返り値に持つ関数を探す
func FindFunctionsReturningInterface(targetInterface *types.Interface, pkgs []*packages.Package) []FunctionInfo {
	if targetInterface == nil {
		return nil
	}

	functions := make([]FunctionInfo, 0)

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)

			// 関数かどうかをチェック
			fn, ok := obj.(*types.Func)
			if !ok {
				continue
			}

			// 関数のシグネチャを取得
			sig, ok := fn.Type().(*types.Signature)
			if !ok {
				continue
			}

			// 返り値をチェック
			results := sig.Results()
			if results == nil {
				continue
			}

			// 各返り値が指定されたinterfaceかどうかをチェック
			for i := 0; i < results.Len(); i++ {
				result := results.At(i)
				if matchesInterface(result.Type(), targetInterface) {
					functions = append(functions, FunctionInfo{
						Name:        fn.Name(),
						PackagePath: pkg.PkgPath,
					})
					break // 同じ関数を複数回追加しないように
				}
			}
		}
	}

	return functions
}

// GetFunctionReturnType は指定された関数の返り値の型を取得する
func GetFunctionReturnType(funcName, packagePath string, pkgs []*packages.Package) (types.Type, bool) {
	for _, pkg := range pkgs {
		if pkg.PkgPath != packagePath {
			continue
		}

		obj := pkg.Types.Scope().Lookup(funcName)
		if obj == nil {
			continue
		}

		fn, ok := obj.(*types.Func)
		if !ok {
			continue
		}

		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Results() == nil || sig.Results().Len() == 0 {
			continue
		}

		return sig.Results().At(0).Type(), true
	}

	return nil, false
}

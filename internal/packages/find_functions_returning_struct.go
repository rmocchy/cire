package packages

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

// FindFunctionsReturningStruct は指定された構造体を返り値に持つ関数を探す
// structName: 構造体名
// structPkgPath: 構造体が定義されているパッケージパス
// pkgs: 検索対象のパッケージ群
// 注意: インターフェースは対象外です
func FindFunctionsReturningStruct(structName, structPkgPath string, pkgs []*packages.Package) []FunctionInfo {
	var functions []FunctionInfo

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
				if matchesStructType(result.Type(), structName, structPkgPath) {
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

// matchesStructType は型が指定された構造体と一致するかチェック
// インターフェースは対象外
func matchesStructType(t types.Type, structName, structPkgPath string) bool {
	// ポインタを剥がす
	t = derefType(t)

	// エイリアスを解決
	t = types.Unalias(t)

	// Named型かチェック
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	// 型名をチェック
	if named.Obj().Name() != structName {
		return false
	}

	// 基底型が構造体であることをチェック（インターフェースを除外）
	underlying := named.Underlying()
	if _, ok := underlying.(*types.Struct); !ok {
		return false
	}

	// パッケージパスをチェック
	if pkg := named.Obj().Pkg(); pkg != nil {
		return pkg.Path() == structPkgPath
	}

	return false
}

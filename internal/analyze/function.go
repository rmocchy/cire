package analyze

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type FunctionCache interface {
	BulkGet(returnType *types.Named) []*types.Func
}

type functionCache struct {
	fns map[string]*types.Func
}

func NewFunctionCache(pkgs []*packages.Package) FunctionCache {
	// ここでは単純に全ての関数をキャッシュする例を示す
	// 実際には必要な関数のみをキャッシュするように最適化することも可能
	fns := make(map[string]*types.Func)

	for _, pkg := range pkgs {
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if fn, ok := obj.(*types.Func); ok {
				fns[pkg.PkgPath+"."+fn.Name()] = fn
			}
		}
	}

	return &functionCache{fns: fns}
}

func (fc *functionCache) BulkGet(returnType *types.Named) []*types.Func {
	// キャッシュから指定された返り値の型を持つ関数を取得
	result := make([]*types.Func, 0)
	for _, fn := range fc.fns {
		ret := fn.Signature().Results()
		for i := 0; i < ret.Len(); i++ {
			paramType := Deref(ret.At(i).Type())

			if types.Identical(paramType, returnType) {
				result = append(result, fn)
			}
		}
	}
	return result
}

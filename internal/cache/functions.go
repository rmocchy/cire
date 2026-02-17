package cache

import (
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
	"golang.org/x/tools/go/packages"
)

type functionCache struct {
	fns map[string]*types.Func
}

func NewFunctionCache(pkgs []*packages.Package) core.FunctionCache {
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

func (fc *functionCache) BulkGetByStructResult(structure *types.Struct) []*types.Func {
	// キャッシュから構造体を返り値に持つ関数を取得
	result := make([]*types.Func, 0)
	for _, fn := range fc.fns {
		ret := fn.Signature().Results()
		for i := 0; i < ret.Len(); i++ {
			paramType := ret.At(i).Type()
			if types.Identical(paramType, structure) {
				result = append(result, fn)
				break
			}
		}
	}
	return result
}

func (fc *functionCache) BulkGetByInterfaceResult(interfaceType *types.Interface) []*types.Func {
	// キャッシュからinterfaceを返り値に持つ関数を取得
	result := make([]*types.Func, 0)
	for _, fn := range fc.fns {
		ret := fn.Signature().Results()
		for i := 0; i < ret.Len(); i++ {
			paramType := ret.At(i).Type()
			if types.Implements(paramType, interfaceType) {
				result = append(result, fn)
				break
			}
		}
	}
	return result
}

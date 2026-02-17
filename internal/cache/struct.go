package cache

import (
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
)

type structCache struct {
	Structs map[string]*types.Struct
}

func NewStructCache(pkgs []*types.Package) core.StructCache {
	structs := make(map[string]*types.Struct)

	for _, pkg := range pkgs {
		scope := pkg.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			// 型名から構造体を取得
			if typeName, ok := obj.(*types.TypeName); ok {
				if st, ok := typeName.Type().Underlying().(*types.Struct); ok {
					structs[pkg.Path()+"."+name] = st
				}
			}
		}
	}

	return &structCache{Structs: structs}
}

func (sc *structCache) Get(name string, pkgPath core.PackagePath) (*types.Struct, bool) {
	key := pkgPath.String() + "." + name
	st, ok := sc.Structs[key]
	return st, ok
}

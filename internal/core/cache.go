package core

import (
	"go/types"
)

type FunctionCache interface {
	BulkGetByStructResult(structure *types.Struct) []*types.Func
	BulkGetByInterfaceResult(iface *types.Interface) []*types.Func
}

type StructCache interface {
	Get(name string, pkgPath PackagePath) (*types.Struct, bool)
}

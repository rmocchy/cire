package core

import (
	"go/types"
)

type PackagePath string

// GetPackagePath は型からパッケージパスを取得する
func GetPackagePath(t types.Type) PackagePath {
	// ポインタを剥がす
	t = Deref(t)

	// Named型からパッケージパスを取得
	if named, ok := ConvertToNamed(t); ok {
		if pkg := named.Obj().Pkg(); pkg != nil {
			return PackagePath(pkg.Path())
		}
	}

	return ""
}

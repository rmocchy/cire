package core

type PackagePath struct {
	path string
}

func NewPackagePath(path string) PackagePath {
	return PackagePath{path: path}
}

func (p PackagePath) String() string {
	return p.path
}

// GetPackagePath は型からパッケージパスを取得する
// func GetPackagePath(t types.Type) PackagePath {
// 	// ポインタを剥がす
// 	t = Deref(t)

// 	// Named型からパッケージパスを取得
// 	if named, ok := ConvertToNamed(t); ok {
// 		if pkg := named.Obj().Pkg(); pkg != nil {
// 			return PackagePath(pkg.Path())
// 		}
// 	}

// 	return ""
// }

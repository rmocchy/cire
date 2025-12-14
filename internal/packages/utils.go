package packages

import "go/types"

// derefType はポインタ型を再帰的に剥がす
func derefType(t types.Type) types.Type {
	for {
		ptr, ok := t.(*types.Pointer)
		if !ok {
			return t
		}
		t = ptr.Elem()
	}
}

// getTypeName は型の名前を取得
func getTypeName(t types.Type) string {
	t = derefType(t)
	t = types.Unalias(t)

	if named, ok := t.(*types.Named); ok {
		return named.Obj().Name()
	}

	return t.String()
}

// getPackagePath は型のパッケージパスを取得
func getPackagePath(t types.Type) string {
	t = derefType(t)
	t = types.Unalias(t)

	if named, ok := t.(*types.Named); ok {
		if pkg := named.Obj().Pkg(); pkg != nil {
			return pkg.Path()
		}
	}

	return ""
}

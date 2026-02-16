package core

import "go/types"

// 型変換ヘルパー - Goの型アサーションのように使える

// ConvertToPointer はポインタ型として取得する
func ConvertToPointer(t types.Type) (*types.Pointer, bool) {
	ptr, ok := t.(*types.Pointer)
	return ptr, ok
}

// ConvertToNamed はNamed型として取得する
func ConvertToNamed(t types.Type) (*types.Named, bool) {
	named, ok := types.Unalias(t).(*types.Named)
	return named, ok
}

// ConvertToNamedInterface はNamed型のInterfaceとして取得（型名も返す）
func ConvertToNamedInterface(t types.Type) (string, *types.Interface, bool) {
	named, ok := ConvertToNamed(t)
	if !ok {
		return "", nil, false
	}

	iface, ok := types.Unalias(named.Underlying()).(*types.Interface)
	if !ok {
		return "", nil, false
	}

	return named.Obj().Name(), iface, true
}

// ConvertToNamedStruct はNamed型のStructとして取得（型名も返す）
func ConvertToNamedStruct(t types.Type) (string, *types.Struct, bool) {
	named, ok := ConvertToNamed(t)
	if !ok {
		return "", nil, false
	}

	st, ok := types.Unalias(named.Underlying()).(*types.Struct)
	if !ok {
		return "", nil, false
	}

	return named.Obj().Name(), st, true
}

// Deref はポインタを剥がす（ポインタでない場合はそのまま返す）
func Deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}

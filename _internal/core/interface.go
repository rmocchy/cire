package core

import (
	"go/types"
)

// matchesInterface は型が指定されたinterfaceと一致するかチェック
func matchesInterface(t types.Type, targetInterface *types.Interface) bool {
	// ポインタを剥がす
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	// エイリアスを解決
	t = types.Unalias(t)

	// Named型かチェック
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	// 基底型がinterfaceであることをチェック
	underlying := named.Underlying()
	iface, ok := underlying.(*types.Interface)
	if !ok {
		return false
	}

	// interfaceが同一かをチェック
	return types.Identical(iface, targetInterface)
}

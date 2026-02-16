package core

import (
	"go/types"
	"testing"
)

func TestGetPackagePath(t *testing.T) {
	// 基本型の場合は空文字列
	if path := GetPackagePath(types.Typ[types.Int]); path != "" {
		t.Errorf("Expected empty path for basic type, got %s", path)
	}

	// Named型でない場合も空文字列
	if path := GetPackagePath(types.NewPointer(types.Typ[types.String])); path != "" {
		t.Errorf("Expected empty path for non-named type, got %s", path)
	}
}

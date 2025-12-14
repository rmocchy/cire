package file

import (
	"path/filepath"
	"testing"
)

func TestParseWireFileStructs(t *testing.T) {
	// sample/basic/wire.goファイルのパスを構築
	sampleWirePath := filepath.Join("..", "..", "sample", "basic", "wire.go")

	functions, err := ParseWireFileStructs(sampleWirePath)
	if err != nil {
		t.Fatalf("ParseWireFileStructs failed: %v", err)
	}

	// 少なくとも1つの関数が見つかるはず
	if len(functions) == 0 {
		t.Fatal("Expected at least one function, got none")
	}

	// InitializeUserHandler関数が存在するか確認
	found := false
	var targetFunc FunctionInfo
	for _, fn := range functions {
		if fn.Name == "InitializeUserHandler" {
			found = true
			targetFunc = fn
			break
		}
	}

	if !found {
		t.Fatal("InitializeUserHandler function not found")
	}

	// 返り値の構造体情報をチェック
	// InitializeUserHandler() (*ControllerSet, error) なので、
	// エラー以外の返り値は *ControllerSet のはず
	if len(targetFunc.ReturnTypes) != 1 {
		t.Fatalf("Expected 1 non-error return type, got %d: %v",
			len(targetFunc.ReturnTypes), targetFunc.ReturnTypes)
	}

	structInfo := targetFunc.ReturnTypes[0]

	// 構造体名のチェック
	expectedName := "ControllerSet"
	if structInfo.Name != expectedName {
		t.Errorf("Expected struct name %s, got %s", expectedName, structInfo.Name)
	}

	// ポインタかどうかのチェック
	if !structInfo.IsPointer {
		t.Error("Expected pointer type, but got non-pointer")
	}

	t.Logf("Successfully parsed function: %s", targetFunc.Name)
	t.Logf("  Struct: %s, IsPointer: %v", structInfo.Name, structInfo.IsPointer)
}

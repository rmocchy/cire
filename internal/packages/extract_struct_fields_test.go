package packages

import (
	"path/filepath"
	"testing"
)

func TestExtractStructFields(t *testing.T) {
	// テスト用の作業ディレクトリ（sample/basicをベースにする）
	workDir := filepath.Join("..", "..", "sample", "basic")

	tests := []struct {
		name         string
		packagePath  string
		structName   string
		wantErr      bool
		validateFunc func(*testing.T, *StructFieldsInfo)
	}{
		{
			name:        "UserHandler struct with interface field",
			packagePath: "./handler",
			structName:  "UserHandler",
			wantErr:     false,
			validateFunc: func(t *testing.T, info *StructFieldsInfo) {
				if info.StructName != "UserHandler" {
					t.Errorf("Expected struct name UserHandler, got %s", info.StructName)
				}

				if len(info.Fields) != 1 {
					t.Fatalf("Expected 1 field, got %d", len(info.Fields))
				}

				field := info.Fields[0]
				if field.Name != "service" {
					t.Errorf("Expected field name 'service', got '%s'", field.Name)
				}

				if field.TypeName != "UserService" {
					t.Errorf("Expected type name 'UserService', got '%s'", field.TypeName)
				}

				// パッケージパスの検証（ローカルパスなので実際のモジュール名が入る）
				if field.PackagePath == "" {
					t.Error("Expected non-empty package path")
				}
				t.Logf("Package path: %s", field.PackagePath)

				if field.IsPointer {
					t.Error("Expected non-pointer field, but got pointer")
				}

				if !field.IsInterface {
					t.Error("Expected interface type, but got non-interface")
				}

				t.Logf("Field: %s, Type: %s, Package: %s, IsPointer: %v, IsInterface: %v",
					field.Name, field.TypeName, field.PackagePath, field.IsPointer, field.IsInterface)
			},
		},
		{
			name:        "userServiceImpl struct with interface field",
			packagePath: "./service",
			structName:  "userServiceImpl",
			wantErr:     false,
			validateFunc: func(t *testing.T, info *StructFieldsInfo) {
				if info.StructName != "userServiceImpl" {
					t.Errorf("Expected struct name userServiceImpl, got %s", info.StructName)
				}

				if len(info.Fields) != 1 {
					t.Fatalf("Expected 1 field, got %d", len(info.Fields))
				}

				field := info.Fields[0]
				if field.Name != "repo" {
					t.Errorf("Expected field name 'repo', got '%s'", field.Name)
				}

				if field.TypeName != "UserRepository" {
					t.Errorf("Expected type name 'UserRepository', got '%s'", field.TypeName)
				}

				// パッケージパスの検証（ローカルパスなので実際のモジュール名が入る）
				if field.PackagePath == "" {
					t.Error("Expected non-empty package path")
				}
				t.Logf("Package path: %s", field.PackagePath)

				if field.IsPointer {
					t.Error("Expected non-pointer field, but got pointer")
				}

				if !field.IsInterface {
					t.Error("Expected interface type, but got non-interface")
				}

				t.Logf("Field: %s, Type: %s, Package: %s, IsPointer: %v, IsInterface: %v",
					field.Name, field.TypeName, field.PackagePath, field.IsPointer, field.IsInterface)
			},
		},
		{
			name:        "Non-existent struct should return error",
			packagePath: "./handler",
			structName:  "NonExistentStruct",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ExtractStructFields(workDir, tt.packagePath, tt.structName)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if info == nil {
				t.Fatal("Expected non-nil StructFieldsInfo")
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, info)
			}
		})
	}
}

func TestExtractStructFieldsWithPointerFields(t *testing.T) {
	// ポインタフィールドを持つ構造体のテスト用
	// 実際のサンプルにポインタフィールドがある場合はここでテスト
	t.Skip("Skipping pointer field test - add when sample with pointer fields is available")
}

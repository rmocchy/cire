package packages

import (
	"go/types"
	"testing"
)

func TestDerefType(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Type
		expected string
	}{
		{
			name:     "基本型",
			input:    types.Typ[types.Int],
			expected: "int",
		},
		{
			name:     "ポインタ型",
			input:    types.NewPointer(types.Typ[types.String]),
			expected: "string",
		},
		{
			name:     "二重ポインタ型",
			input:    types.NewPointer(types.NewPointer(types.Typ[types.Bool])),
			expected: "bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := derefType(tt.input)
			if result.String() != tt.expected {
				t.Errorf("derefType() = %v, want %v", result.String(), tt.expected)
			}
		})
	}
}

func TestGetTypeName(t *testing.T) {
	// Named型のテストのためにパッケージをロード
	pkg := types.NewPackage("test/pkg", "pkg")
	
	// 構造体型を作成
	structType := types.NewStruct(nil, nil)
	namedType := types.NewNamed(types.NewTypeName(0, pkg, "MyStruct", nil), structType, nil)

	tests := []struct {
		name     string
		input    types.Type
		expected string
	}{
		{
			name:     "基本型",
			input:    types.Typ[types.Int],
			expected: "int",
		},
		{
			name:     "Named型",
			input:    namedType,
			expected: "MyStruct",
		},
		{
			name:     "ポインタ型のNamed型",
			input:    types.NewPointer(namedType),
			expected: "MyStruct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTypeName(tt.input)
			if result != tt.expected {
				t.Errorf("getTypeName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPackagePath(t *testing.T) {
	// パッケージパスのテスト用にパッケージを作成
	pkg1 := types.NewPackage("github.com/example/pkg1", "pkg1")
	pkg2 := types.NewPackage("github.com/example/pkg2", "pkg2")

	// 構造体型を作成
	structType1 := types.NewStruct(nil, nil)
	namedType1 := types.NewNamed(types.NewTypeName(0, pkg1, "Type1", nil), structType1, nil)

	structType2 := types.NewStruct(nil, nil)
	namedType2 := types.NewNamed(types.NewTypeName(0, pkg2, "Type2", nil), structType2, nil)

	tests := []struct {
		name     string
		input    types.Type
		expected string
	}{
		{
			name:     "基本型（パッケージなし）",
			input:    types.Typ[types.Int],
			expected: "",
		},
		{
			name:     "Named型 - pkg1",
			input:    namedType1,
			expected: "github.com/example/pkg1",
		},
		{
			name:     "Named型 - pkg2",
			input:    namedType2,
			expected: "github.com/example/pkg2",
		},
		{
			name:     "ポインタ型のNamed型",
			input:    types.NewPointer(namedType1),
			expected: "github.com/example/pkg1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPackagePath(tt.input)
			if result != tt.expected {
				t.Errorf("getPackagePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

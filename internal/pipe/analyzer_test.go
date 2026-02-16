package pipe

import (
	"testing"
)

func TestNewWireAnalyzer(t *testing.T) {
	tests := []struct {
		name          string
		workDir       string
		searchPattern string
		wantErr       bool
	}{
		{
			name:          "valid sample/basic",
			workDir:       "../../sample/basic",
			searchPattern: "./...",
			wantErr:       false,
		},
		{
			name:          "invalid directory",
			workDir:       "/nonexistent",
			searchPattern: "./...",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer, err := NewWireAnalyzer(tt.workDir, tt.searchPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWireAnalyzer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && analyzer == nil {
				t.Error("NewWireAnalyzer() returned nil analyzer")
			}
		})
	}
}

func TestWireAnalyzer_AnalyzeStruct(t *testing.T) {
	workDir := "../../sample/basic"
	searchPattern := "./..."

	analyzer, err := NewWireAnalyzer(workDir, searchPattern)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	tests := []struct {
		name        string
		packagePath string
		structName  string
		wantErr     bool
		validate    func(*testing.T, *StructNode)
	}{
		{
			name:        "analyze ControllerSet",
			packagePath: "",
			structName:  "ControllerSet",
			wantErr:     false,
			validate: func(t *testing.T, node *StructNode) {
				if node.StructName != "ControllerSet" {
					t.Errorf("Expected StructName=ControllerSet, got %s", node.StructName)
				}
				if len(node.Fields) == 0 {
					t.Error("Expected at least one field")
				}
			},
		},
		{
			name:        "analyze UserHandler",
			packagePath: "github.com/rmocchy/convinient_wire/sample/basic/handler",
			structName:  "UserHandler",
			wantErr:     false,
			validate: func(t *testing.T, node *StructNode) {
				if node.StructName != "UserHandler" {
					t.Errorf("Expected StructName=UserHandler, got %s", node.StructName)
				}
				hasInitFunc := false
				for _, fn := range node.InitFunctions {
					if fn.Name == "NewUserHandler" {
						hasInitFunc = true
						break
					}
				}
				if !hasInitFunc {
					t.Error("Expected NewUserHandler in InitFunctions")
				}
			},
		},
		{
			name:        "non-existent struct",
			packagePath: "",
			structName:  "NonExistentStruct",
			wantErr:     true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzer.AnalyzeStruct(tt.packagePath, tt.structName)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestIsBuiltinType(t *testing.T) {
	tests := []struct {
		typeName string
		want     bool
	}{
		{"string", true},
		{"int", true},
		{"bool", true},
		{"error", true},
		{"UserHandler", false},
		{"Config", false},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			if got := isBuiltinType(tt.typeName); got != tt.want {
				t.Errorf("isBuiltinType(%s) = %v, want %v", tt.typeName, got, tt.want)
			}
		})
	}
}

func TestFindInitFunctions(t *testing.T) {
	workDir := "../../sample/basic"
	searchPattern := "./..."

	analyzer, err := NewWireAnalyzer(workDir, searchPattern)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Config構造体を解析して初期化関数が見つかることを確認
	result, err := analyzer.AnalyzeStruct(
		"github.com/rmocchy/convinient_wire/sample/basic/repository",
		"Config",
	)
	if err != nil {
		t.Fatalf("Failed to analyze Config struct: %v", err)
	}

	hasNewConfig := false
	for _, fn := range result.InitFunctions {
		if fn.Name == "NewConfig" {
			hasNewConfig = true
			if fn.PackagePath != "github.com/rmocchy/convinient_wire/sample/basic/repository" {
				t.Errorf("Expected PackagePath=...repository, got %s", fn.PackagePath)
			}
		}
	}

	if !hasNewConfig {
		t.Error("Expected NewConfig in init functions")
	}
}

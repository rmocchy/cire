package pipe

import (
	"go/types"
	"testing"

	"github.com/rmocchy/convinient_wire/internal/cache"
	"github.com/rmocchy/convinient_wire/internal/core"
	"golang.org/x/tools/go/packages"
)

// テスト用のヘルパー関数: パッケージをロードしてキャッシュを作成
func setupTestAnalyzer(t *testing.T, workDir string) (*WireAnalyzer, error) {
	t.Helper()

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedDeps,
		Dir: workDir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	// エラーチェック
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, pkg.Errors[0]
		}
	}

	// types.Packageの配列を作成
	typesPkgs := make([]*types.Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		if pkg.Types != nil {
			typesPkgs = append(typesPkgs, pkg.Types)
		}
	}

	// キャッシュを作成
	functionCache := cache.NewFunctionCache(pkgs)
	structCache := cache.NewStructCache(typesPkgs)

	// WireAnalyzerを作成
	return NewWireAnalyzer(functionCache, structCache)
}

func TestNewWireAnalyzer(t *testing.T) {
	tests := []struct {
		name    string
		workDir string
		wantErr bool
	}{
		{
			name:    "valid sample/basic",
			workDir: "../../sample/basic",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer, err := setupTestAnalyzer(t, tt.workDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupTestAnalyzer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && analyzer == nil {
				t.Error("setupTestAnalyzer() returned nil analyzer")
			}
		})
	}
}

func TestWireAnalyzer_AnalyzeStruct(t *testing.T) {
	workDir := "../../sample/basic"

	analyzer, err := setupTestAnalyzer(t, workDir)
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
					if fn.Name() == "NewUserHandler" {
						hasInitFunc = true
						break
					}
				}
				if !hasInitFunc {
					t.Error("Expected NewUserHandler in InitFunctions")
				}
				// UserHandlerはserviceフィールドを持つ
				if len(node.Fields) == 0 {
					t.Error("Expected at least one field")
				}
				// serviceフィールドがインターフェースであることを確認
				hasServiceField := false
				for _, field := range node.Fields {
					if field.GetFieldName() == "service" {
						hasServiceField = true
						if field.NodeType() != NodeTypeInterface {
							t.Errorf("Expected service field to be interface, got %v", field.NodeType())
						}
					}
				}
				if !hasServiceField {
					t.Error("Expected service field in UserHandler")
				}
			},
		},
		{
			name:        "analyze Config",
			packagePath: "github.com/rmocchy/convinient_wire/sample/basic/repository",
			structName:  "Config",
			wantErr:     false,
			validate: func(t *testing.T, node *StructNode) {
				if node.StructName != "Config" {
					t.Errorf("Expected StructName=Config, got %s", node.StructName)
				}
				hasInitFunc := false
				for _, fn := range node.InitFunctions {
					if fn.Name() == "NewConfig" {
						hasInitFunc = true
						break
					}
				}
				if !hasInitFunc {
					t.Error("Expected NewConfig in InitFunctions")
				}
				// Configはビルトイン型のフィールドを持つ
				if len(node.Fields) == 0 {
					t.Error("Expected at least one field")
				}
			},
		},
		{
			name:        "non-existent struct",
			packagePath: "github.com/rmocchy/convinient_wire/sample/basic",
			structName:  "NonExistentStruct",
			wantErr:     true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzer.AnalyzeStruct(tt.structName, core.NewPackagePath(tt.packagePath))
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
		name     string
		typeName string
		want     bool
	}{
		{"string type", "string", true},
		{"int type", "int", true},
		{"bool type", "bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実際のtypes.Basicを作成してテスト
			var typ types.Type
			switch tt.typeName {
			case "string":
				typ = types.Typ[types.String]
			case "int":
				typ = types.Typ[types.Int]
			case "bool":
				typ = types.Typ[types.Bool]
			}
			if got := isBuiltinType(typ); got != tt.want {
				t.Errorf("isBuiltinType(%s) = %v, want %v", tt.typeName, got, tt.want)
			}
		})
	}
}

func TestWireAnalyzer_AnalyzeStructWithNestedFields(t *testing.T) {
	workDir := "../../sample/basic"

	analyzer, err := setupTestAnalyzer(t, workDir)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Config構造体を解析（ビルトインフィールドを含む）
	result, err := analyzer.AnalyzeStruct(
		"Config",
		core.NewPackagePath("github.com/rmocchy/convinient_wire/sample/basic/repository"),
	)
	if err != nil {
		t.Fatalf("Failed to analyze Config struct: %v", err)
	}

	if result.StructName != "Config" {
		t.Errorf("Expected StructName=Config, got %s", result.StructName)
	}

	// DSNとMaxPoolSizeフィールドがビルトイン型として検出されることを確認
	hasDSN := false
	hasMaxPoolSize := false
	for _, field := range result.Fields {
		if field.GetFieldName() == "DSN" && field.NodeType() == NodeTypeBuiltin {
			hasDSN = true
		}
		if field.GetFieldName() == "MaxPoolSize" && field.NodeType() == NodeTypeBuiltin {
			hasMaxPoolSize = true
		}
	}

	if !hasDSN {
		t.Error("Expected DSN field of builtin type")
	}
	if !hasMaxPoolSize {
		t.Error("Expected MaxPoolSize field of builtin type")
	}
}

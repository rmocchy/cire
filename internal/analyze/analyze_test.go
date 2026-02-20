package analyze

import (
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
)

// loadTestPackages はテスト用のパッケージをロードする
func loadTestPackages(t *testing.T, workDir string) []*packages.Package {
	t.Helper()

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes |
			packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedDeps,
		Dir: workDir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			t.Fatalf("Package error: %v", pkg.Errors[0])
		}
	}

	return pkgs
}

// setupTestAnalyzer はテスト用のAnalyzerをセットアップする
func setupTestAnalyzer(t *testing.T, workDir string) (Analyze, []*packages.Package) {
	t.Helper()

	pkgs := loadTestPackages(t, workDir)
	functionCache := NewFunctionCache(pkgs)
	analysisCache := NewAnalysisCache()
	analyzer := NewAnalyze(functionCache, analysisCache)

	return analyzer, pkgs
}

// findNamedType は指定したパッケージパスと型名から *types.Named を探す
func findNamedType(t *testing.T, pkgs []*packages.Package, pkgPath, typeName string) *types.Named {
	t.Helper()

	for _, pkg := range pkgs {
		if pkg.PkgPath != pkgPath {
			continue
		}
		obj := pkg.Types.Scope().Lookup(typeName)
		if obj == nil {
			continue
		}
		named, ok := obj.Type().(*types.Named)
		if !ok {
			t.Fatalf("%s in %s is not a named type", typeName, pkgPath)
		}
		return named
	}

	t.Fatalf("Named type %s not found in package %s", typeName, pkgPath)
	return nil
}

func TestNewAnalyze(t *testing.T) {
	tests := []struct {
		name    string
		workDir string
	}{
		{
			name:    "valid sample/basic",
			workDir: "../../sample/basic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer, _ := setupTestAnalyzer(t, tt.workDir)
			if analyzer == nil {
				t.Error("setupTestAnalyzer() returned nil analyzer")
			}
		})
	}
}

func TestFunctionCache_BulkGet(t *testing.T) {
	workDir := "../../sample/basic"
	pkgs := loadTestPackages(t, workDir)
	functionCache := NewFunctionCache(pkgs)

	tests := []struct {
		name        string
		packagePath string
		typeName    string
		wantFuncs   []string
	}{
		{
			// NewConfig() *Config — ポインタ返し
			name:        "get function returning *Config",
			packagePath: "github.com/rmocchy/cire/sample/basic/repository",
			typeName:    "Config",
			wantFuncs:   []string{"NewConfig"},
		},
		{
			// NewUserService() UserService — インターフェース返し
			name:        "get function returning UserService",
			packagePath: "github.com/rmocchy/cire/sample/basic/service",
			typeName:    "UserService",
			wantFuncs:   []string{"NewUserService"},
		},
		{
			// NewUserRepository() (UserRepository, error) — インターフェース返し
			name:        "get function returning UserRepository",
			packagePath: "github.com/rmocchy/cire/sample/basic/repository",
			typeName:    "UserRepository",
			wantFuncs:   []string{"NewUserRepository"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnType := findNamedType(t, pkgs, tt.packagePath, tt.typeName)
			fns := functionCache.BulkGet(returnType)

			for _, wantFunc := range tt.wantFuncs {
				found := false
				for _, fn := range fns {
					t.Logf("Found function: %s in %s", fn.Name(), fn.Pkg().Path())
					if fn.Name() == wantFunc {
						found = true
					}
				}
				if !found {
					t.Errorf("Expected function %s not found", wantFunc)
				}
			}
		})
	}
}

func TestAnalyze_ExecuteFromStruct(t *testing.T) {
	workDir := "../../sample/basic"
	analyzer, pkgs := setupTestAnalyzer(t, workDir)

	tests := []struct {
		name        string
		packagePath string
		structName  string
		wantErr     bool
		wantFuncs   []string
	}{
		{
			name:        "analyze App struct — all providers found",
			packagePath: "github.com/rmocchy/cire/sample/basic",
			structName:  "App",
			wantErr:     false,
			wantFuncs:   []string{"NewUserHandler", "NewUserService", "NewUserRepository", "NewConfig"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namedType := findNamedType(t, pkgs, tt.packagePath, tt.structName)

			nodes, err := analyzer.ExecuteFromStruct(namedType)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ExecuteFromStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// 全ノードから名前を収集（再帰的に）
			allNames := collectNodeNames(nodes)
			for _, wantFunc := range tt.wantFuncs {
				if !allNames[wantFunc] {
					t.Errorf("Expected function %s not found in nodes: %v", wantFunc, allNames)
				}
			}
		})
	}
}

// collectNodeNames はノードツリーから全ての関数名を再帰的に収集する
func collectNodeNames(nodes []*FnDITreeNode) map[string]bool {
	names := make(map[string]bool)
	for _, node := range nodes {
		if node == nil {
			continue
		}
		names[node.Name] = true
		for k, v := range collectNodeNames(node.Childs) {
			names[k] = v
		}
	}
	return names
}

func TestAnalysisCache(t *testing.T) {
	workDir := "../../sample/basic"
	pkgs := loadTestPackages(t, workDir)
	cache := NewAnalysisCache()

	namedType := findNamedType(t, pkgs, "github.com/rmocchy/cire/sample/basic/repository", "Config")

	// 初期状態では空
	_, found := cache.Get(namedType)
	if found {
		t.Error("Expected cache to be empty initially")
	}

	// キャッシュに追加
	testNodes := []*FnDITreeNode{
		{Name: "TestFunc", PkgPath: "test/path", Childs: nil},
	}
	cache.Set(namedType, testNodes)

	// キャッシュから取得
	result, found := cache.Get(namedType)
	if !found {
		t.Error("Expected cache to contain the value")
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 node, got %d", len(result))
	}
	if result[0].Name != "TestFunc" {
		t.Errorf("Expected TestFunc, got %s", result[0].Name)
	}
}

func TestFnDITreeNode(t *testing.T) {
	child := &FnDITreeNode{
		Name:    "ChildFunc",
		PkgPath: "test/child",
		Childs:  nil,
	}

	parent := &FnDITreeNode{
		Name:    "ParentFunc",
		PkgPath: "test/parent",
		Childs:  []*FnDITreeNode{child},
	}

	if parent.Name != "ParentFunc" {
		t.Errorf("Expected ParentFunc, got %s", parent.Name)
	}
	if parent.PkgPath != "test/parent" {
		t.Errorf("Expected test/parent, got %s", parent.PkgPath)
	}
	if len(parent.Childs) != 1 {
		t.Errorf("Expected 1 child, got %d", len(parent.Childs))
	}
	if parent.Childs[0].Name != "ChildFunc" {
		t.Errorf("Expected ChildFunc, got %s", parent.Childs[0].Name)
	}
}

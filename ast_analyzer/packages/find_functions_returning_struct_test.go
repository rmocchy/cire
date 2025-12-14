package packages

import (
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestFindFunctionsReturningStruct(t *testing.T) {
	tests := []struct {
		name          string
		structName    string
		structPkgPath string
		wantFuncNames []string
	}{
		{
			name:          "Find functions returning Config struct",
			structName:    "Config",
			structPkgPath: "github.com/rmocchy/convinient_wire/sample/basic/repository",
			wantFuncNames: []string{"NewConfig"},
		},
		{
			name:          "Find functions returning UserHandler struct",
			structName:    "UserHandler",
			structPkgPath: "github.com/rmocchy/convinient_wire/sample/basic/handler",
			wantFuncNames: []string{"NewUserHandler"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &packages.Config{
				Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports |
					packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
				Dir: "../../sample/basic",
			}

			pkgs, err := packages.Load(cfg, "./...")
			if err != nil {
				t.Fatalf("failed to load packages: %v", err)
			}

			if packages.PrintErrors(pkgs) > 0 {
				t.Fatal("packages have errors")
			}

			functions := FindFunctionsReturningStruct(tt.structName, tt.structPkgPath, pkgs)

			if len(functions) == 0 {
				t.Errorf("no functions found returning %s", tt.structName)
				return
			}

			// 関数名をマップに変換してチェックしやすくする
			foundNames := make(map[string]bool)
			for _, fn := range functions {
				foundNames[fn.Name] = true
				t.Logf("Found function: %s (package: %s)", fn.Name, fn.PackagePath)
			}

			// 期待される関数名がすべて見つかったかチェック
			for _, wantName := range tt.wantFuncNames {
				if !foundNames[wantName] {
					t.Errorf("expected function %s not found", wantName)
				}
			}
		})
	}
}

func TestMatchesStructType(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports |
			packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir: "../../sample/basic",
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		t.Fatalf("failed to load packages: %v", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("packages have errors")
	}

	if len(pkgs) == 0 {
		t.Fatal("no packages loaded")
	}

	// repositoryパッケージを探す
	var repoPkg *packages.Package
	for _, pkg := range pkgs {
		if pkg.PkgPath == "github.com/rmocchy/convinient_wire/sample/basic/repository" {
			repoPkg = pkg
			break
		}
	}

	if repoPkg == nil {
		t.Fatal("repository package not found")
	}

	scope := repoPkg.Types.Scope()

	// NewConfig関数を取得（構造体を返す関数）
	obj := scope.Lookup("NewConfig")
	if obj == nil {
		t.Fatal("NewConfig not found")
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		t.Fatal("NewConfig is not a function")
	}

	sig := fn.Type().(*types.Signature)
	results := sig.Results()

	if results.Len() == 0 {
		t.Fatal("NewConfig has no return values")
	}

	returnType := results.At(0).Type()

	tests := []struct {
		name          string
		structName    string
		structPkgPath string
		want          bool
	}{
		{
			name:          "Matches Config struct",
			structName:    "Config",
			structPkgPath: "github.com/rmocchy/convinient_wire/sample/basic/repository",
			want:          true,
		},
		{
			name:          "Does not match wrong struct name",
			structName:    "WrongConfig",
			structPkgPath: "github.com/rmocchy/convinient_wire/sample/basic/repository",
			want:          false,
		},
		{
			name:          "Does not match wrong package path",
			structName:    "Config",
			structPkgPath: "github.com/wrong/package",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesStructType(returnType, tt.structName, tt.structPkgPath)
			if got != tt.want {
				t.Errorf("matchesStructType() = %v, want %v", got, tt.want)
			}
		})
	}
}

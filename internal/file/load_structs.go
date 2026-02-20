package file

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func LoadNamedStructs(path string, pkgs []*packages.Package) ([]*types.Named, error) {
	structNames, err := getStructNames(path)
	if err != nil {
		return nil, err
	}
	pkgPath, err := ResolvePackagePath(path)
	if err != nil {
		return nil, err
	}

	namedStructs := make([]*types.Named, 0)
	for _, p := range pkgs {
		if p.PkgPath != *pkgPath {
			continue
		}
		for _, name := range structNames {
			obj := p.Types.Scope().Lookup(name)
			named, ok := obj.Type().(*types.Named)
			if !ok {
				continue
			}
			if _, ok := named.Underlying().(*types.Struct); ok {
				namedStructs = append(namedStructs, named)
			}
		}
	}

	return namedStructs, nil
}

func getStructNames(path string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	results := make([]string, 0)

	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 構造体のみ対象
			if _, ok := typeSpec.Type.(*ast.StructType); !ok {
				continue
			}

			// すべての構造体を追加
			results = append(results, typeSpec.Name.Name)
		}

		return true
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no structs found in %s", path)
	}

	return results, nil
}

func ResolvePackagePath(filePath string) (*string, error) {
	dir := filepath.Dir(filePath)

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedModule,
		Dir:  dir,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no package found for file: %s", filePath)
	}

	pkgPath := pkgs[0].PkgPath
	return &pkgPath, nil
}

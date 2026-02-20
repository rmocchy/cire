package file

import (
	"fmt"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func LoadNamedStructs(path string, pkgs []*packages.Package) ([]*types.Named, error) {
	pkgPath, err := resolvePackagePath(path)
	if err != nil {
		return nil, err
	}

	namedStructs := make([]*types.Named, 0)
	for _, p := range pkgs {
		if p.PkgPath != *pkgPath {
			continue
		}
		for _, name := range p.Types.Scope().Names() {
			obj := p.Types.Scope().Lookup(name)
			// 型宣言（TypeName）のみを対象とする。変数（Var）はその型が外部パッケージの
			// 構造体であっても解析対象に含めない（例: var AppSet = wire.NewSet(...) が
			// wire.ProviderSet として混入するのを防ぐ）
			if _, ok := obj.(*types.TypeName); !ok {
				continue
			}
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

func resolvePackagePath(filePath string) (*string, error) {
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

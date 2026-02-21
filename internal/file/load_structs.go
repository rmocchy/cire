package file

import (
	"fmt"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func LoadNamedStructs(path string, pkgs []*packages.Package) ([]*types.Named, error) {
	fileName := filepath.Base(path)
	if fileName == "" {
		return nil, fmt.Errorf("invalid file path: %s", path)
	}
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
			if _, ok := obj.(*types.TypeName); !ok {
				continue
			}
			named, ok := obj.Type().(*types.Named)
			if !ok {
				continue
			}
			pos := named.Obj().Pos()
			position := p.Fset.Position(pos)
			defFileName := filepath.Base(position.Filename)
			if defFileName != fileName {
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

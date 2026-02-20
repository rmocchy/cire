package file

import (
	"fmt"
	"go/parser"
	"go/token"
	"path"
)

func ExtractPackageName(filePath string) (*string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	return &f.Name.Name, nil
}

// PkgNameFromPath は pkgPath（例: "github.com/foo/bar/baz"）から
// パッケージ名（例: "baz"）を返す。
func PkgNameFromPath(pkgPath string) string {
	return path.Base(pkgPath)
}

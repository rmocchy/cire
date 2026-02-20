package file

import (
	"fmt"
	"go/parser"
	"go/token"
)

func ExtractPackageName(filePath string) (*string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}
	return &f.Name.Name, nil
}

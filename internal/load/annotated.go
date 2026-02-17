package load

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
)

// FindAnnotatedStructs は指定されたファイルからすべての構造体を検出する
// ファイルは //go:build cire ビルドタグを持つ専用ファイルであることを想定
func FindAnnotatedStructs(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
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
		return nil, fmt.Errorf("no structs found in %s", filePath)
	}

	return results, nil
}

// hasCireBuildTag はファイルに //go:build cire ビルドタグがあるかをチェックする
func hasCireBuildTag(filePath string) (bool, error) {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)

	ctxDefault := build.Default
	ctxCire := build.Default
	ctxCire.BuildTags = append(ctxCire.BuildTags, "wireinject")

	defaultIncluded, err := ctxDefault.MatchFile(dir, base)
	if err != nil {
		return false, err
	}
	cireIncluded, err := ctxCire.MatchFile(dir, base)
	if err != nil {
		return false, err
	}
	return !defaultIncluded && cireIncluded, nil
}

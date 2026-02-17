package load

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// FindAnnotatedStructs は指定されたファイルから @cire アノテーション付きの構造体を検出する
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

			// @cire アノテーションがあれば追加（typeSpec.Doc または genDecl.Doc）
			if hasAnnotation(typeSpec.Doc, "@cire") || hasAnnotation(genDecl.Doc, "@cire") {
				results = append(results, typeSpec.Name.Name)
			}
		}

		return true
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no structs with @cire annotation found in %s", filePath)
	}

	return results, nil
}

// hasAnnotation はコメントに指定されたアノテーションが含まれるかをチェックする
func hasAnnotation(doc *ast.CommentGroup, annotation string) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, annotation) {
			return true
		}
	}
	return false
}

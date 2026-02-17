package load

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// AnnotatedStruct はアノテーション付きの構造体情報
type AnnotatedStruct struct {
	Name        string
	PackagePath string
	PackageName string
}

// FindAnnotatedStructs は指定されたファイルから @cire アノテーション付きの構造体を検出する
func FindAnnotatedStructs(filePath string) ([]AnnotatedStruct, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var results []AnnotatedStruct

	// AST をトラバースして型宣言を探す
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		// 各型定義をチェック
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 構造体でない場合はスキップ
			if _, ok := typeSpec.Type.(*ast.StructType); !ok {
				continue
			}

			// @cire アノテーションのチェック
			if !hasCireAnnotation(typeSpec.Doc) {
				continue
			}

			results = append(results, AnnotatedStruct{
				Name:        typeSpec.Name.Name,
				PackageName: node.Name.Name,
				PackagePath: "", // 後で解決
			})
		}

		return true
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no structs with @cire annotation found in %s", filePath)
	}

	return results, nil
}

// hasCireAnnotation はコメントに @cire アノテーションが含まれるかをチェックする
func hasCireAnnotation(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}

	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "@cire") {
			return true
		}
	}

	return false
}

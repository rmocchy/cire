package file

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// ParseWireFileStructs はwire.goファイルをパースしてpublic関数の返り値構造体情報を取得する
func ParseWireFileStructs(filepath string) ([]FunctionInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// importパスのマッピングを取得
	importMap := extractImports(node)

	var functions []FunctionInfo

	// ASTを走査して関数宣言を探す
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// public関数のみ(先頭が大文字)
		if !isPublicFunction(funcDecl.Name.Name) {
			return true
		}

		// 返り値の構造体情報を取得
		returnTypes := extractStructTypes(funcDecl.Type.Results, importMap)

		functions = append(functions, FunctionInfo{
			Name:        funcDecl.Name.Name,
			ReturnTypes: returnTypes,
		})

		return true
	})

	return functions, nil
}

// extractImports はimport文からパッケージ名とパスのマッピングを作成
func extractImports(node *ast.File) map[string]string {
	importMap := make(map[string]string)

	for _, imp := range node.Imports {
		if imp.Path == nil {
			continue
		}

		// import pathから引用符を除去
		path := imp.Path.Value[1 : len(imp.Path.Value)-1]

		// パッケージ名を取得
		var pkgName string
		if imp.Name != nil {
			// エイリアスがある場合
			pkgName = imp.Name.Name
		} else {
			// デフォルトはパスの最後の部分
			pkgName = getPackageNameFromPath(path)
		}

		importMap[pkgName] = path
	}

	return importMap
}

// getPackageNameFromPath はimportパスからパッケージ名を取得
func getPackageNameFromPath(path string) string {
	// 最後の/以降を取得
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

// isPublicFunction は関数名がpublic(先頭が大文字)かどうかを判定
func isPublicFunction(name string) bool {
	if len(name) == 0 {
		return false
	}
	firstChar := rune(name[0])
	return firstChar >= 'A' && firstChar <= 'Z'
}

// extractStructTypes は返り値から構造体型の情報を抽出
func extractStructTypes(results *ast.FieldList, importMap map[string]string) []StructInfo {
	if results == nil {
		return nil
	}

	var structs []StructInfo
	for _, field := range results.List {
		structInfo := parseStructType(field.Type, importMap)

		// error型と空の構造体情報は除外
		if structInfo.Name != "" && structInfo.Name != "error" {
			structs = append(structs, structInfo)
		}
	}

	return structs
}

// parseStructType は式から構造体情報を抽出
func parseStructType(expr ast.Expr, importMap map[string]string) StructInfo {
	switch t := expr.(type) {
	case *ast.Ident:
		// ローカルパッケージの型またはerror
		return StructInfo{
			Name:      t.Name,
			IsPointer: false,
		}
	case *ast.StarExpr:
		// ポインタ型
		info := parseStructType(t.X, importMap)
		info.IsPointer = true
		return info
	case *ast.SelectorExpr:
		// パッケージ.型名の形式
		if _, ok := t.X.(*ast.Ident); ok {
			structName := t.Sel.Name

			return StructInfo{
				Name:      structName,
				IsPointer: false,
			}
		}
	}

	return StructInfo{}
}

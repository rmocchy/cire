package core

// 作ろうと思ったが複雑化するのでv1ではやめる

// type retType struct {
// 	Type types.Type
// 	IsInterface bool
// 	IsPointer bool
// }

// // FindStructsReturnedByFunction は指定された関数の return 文から実際に返される構造体の型を取得する
// // 関数の返り値がインターフェースの場合、そのインターフェースを実装している型のみを返す
// func FindStructsReturnedByFunction(pkg *packages.Package, fn *types.Func) [][]types.Type {
// 	rets := fn.Signature().Results()
// 	if rets == nil {
// 		return nil
// 	}

// 	// 返り値の型をインデックスごとに保存
// 	declaredTypes := make([]types.Type, 0, rets.Len())

// 	for i := 0; i < rets.Len(); i++ {
// 		t := rets.At(i).Type()
// 		declaredTypes[i] = t

// 		if _, ok := t.Underlying().(*types.Interface); ok {
// 	}

// 	result := make([]types.Type, 0)

// 	// 関数のシグネチャから返り値の型を取得
// 	sig := fn.Signature()
// 	results := sig.Results()
// 	if results == nil {
// 		return result
// 	}

// 	// 返り値の型をインデックスごとに保存
// 	declaredTypes := make([]types.Type, results.Len())
// 	for i := 0; i < results.Len(); i++ {
// 		declaredTypes[i] = results.At(i).Type()
// 	}

// 	// パッケージの AST を走査
// 	for _, file := range pkg.Syntax {
// 		ast.Inspect(file, func(n ast.Node) bool {
// 			// 関数宣言を探す
// 			funcDecl, ok := n.(*ast.FuncDecl)
// 			if !ok {
// 				return true
// 			}

// 			// 対象の関数か確認
// 			funcObj := pkg.TypesInfo.Defs[funcDecl.Name]
// 			if funcObj != fn {
// 				return true
// 			}
// 			if funcDecl.Body == nil {
// 				return true
// 			}

// 			// return 文を探す
// 			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
// 				retStmt, ok := n.(*ast.ReturnStmt)
// 				if !ok {
// 					return true
// 				}

// 				// 各 return 式の型を取得
// 				for i, expr := range retStmt.Results {
// 					// TypeOf で式の型を解決
// 					exprType := pkg.TypesInfo.TypeOf(expr)
// 					if exprType == nil {
// 						continue
// 					}

// 					// 宣言された型がインターフェースの場合、実装チェック
// 					if i < len(declaredTypes) {
// 						if ifaceType, ok := declaredTypes[i].Underlying().(*types.Interface); ok {
// 							// インターフェースを実装しているかチェック
// 							if types.Implements(exprType, ifaceType) {
// 								result = append(result, exprType)
// 							}
// 						} else {
// 							// インターフェースでない場合はそのまま追加
// 							result = append(result, exprType)
// 						}
// 					} else {
// 						result = append(result, exprType)
// 					}
// 				}

// 				return true
// 			})

// 			return false
// 		})
// 	}

// 	return result
// }

package packages

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// InterfaceReference はインターフェースを参照する関数の情報を保持する
type InterfaceReference struct {
	FunctionName        string // 関数名
	PackagePath         string // 関数が定義されているパッケージパス
	ImplementingType    string // 対応づけられた実装型の名前
	ImplementingPkgPath string // 実装型のパッケージパス
}

// FindInterfaceReferences は指定されたインターフェースを参照する関数とそこで対応づけられた構造体を返す
// workDir: パッケージ解決の基準となる作業ディレクトリ
// interfaceName: 検索するインターフェースの名前
// interfacePkgPath: インターフェースが定義されているパッケージパス
// searchPattern: 検索対象のパッケージパターン（例: "./...", "github.com/user/repo/..."）
func FindInterfaceReferences(workDir, interfaceName, interfacePkgPath, searchPattern string) ([]InterfaceReference, error) {
	// パッケージをロード
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports,
		Dir: workDir,
	}

	pkgs, err := packages.Load(cfg, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found for pattern: %s", searchPattern)
	}

	var references []InterfaceReference

	// 各パッケージを検索
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			// エラーがあっても他のパッケージを処理
			continue
		}

		pkgRefs := findReferencesInPackage(pkg, interfaceName, interfacePkgPath)
		references = append(references, pkgRefs...)
	}

	return references, nil
}

// findReferencesInPackage は特定のパッケージ内でインターフェース参照を検索
func findReferencesInPackage(pkg *packages.Package, interfaceName, interfacePkgPath string) []InterfaceReference {
	var references []InterfaceReference

	// パッケージ内の全ての関数を走査
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			// 関数宣言を探す
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// コンストラクタ関数かどうか（Newで始まる）をチェック
			// また、戻り値の型がインターフェースかどうかをチェック
			refs := checkFunctionForInterface(pkg, funcDecl, interfaceName, interfacePkgPath)
			references = append(references, refs...)

			return true
		})
	}

	return references
}

// checkFunctionForInterface は関数がインターフェースを参照しているかチェック
func checkFunctionForInterface(pkg *packages.Package, funcDecl *ast.FuncDecl, interfaceName, interfacePkgPath string) []InterfaceReference {
	var references []InterfaceReference

	if funcDecl.Type.Results == nil {
		return references
	}

	// 戻り値の型をチェック
	for _, result := range funcDecl.Type.Results.List {
		resultType := pkg.TypesInfo.TypeOf(result.Type)
		if resultType == nil {
			continue
		}

		// インターフェース型かチェック
		if isTargetInterface(resultType, interfaceName, interfacePkgPath) {
			// 関数本体から実装型を探す
			implType := findImplementingType(pkg, funcDecl, resultType)
			if implType != nil {
				ref := InterfaceReference{
					FunctionName:        funcDecl.Name.Name,
					PackagePath:         pkg.PkgPath,
					ImplementingType:    getTypeName(implType),
					ImplementingPkgPath: getPackagePath(implType),
				}
				references = append(references, ref)
			}
		}
	}

	return references
}

// isTargetInterface は型が指定されたインターフェースかどうかをチェック
func isTargetInterface(t types.Type, interfaceName, interfacePkgPath string) bool {
	// ポインタ型の場合は剥がす
	t = derefType(t)
	t = types.Unalias(t)

	// Named型かチェック
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	// インターフェース型かチェック
	_, isInterface := named.Underlying().(*types.Interface)
	if !isInterface {
		return false
	}

	// 名前とパッケージパスが一致するかチェック
	obj := named.Obj()
	if obj.Name() != interfaceName {
		return false
	}

	if obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == interfacePkgPath
}

// findImplementingType は関数本体から実装型を探す
func findImplementingType(pkg *packages.Package, funcDecl *ast.FuncDecl, interfaceType types.Type) types.Type {
	if funcDecl.Body == nil {
		return nil
	}

	var implType types.Type

	// 関数本体を走査してreturn文を探す
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		retStmt, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}

		// return文の値を検査
		for _, expr := range retStmt.Results {
			exprType := pkg.TypesInfo.TypeOf(expr)
			if exprType == nil {
				continue
			}

			// インターフェース型に代入可能な具象型を見つける
			if types.AssignableTo(exprType, interfaceType) {
				// ポインタの場合は元の型を取得
				if ptr, ok := exprType.(*types.Pointer); ok {
					if named, ok := ptr.Elem().(*types.Named); ok {
						implType = named
						return false
					}
				}
				if named, ok := exprType.(*types.Named); ok {
					implType = named
					return false
				}
			}
		}

		return true
	})

	return implType
}

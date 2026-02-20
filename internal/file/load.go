package file

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// LoadPackagesFromFile は指定されたファイルからパッケージをロードする
func LoadAllPkgsFromPath(path string) ([]*packages.Package, error) {
	// ファイルのディレクトリを取得
	dir := filepath.Dir(path)

	// go.modファイルを探してモジュールルートを見つける
	moduleRoot := findModuleRoot(dir)
	if moduleRoot == "" {
		moduleRoot = dir
	}

	// パッケージのロード設定
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Dir: moduleRoot, // モジュールルートを設定
	}

	// ファイルが含まれるパッケージとその依存関係をロード
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}

	return pkgs, nil
}

// findModuleRoot はgo.modファイルを探してモジュールルートを返す
func findModuleRoot(startDir string) string {
	dir := startDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// ルートディレクトリに到達した
			return ""
		}
		dir = parent
	}
}

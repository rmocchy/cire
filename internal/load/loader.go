package load

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmocchy/convinient_wire/internal/cache"
	"github.com/rmocchy/convinient_wire/internal/core"
	"golang.org/x/tools/go/packages"
)

// PackageLoader はパッケージをロードしてキャッシュを構築する
type PackageLoader struct {
	FunctionCache core.FunctionCache
	StructCache   core.StructCache
}

// LoadPackages は指定されたパッケージパターンからパッケージをロードし、キャッシュを構築する
// packagePath が空文字列の場合は "./..." をロードする
func LoadPackages(packagePath string) (*PackageLoader, error) {
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
	}

	// カレントディレクトリまたは指定されたパッケージからロード
	pattern := "./..."
	if packagePath != "" {
		pattern = packagePath
	}

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}

	functionCache := cache.NewFunctionCache(pkgs)
	structCache := cache.NewStructCache(pkgs)

	return &PackageLoader{
		FunctionCache: functionCache,
		StructCache:   structCache,
	}, nil
}

// LoadPackagesFromFile は指定されたファイルからパッケージをロードする
func LoadPackagesFromFile(path string) (*PackageLoader, error) {
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

	// キャッシュの構築
	functionCache := cache.NewFunctionCache(pkgs)
	structCache := cache.NewStructCache(pkgs)

	return &PackageLoader{
		FunctionCache: functionCache,
		StructCache:   structCache,
	}, nil
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

// ResolvePackagePath は指定されたファイルの実際のパッケージパスを解決する
func ResolvePackagePath(filePath string) (*core.PackagePath, error) {
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

	pkgPath := core.NewPackagePath(pkgs[0].PkgPath)
	return &pkgPath, nil
}

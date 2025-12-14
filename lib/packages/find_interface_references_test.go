package packages

import (
	"testing"
)

func TestFindInterfaceReferences(t *testing.T) {
	tests := []struct {
		name             string
		workDir          string
		interfaceName    string
		interfacePkgPath string
		searchPattern    string
		wantFuncs        []string // 期待される関数名のリスト
		wantImplTypes    []string // 期待される実装型のリスト
		wantErr          bool
	}{
		{
			name:             "UserService インターフェースの参照を検索",
			workDir:          "../../sample/basic",
			interfaceName:    "UserService",
			interfacePkgPath: "github.com/rmocchy/convinient_wire/sample/basic/service",
			searchPattern:    "./...",
			wantFuncs:        []string{"NewUserService"},
			wantImplTypes:    []string{"userServiceImpl"},
			wantErr:          false,
		},
		{
			name:             "UserRepository インターフェースの参照を検索",
			workDir:          "../../sample/basic",
			interfaceName:    "UserRepository",
			interfacePkgPath: "github.com/rmocchy/convinient_wire/sample/basic/repository",
			searchPattern:    "./...",
			wantFuncs:        []string{"NewUserRepository"},
			wantImplTypes:    []string{"userRepositoryImpl"},
			wantErr:          false,
		},
		{
			name:             "存在しないインターフェースの検索",
			workDir:          "../../sample/basic",
			interfaceName:    "NonExistentInterface",
			interfacePkgPath: "github.com/rmocchy/convinient_wire/sample/basic/service",
			searchPattern:    "./...",
			wantFuncs:        []string{},
			wantImplTypes:    []string{},
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs, err := FindInterfaceReferences(tt.workDir, tt.interfaceName, tt.interfacePkgPath, tt.searchPattern)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindInterfaceReferences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// 関数名の検証
			gotFuncs := make(map[string]bool)
			for _, ref := range refs {
				gotFuncs[ref.FunctionName] = true
			}

			for _, wantFunc := range tt.wantFuncs {
				if !gotFuncs[wantFunc] {
					t.Errorf("期待される関数 %s が見つかりませんでした。取得した参照: %+v", wantFunc, refs)
				}
			}

			// 実装型の検証
			gotImplTypes := make(map[string]bool)
			for _, ref := range refs {
				gotImplTypes[ref.ImplementingType] = true
			}

			for _, wantType := range tt.wantImplTypes {
				if !gotImplTypes[wantType] {
					t.Errorf("期待される実装型 %s が見つかりませんでした。取得した参照: %+v", wantType, refs)
				}
			}

			// 見つかった参照の詳細を出力（デバッグ用）
			if len(refs) > 0 {
				t.Logf("見つかった参照:")
				for _, ref := range refs {
					t.Logf("  関数: %s (パッケージ: %s)", ref.FunctionName, ref.PackagePath)
					t.Logf("    実装型: %s (パッケージ: %s)", ref.ImplementingType, ref.ImplementingPkgPath)
				}
			}
		})
	}
}

func TestFindInterfaceReferences_SpecificPackage(t *testing.T) {
	// 特定のパッケージのみを検索
	refs, err := FindInterfaceReferences(
		"../../sample/basic",
		"UserService",
		"github.com/rmocchy/convinient_wire/sample/basic/service",
		"github.com/rmocchy/convinient_wire/sample/basic/service",
	)

	if err != nil {
		t.Fatalf("FindInterfaceReferences() error = %v", err)
	}

	if len(refs) == 0 {
		t.Error("期待される参照が見つかりませんでした")
		return
	}

	// NewUserService 関数が見つかることを確認
	found := false
	for _, ref := range refs {
		if ref.FunctionName == "NewUserService" {
			found = true
			if ref.ImplementingType != "userServiceImpl" {
				t.Errorf("実装型が期待と異なります。got = %s, want = userServiceImpl", ref.ImplementingType)
			}
			if ref.PackagePath != "github.com/rmocchy/convinient_wire/sample/basic/service" {
				t.Errorf("パッケージパスが期待と異なります。got = %s", ref.PackagePath)
			}
		}
	}

	if !found {
		t.Error("NewUserService 関数が見つかりませんでした")
	}
}

func TestFindInterfaceReferences_InvalidWorkDir(t *testing.T) {
	_, err := FindInterfaceReferences(
		"/invalid/path/that/does/not/exist",
		"UserService",
		"github.com/rmocchy/convinient_wire/sample/basic/service",
		"./...",
	)

	if err == nil {
		t.Error("無効な作業ディレクトリでエラーが発生しませんでした")
	}
}

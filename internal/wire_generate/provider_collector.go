package wiregenerate

import (
	"go/types"

	pipe "github.com/rmocchy/convinient_wire/internal/analyze"
)

// collectProviderSets はルート構造体からプロバイダーセットを収集する
func collectProviderSets(rootStructs []*pipe.StructNode, importMap map[string]bool) []ProviderSetData {
	providerSets := make([]ProviderSetData, 0, len(rootStructs))

	for _, root := range rootStructs {
		providerMap := make(map[string]bool)
		providers := []string{}

		// フィールドの依存関係から再帰的に initFunction を収集
		collectInitFunctionsFromFields(root.Fields, importMap, providerMap, &providers)

		if len(providers) > 0 {
			providerSets = append(providerSets, ProviderSetData{
				StructName: root.StructName,
				Providers:  providers,
			})
		}
	}

	return providerSets
}

// collectInitFunctionsFromFields は FieldNode から再帰的に initFunction を収集する
func collectInitFunctionsFromFields(fields []pipe.FieldNode, importMap map[string]bool, providerMap map[string]bool, providers *[]string) {
	for _, field := range fields {
		switch f := field.(type) {
		case *pipe.StructNode:
			collectInitFunctionsFromFields(f.Dependencies, importMap, providerMap, providers)
			addInitFunctions(f.InitFunctions, importMap, providerMap, providers)
		case *pipe.InterfaceNode:
			collectInitFunctionsFromFields(f.Dependencies, importMap, providerMap, providers)
			addInitFunctions(f.InitFunctions, importMap, providerMap, providers)
		}
	}
}

// addInitFunctions は initFunction のリストをプロバイダーリストに追加する
func addInitFunctions(initFuncs []*types.Func, importMap map[string]bool, providerMap map[string]bool, providers *[]string) {
	for _, initFunc := range initFuncs {
		addInitFunction(initFunc, importMap, providerMap, providers)
	}
}

// addInitFunction は initFunction をプロバイダーリストに追加する（重複チェック付き）
func addInitFunction(initFunc *types.Func, importMap map[string]bool, providerMap map[string]bool, providers *[]string) {
	pkg := initFunc.Pkg()
	if pkg == nil {
		return
	}

	pkgPath := pkg.Path()
	pkgName := pkg.Name()
	fullName := pkgName + "." + initFunc.Name()

	if providerMap[fullName] {
		return
	}

	importMap[pkgPath] = true
	providerMap[fullName] = true
	*providers = append(*providers, fullName)
}

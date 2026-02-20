package analyze

import "go/types"

type AnalysisCache interface {
	Get(namedType *types.Named) ([]*FnDITreeNode, bool)
	Set(namedType *types.Named, functions []*FnDITreeNode)
}

type analysisCache struct {
	cache map[string][]*FnDITreeNode
}

func NewAnalysisCache() AnalysisCache {
	return &analysisCache{
		cache: make(map[string][]*FnDITreeNode),
	}
}

func (ac *analysisCache) Get(namedType *types.Named) ([]*FnDITreeNode, bool) {
	key := getIdenticalTypeName(namedType)
	functions, found := ac.cache[key]
	return functions, found
}

func (ac *analysisCache) Set(namedType *types.Named, functions []*FnDITreeNode) {
	key := getIdenticalTypeName(namedType)
	ac.cache[key] = functions
}

// pkgPath + defName
func getIdenticalTypeName(t types.Type) string {
	switch tt := t.(type) {
	case *types.Named:
		obj := tt.Obj()
		return obj.Pkg().Path() + "." + obj.Name()
	default:
		return ""
	}
}

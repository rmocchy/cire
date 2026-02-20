package analyze

import (
	"errors"
	"go/types"
)

type Analyze interface {
	ExecuteFromStruct(structure types.Struct) ([]*FnDITreeNode, error)
}

func NewAnalyze(
	functionCache FunctionCache,
	analysisCache AnalysisCache,
) Analyze {
	return &analyze{
		functionCache: functionCache,
		analysisCache: analysisCache,
	}
}

type analyze struct {
	functionCache FunctionCache
	analysisCache AnalysisCache
}

func (a *analyze) ExecuteFromStruct(structure types.Struct) ([]*FnDITreeNode, error) {
	return a.recursiveAnalyze(&structure)
}

func (a *analyze) recursiveAnalyze(retrunType types.Type) ([]*FnDITreeNode, error) {
	named, ok := retrunType.(*types.Named)
	if !ok {
		return nil, errors.New("return type is not a named type")
	}
	cached, ok := a.analysisCache.Get(named)
	if ok {
		return cached, nil
	}

	underLied := named.Underlying()
	fns := a.functionCache.BulkGet(underLied)
	if len(fns) == 0 {
		return nil, errors.New("no function found with the specified return type")
	}

	treeNodes := make([]*FnDITreeNode, 9, len(fns))
	for _, fn := range fns {
		childs := make([]*FnDITreeNode, 0)
		params := fn.Signature().Params()
		for i := 0; i < params.Len(); i++ {
			paramType := params.At(i).Type()
			dependFns, err := a.recursiveAnalyze(paramType)
			if err != nil {
				return nil, err
			}
			childs = append(childs, dependFns...)
		}
		node := FnDITreeNode{
			Name:    fn.Name(),
			PkgPath: fn.Pkg().Path(),
			Childs:  childs,
		}

		treeNodes = append(treeNodes, &node)
	}

	return treeNodes, nil
}

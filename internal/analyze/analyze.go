package analyze

import (
	"errors"
	"go/types"
)

type Analyze interface {
	ExecuteFromStruct(structure *types.Named) ([]*FnDITreeNode, error)
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

func (a *analyze) ExecuteFromStruct(structure *types.Named) ([]*FnDITreeNode, error) {
	st, ok := structure.Underlying().(*types.Struct)
	if !ok {
		return nil, errors.New("not a struct type")
	}
	var allNodes []*FnDITreeNode
	for i := 0; i < st.NumFields(); i++ {
		fieldType, ok := Deref(st.Field(i).Type()).(*types.Named)
		if !ok {
			continue
		}
		nodes, err := a.recursiveAnalyze(fieldType)
		if err != nil {
			return nil, err
		}
		allNodes = append(allNodes, nodes...)
	}
	return allNodes, nil
}

func (a *analyze) recursiveAnalyze(retrunType *types.Named) ([]*FnDITreeNode, error) {
	cached, ok := a.analysisCache.Get(retrunType)
	if ok {
		return cached, nil
	}
	fns := a.functionCache.BulkGet(retrunType)
	if len(fns) == 0 {
		return nil, errors.New("no function found with the specified return type")
	}

	treeNodes := make([]*FnDITreeNode, 0, len(fns))
	for _, fn := range fns {
		childs := make([]*FnDITreeNode, 0)
		params := fn.Signature().Params()
		for i := 0; i < params.Len(); i++ {
			paramType := Deref(params.At(i).Type())
			named, ok := paramType.(*types.Named)
			if !ok {
				continue
			}
			dependFns, err := a.recursiveAnalyze(named)
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

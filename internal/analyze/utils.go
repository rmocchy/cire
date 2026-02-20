package analyze

import "go/types"

type ConvertTreeToUniqueList struct {
	Visited map[string]bool
	List    []*FnDITreeNode
}

func NewConvertTreeToUniqueList() *ConvertTreeToUniqueList {
	return &ConvertTreeToUniqueList{
		Visited: make(map[string]bool),
		List:    []*FnDITreeNode{},
	}
}

func (c *ConvertTreeToUniqueList) Execute(node *FnDITreeNode) {
	key := node.PkgPath + "." + node.Name
	if c.Visited[key] {
		return
	}
	c.Visited[key] = true
	c.List = append(c.List, node)
	for _, child := range node.Childs {
		c.Execute(child)
	}
}

// Deref は、ポインタ型の場合はその要素の型を返し、そうでない場合はそのままの型を返す
func Deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}

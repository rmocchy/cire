package analyze

import "go/types"

type convertTreeToUniqueList struct {
	visited map[string]bool
	list    []*FnDITreeNode
}

func NewConvertTreeToUniqueList() *convertTreeToUniqueList {
	return &convertTreeToUniqueList{
		visited: make(map[string]bool),
		list:    []*FnDITreeNode{},
	}
}

func (c *convertTreeToUniqueList) Execute(node *FnDITreeNode) {
	key := node.PkgPath + "." + node.Name
	if c.visited[key] {
		return
	}
	c.visited[key] = true
	c.list = append(c.list, node)
	for _, child := range node.Childs {
		c.Execute(child)
	}
}

func (c *convertTreeToUniqueList) List() []*FnDITreeNode {
	return c.list
}

// Deref は、ポインタ型の場合はその要素の型を返し、そうでない場合はそのままの型を返す
func Deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}

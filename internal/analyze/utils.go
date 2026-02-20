package analyze

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

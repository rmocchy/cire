package analyze

type FnDITreeNode struct {
	Name    string          `json:"name"`
	PkgPath string          `json:"pkg_path"`
	Childs  []*FnDITreeNode `json:"childs"`
}

package analyze

type Function struct {
	Name        string `json:"name"`
	PackagePath string `json:"package_path"`
}

type FnDITreeNode struct {
	Name    string          `json:"name"`
	PkgPath string          `json:"pkg_path"`
	Childs  []*FnDITreeNode `json:"childs"`
}

package core

type PackagePath struct {
	path string
}

func NewPackagePath(path string) PackagePath {
	return PackagePath{path: path}
}

func (p PackagePath) String() string {
	return p.path
}

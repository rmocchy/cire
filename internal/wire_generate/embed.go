package wiregenerate

import (
	_ "embed"
)

//go:embed wire.go.tmpl
var wireTemplate string

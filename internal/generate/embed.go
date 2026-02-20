package generate

import (
	_ "embed"
)

//go:embed wire.go.tmpl
var wireTemplate string

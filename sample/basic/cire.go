package main

import (
	"github.com/rmocchy/convinient_wire/sample/basic/handler"
)

// App は依存関係の解析対象となるルート構造体
type App struct {
	handler *handler.UserHandler
}

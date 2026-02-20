package main

import (
	"github.com/rmocchy/cire/sample/duplicate/handler"
)

// App は依存関係の解析対象となるルート構造体
type App struct {
	handler *handler.UserHandler
}

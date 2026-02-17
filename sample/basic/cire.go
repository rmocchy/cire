//go:build cire
// +build cire

package main

import (
	"github.com/rmocchy/convinient_wire/sample/basic/handler"
)

// ControllerSet は依存関係の解析対象となるルート構造体
type ControllerSet struct {
	handler *handler.UserHandler
}

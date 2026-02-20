package analyze

import (
	"fmt"
)

// 特定の構造体が複数の異なる関数によって生成されている場合はエラーにする
// 構造体名のマップを埋めていき, 全てが過分なく埋まらなければエラーにする
// 平滑化された前提なので深さ1のみ探索する
func IsDepTreeSatisfiable(nodes []*FnDITreeNode) error {
	retToFnName := make(map[string]string)
	for _, node := range nodes {
		if node == nil {
			continue
		}
		for _, ret := range node.ReturnTypes {
			if existingFn, exists := retToFnName[ret]; exists {
				if existingFn != node.Name {
					return fmt.Errorf("multiple functions found for return type %s: %s and %s", ret, existingFn, node.Name)
				}
			} else {
				retToFnName[ret] = node.Name
			}
		}
	}
	return nil
}

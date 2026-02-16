package app

import (
	"fmt"

	file "github.com/rmocchy/convinient_wire/ast_analyzer/files"
	"github.com/rmocchy/convinient_wire/internal/pipe"
)

// WireAnalyzer はwire.goの解析を行う
type WireAnalyzer struct {
	analyzer *pipe.WireAnalyzer
}

// NewWireAnalyzer は新しいWireAnalyzerを作成する
func NewWireAnalyzer(workDir, searchPattern string) (*WireAnalyzer, error) {
	analyzer, err := pipe.NewWireAnalyzer(workDir, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create analyzer: %w", err)
	}

	return &WireAnalyzer{
		analyzer: analyzer,
	}, nil
}

// AnalyzeWireFile はwire.goファイルを解析する
func (wa *WireAnalyzer) AnalyzeWireFile(wireFilePath string) ([]*StructNode, error) {
	// wire.goから構造体を取得
	functions, err := file.ParseWireFileStructs(wireFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wire file: %w", err)
	}

	var results []*StructNode

	// 各関数の返り値構造体を解析
	for _, funcInfo := range functions {
		for _, structInfo := range funcInfo.ReturnTypes {
			// 構造体を再帰的に解析
			pipeNode, err := wa.analyzer.AnalyzeStruct(structInfo.PackagePath, structInfo.Name)
			if err != nil {
				// エラーがあっても他の構造体の解析を続ける
				results = append(results, &StructNode{
					StructName: structInfo.Name,
					Skipped:    true,
					SkipReason: fmt.Sprintf("failed to analyze: %v", err),
				})
				continue
			}
			// pipe.StructNodeからapp.StructNodeに変換
			results = append(results, convertPipeNodeToAppNode(pipeNode))
		}
	}

	return results, nil
}

// convertPipeNodeToAppNode はpipe.StructNodeをapp.StructNodeに変換する
func convertPipeNodeToAppNode(pipeNode *pipe.StructNode) *StructNode {
	if pipeNode == nil {
		return nil
	}

	appNode := &StructNode{
		FieldName:     pipeNode.FieldName,
		StructName:    pipeNode.StructName,
		PackagePath:   pipeNode.PackagePath,
		InitFunctions: make([]InitFunctionInfo, 0, len(pipeNode.InitFunctions)),
		Fields:        make([]FieldNode, 0, len(pipeNode.Fields)),
		Skipped:       pipeNode.Skipped,
		SkipReason:    pipeNode.SkipReason,
	}

	// InitFunctionsを変換
	for _, fn := range pipeNode.InitFunctions {
		appNode.InitFunctions = append(appNode.InitFunctions, InitFunctionInfo{
			Name:        fn.Name,
			PackagePath: fn.PackagePath,
		})
	}

	// Fieldsを変換
	for _, field := range pipeNode.Fields {
		appNode.Fields = append(appNode.Fields, convertPipeFieldToAppField(field))
	}

	return appNode
}

// convertPipeFieldToAppField はpipe.FieldNodeをapp.FieldNodeに変換する
func convertPipeFieldToAppField(pipeField pipe.FieldNode) FieldNode {
	if pipeField == nil {
		return nil
	}

	switch pipeField.NodeType() {
	case pipe.NodeTypeStruct:
		if pipeStruct, ok := pipeField.(*pipe.StructNode); ok {
			return convertPipeNodeToAppNode(pipeStruct)
		}
	case pipe.NodeTypeInterface:
		if pipeInterface, ok := pipeField.(*pipe.InterfaceNode); ok {
			return &InterfaceNode{
				FieldName:      pipeInterface.FieldName,
				TypeName:       pipeInterface.TypeName,
				PackagePath:    pipeInterface.PackagePath,
				ResolvedStruct: convertPipeNodeToAppNode(pipeInterface.ResolvedStruct),
				Skipped:        pipeInterface.Skipped,
				SkipReason:     pipeInterface.SkipReason,
			}
		}
	}

	return nil
}

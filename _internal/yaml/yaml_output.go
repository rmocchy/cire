package yaml

import (
	"fmt"
	"go/types"
	"os"

	"github.com/rmocchy/cire/internal/analyze"
	"gopkg.in/yaml.v3"
)

// YAMLRootOutput はルート出力構造体（複数の@cire構造体に対応）
type YAMLRootOutput struct {
	Root []YAMLOutput `yaml:"root"`
}

// YAMLOutput は構造体の依存関係をYAML形式で出力するための構造体
type YAMLOutput struct {
	StructName   string          `yaml:"struct_name"`
	PackagePath  string          `yaml:"package_path"`
	Functions    []YAMLFunction  `yaml:"init_functions,omitempty"`
	Dependencies []YAMLFieldNode `yaml:"dependencies,omitempty"`
	Fields       []YAMLFieldNode `yaml:"fields,omitempty"`
	Skipped      bool            `yaml:"skipped,omitempty"`
	SkipReason   string          `yaml:"skip_reason,omitempty"`
}

// YAMLFunction は関数情報をYAML形式で出力するための構造体
type YAMLFunction struct {
	Name        string `yaml:"name"`
	PackagePath string `yaml:"package_path"`
}

// YAMLFieldNode はフィールドをYAML形式で出力するための構造体
type YAMLFieldNode struct {
	FieldName    string          `yaml:"field_name"`
	Type         string          `yaml:"type"`
	NodeType     string          `yaml:"node_type"`
	PackagePath  string          `yaml:"package_path,omitempty"`
	Functions    []YAMLFunction  `yaml:"init_functions,omitempty"`
	Dependencies []YAMLFieldNode `yaml:"dependencies,omitempty"`
	Fields       []YAMLFieldNode `yaml:"fields,omitempty"`
	Skipped      bool            `yaml:"skipped,omitempty"`
	SkipReason   string          `yaml:"skip_reason,omitempty"`
}

// ToYAMLOutput はStructNodeをYAMLOutput形式に変換する
func ToYAMLOutput(node *analyze.StructNode) YAMLOutput {
	return YAMLOutput{
		StructName:   node.StructName,
		PackagePath:  node.PackagePath,
		Functions:    convertFunctions(node.InitFunctions),
		Fields:       convertFields(node.Fields),
		Dependencies: convertFields(node.Dependencies),
		Skipped:      node.Skipped,
		SkipReason:   node.SkipReason,
	}
}

// convertFunctions は関数のリストをYAMLFunction形式に変換する
func convertFunctions(fns []*types.Func) []YAMLFunction {
	if len(fns) == 0 {
		return nil
	}

	result := make([]YAMLFunction, 0, len(fns))
	for _, fn := range fns {
		result = append(result, YAMLFunction{
			Name:        fn.Name(),
			PackagePath: fn.Pkg().Path(),
		})
	}
	return result
}

// convertFields はフィールドのリストをYAMLFieldNode形式に変換する
func convertFields(fields []analyze.FieldNode) []YAMLFieldNode {
	if len(fields) == 0 {
		return nil
	}

	result := make([]YAMLFieldNode, 0, len(fields))
	for _, field := range fields {
		switch f := field.(type) {
		case *analyze.StructNode:
			result = append(result, YAMLFieldNode{
				FieldName:    f.FieldName,
				Type:         f.StructName,
				NodeType:     "struct",
				PackagePath:  f.PackagePath,
				Functions:    convertFunctions(f.InitFunctions),
				Dependencies: convertFields(f.Dependencies),
				Fields:       convertFields(f.Fields),
				Skipped:      f.Skipped,
				SkipReason:   f.SkipReason,
			})
		case *analyze.InterfaceNode:
			result = append(result, YAMLFieldNode{
				FieldName:    f.FieldName,
				Type:         f.TypeName,
				NodeType:     "interface",
				PackagePath:  f.PackagePath,
				Functions:    convertFunctions(f.InitFunctions),
				Dependencies: convertFields(f.Dependencies),
				Skipped:      f.Skipped,
				SkipReason:   f.SkipReason,
			})
		case *analyze.BuiltinNode:
			result = append(result, YAMLFieldNode{
				FieldName: f.FieldName,
				Type:      f.TypeName,
				NodeType:  "builtin",
			})
		}
	}
	return result
}

// OutputToYAML は解析結果をYAML形式で出力する
// outputFile が空文字列の場合は標準出力に出力
func OutputToYAML(node *analyze.StructNode, outputFile string) error {
	// YAMLに変換（root配列形式）
	yamlOutput := ToYAMLOutput(node)
	rootOutput := YAMLRootOutput{
		Root: []YAMLOutput{yamlOutput},
	}
	yamlData, err := yaml.Marshal(&rootOutput)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	// 出力
	if outputFile != "" {
		if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Successfully wrote dependencies to %s\n", outputFile)
	} else {
		fmt.Print(string(yamlData))
	}

	return nil
}

// OutputMultipleToYAML は複数の解析結果をYAML形式で出力する
// outputFile が空文字列の場合は標準出力に出力
func OutputMultipleToYAML(nodes []*analyze.StructNode, outputFile string) error {
	if len(nodes) == 0 {
		return fmt.Errorf("no nodes to output")
	}

	// すべてのノードをYAMLに変換
	yamlOutputs := make([]YAMLOutput, 0, len(nodes))
	for _, node := range nodes {
		yamlOutputs = append(yamlOutputs, ToYAMLOutput(node))
	}

	rootOutput := YAMLRootOutput{
		Root: yamlOutputs,
	}
	yamlData, err := yaml.Marshal(&rootOutput)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	// 出力
	if outputFile != "" {
		if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Successfully wrote dependencies to %s\n", outputFile)
	} else {
		fmt.Print(string(yamlData))
	}

	return nil
}

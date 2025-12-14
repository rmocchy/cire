package app

import (
	"fmt"
	"testing"
)

func TestWireAnalyzer_AnalyzeWireFile(t *testing.T) {
	// サンプルディレクトリのwire.goを解析
	workDir := "../../sample/basic"
	wireFilePath := "../../sample/basic/wire.go"
	searchPattern := "./..."

	analyzer := NewWireAnalyzer(workDir, searchPattern)
	results, err := analyzer.AnalyzeWireFile(wireFilePath)
	if err != nil {
		t.Fatalf("AnalyzeWireFile failed: %v", err)
	}

	// 結果を表示
	for _, result := range results {
		printStructAnalysis(t, &result, 0)
	}
}

// printStructAnalysis は構造体の解析結果を階層的に表示する
func printStructAnalysis(t *testing.T, result *StructAnalysisResult, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	if result.Skipped {
		t.Logf("%s[SKIPPED] %s: %s", prefix, result.StructName, result.SkipReason)
		return
	}

	t.Logf("%sStruct: %s (Package: %s)", prefix, result.StructName, result.PackagePath)

	for _, field := range result.Fields {
		fieldPrefix := prefix + "  "
		pointer := ""
		if field.IsPointer {
			pointer = "*"
		}

		if field.IsInterface {
			t.Logf("%sField: %s %s%s (interface, Package: %s)",
				fieldPrefix, field.Name, pointer, field.TypeName, field.PackagePath)

			if field.InterfaceSkipped {
				t.Logf("%s  -> [SKIPPED] %s", fieldPrefix, field.InterfaceSkipReason)
			} else if field.ResolvedStruct != nil {
				t.Logf("%s  -> Resolved to:", fieldPrefix)
				printStructAnalysis(t, field.ResolvedStruct, indent+3)
			}
		} else if field.ResolvedStruct != nil {
			t.Logf("%sField: %s %s%s (Package: %s)",
				fieldPrefix, field.Name, pointer, field.TypeName, field.PackagePath)
			printStructAnalysis(t, field.ResolvedStruct, indent+2)
		} else {
			t.Logf("%sField: %s %s%s (Package: %s)",
				fieldPrefix, field.Name, pointer, field.TypeName, field.PackagePath)
		}
	}
}

func TestWireAnalyzer_AnalyzeStruct(t *testing.T) {
	workDir := "../../sample/basic"
	searchPattern := "./..."

	analyzer := NewWireAnalyzer(workDir, searchPattern)

	// ControllerSetを直接解析
	result, err := analyzer.analyzeStruct("", "ControllerSet")
	if err != nil {
		t.Fatalf("analyzeStruct failed: %v", err)
	}

	printStructAnalysis(t, result, 0)
}

func ExampleWireAnalyzer_AnalyzeWireFile() {
	workDir := "../../sample/basic"
	wireFilePath := "../../sample/basic/wire.go"
	searchPattern := "./..."

	analyzer := NewWireAnalyzer(workDir, searchPattern)
	results, err := analyzer.AnalyzeWireFile(wireFilePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, result := range results {
		printStructAnalysisExample(&result, 0)
	}
}

func printStructAnalysisExample(result *StructAnalysisResult, indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	if result.Skipped {
		fmt.Printf("%s[SKIPPED] %s: %s\n", prefix, result.StructName, result.SkipReason)
		return
	}

	fmt.Printf("%sStruct: %s\n", prefix, result.StructName)

	for _, field := range result.Fields {
		fieldPrefix := prefix + "  "
		pointer := ""
		if field.IsPointer {
			pointer = "*"
		}

		if field.IsInterface {
			fmt.Printf("%sField: %s %s%s (interface)\n",
				fieldPrefix, field.Name, pointer, field.TypeName)

			if field.InterfaceSkipped {
				fmt.Printf("%s  -> [SKIPPED] %s\n", fieldPrefix, field.InterfaceSkipReason)
			} else if field.ResolvedStruct != nil {
				fmt.Printf("%s  -> Resolved to:\n", fieldPrefix)
				printStructAnalysisExample(field.ResolvedStruct, indent+3)
			}
		} else if field.ResolvedStruct != nil {
			fmt.Printf("%sField: %s %s%s\n",
				fieldPrefix, field.Name, pointer, field.TypeName)
			printStructAnalysisExample(field.ResolvedStruct, indent+2)
		} else {
			fmt.Printf("%sField: %s %s%s\n",
				fieldPrefix, field.Name, pointer, field.TypeName)
		}
	}
}

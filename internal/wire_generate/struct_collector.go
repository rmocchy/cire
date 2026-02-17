package wiregenerate

import (
	"path/filepath"

	pipe "github.com/rmocchy/convinient_wire/internal/analyze"
)

// collectStructDefs は StructNode から構造体定義を収集する
func collectStructDefs(results []*pipe.StructNode, importMap map[string]bool) []StructDef {
	structs := make([]StructDef, 0, len(results))

	for _, result := range results {
		structDef := StructDef{
			Name:   result.StructName,
			Fields: convertFieldsToStructFieldDefs(result.Fields, importMap),
		}
		structs = append(structs, structDef)
	}

	return structs
}

// convertFieldsToStructFieldDefs はフィールドリストを StructFieldDef のリストに変換する
func convertFieldsToStructFieldDefs(fields []pipe.FieldNode, importMap map[string]bool) []StructFieldDef {
	fieldDefs := make([]StructFieldDef, 0, len(fields))

	for _, field := range fields {
		if fieldDef := convertFieldToStructFieldDef(field, importMap); fieldDef != nil {
			fieldDefs = append(fieldDefs, *fieldDef)
		}
	}

	return fieldDefs
}

// convertFieldToStructFieldDef は FieldNode を StructFieldDef に変換する
func convertFieldToStructFieldDef(field pipe.FieldNode, importMap map[string]bool) *StructFieldDef {
	switch f := field.(type) {
	case *pipe.StructNode:
		if f.PackagePath != "" {
			importMap[f.PackagePath] = true
		}
		pkgName := filepath.Base(f.PackagePath)
		return &StructFieldDef{
			Name:    f.FieldName,
			Type:    pkgName + "." + f.StructName,
			Pointer: true,
		}
	case *pipe.InterfaceNode:
		if f.PackagePath != "" {
			importMap[f.PackagePath] = true
		}
		pkgName := filepath.Base(f.PackagePath)
		return &StructFieldDef{
			Name:    f.FieldName,
			Type:    pkgName + "." + f.TypeName,
			Pointer: false,
		}
	case *pipe.BuiltinNode:
		return &StructFieldDef{
			Name:    f.FieldName,
			Type:    f.TypeName,
			Pointer: false,
		}
	}
	return nil
}

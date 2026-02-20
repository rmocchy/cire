package analyze

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type JsonConfig struct {
	Dir  string
	Data map[string][]*FnDITreeNode
}

func WriteOnJsonFile(config *JsonConfig) error {
	fileName := "dep_tree.json"
	filePath := filepath.Join(config.Dir, fileName)

	data, err := json.MarshalIndent(config.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("JSON file generated: %s\n", filePath)
	return nil
}

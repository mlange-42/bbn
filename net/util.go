package net

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FromFile reads a [Network] from an YAML or XML file.
func FromFile(path string) (*Network, []Variable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	ext := filepath.Ext(path)

	switch strings.ToLower(ext) {
	case ".yml":
		n, err := FromYAML(data)
		if err != nil {
			return nil, nil, err
		}
		return n, n.variables, nil
	case ".xml", ".bifxml":
		n, err := FromBIFXML(data)
		if err != nil {
			return nil, nil, err
		}
		return n, n.variables, nil
	default:
		return nil, nil, fmt.Errorf("unsupported file format '%s'", ext)
	}
}

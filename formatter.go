package main

import (
	"encoding/json"
	"fmt"
)

// GenerateJSON creates a JSON representation of the struct layout
func GenerateJSON(layout *StructLayout, pretty bool) (string, error) {
	var data []byte
	var err error
	
	if pretty {
		data, err = json.MarshalIndent(layout, "", "  ")
	} else {
		data, err = json.Marshal(layout)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	return string(data), nil
}

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatASCII   OutputFormat = "ascii"
	FormatCompact OutputFormat = "compact"
	FormatJSON    OutputFormat = "json"
)

// FormatOutput formats the layout according to the specified format
func FormatOutput(layout *StructLayout, format OutputFormat) (string, error) {
	switch format {
	case FormatASCII:
		return GenerateASCIIDiagram(layout), nil
	case FormatCompact:
		return GenerateCompactDiagram(layout), nil
	case FormatJSON:
		return GenerateJSON(layout, true)
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

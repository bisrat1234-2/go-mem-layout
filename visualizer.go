package main

import (
	"fmt"
	"strings"
)

// GenerateASCIIDiagram creates an ASCII diagram of the struct layout
func GenerateASCIIDiagram(layout *StructLayout) string {
	var sb strings.Builder
	
	// Header
	sb.WriteString(fmt.Sprintf("Struct: %s\n", layout.Name))
	sb.WriteString(fmt.Sprintf("Total Size: %d bytes\n", layout.TotalSize))
	sb.WriteString(fmt.Sprintf("Alignment: %d bytes\n\n", layout.Alignment))
	
	if len(layout.Fields) == 0 {
		sb.WriteString("No fields\n")
		return sb.String()
	}
	
	// Calculate the width for the diagram
	maxNameLen := 0
	maxTypeLen := 0
	for _, field := range layout.Fields {
		if len(field.Name) > maxNameLen {
			maxNameLen = len(field.Name)
		}
		if len(field.Type) > maxTypeLen {
			maxTypeLen = len(field.Type)
		}
	}
	
	// Ensure minimum widths
	if maxNameLen < 10 {
		maxNameLen = 10
	}
	if maxTypeLen < 10 {
		maxTypeLen = 10
	}
	
	// Build the diagram
	sb.WriteString("Memory Layout:\n")
	sb.WriteString("┌────────┬─────────────────────────────────────────────────────────────┐\n")
	sb.WriteString("│ Offset │ Field                                                       │\n")
	sb.WriteString("├────────┼─────────────────────────────────────────────────────────────┤\n")
	
	// Track current offset for padding detection
	currentOffset := uintptr(0)
	paddingIndex := 0
	
	for i, field := range layout.Fields {
		// Check if there's padding before this field
		if paddingIndex < len(layout.Padding) {
			pad := layout.Padding[paddingIndex]
			if i > 0 && pad.Offset == currentOffset {
				// Print padding
				for j := uintptr(0); j < pad.Size; j++ {
					sb.WriteString(fmt.Sprintf("│ %6d │ [padding]                                                   │\n", 
						currentOffset+j))
				}
				currentOffset += pad.Size
				paddingIndex++
			}
		}
		
		// Print field
		sb.WriteString(fmt.Sprintf("│ %6d │ %-20s %-20s (%d bytes)%s│\n",
			field.Offset,
			field.Name,
			field.Type,
			field.Size,
			strings.Repeat(" ", max(0, 17-len(fmt.Sprintf("%d", field.Size))))))
		
		// Print continuation for multi-byte fields
		for j := uintptr(1); j < field.Size; j++ {
			sb.WriteString(fmt.Sprintf("│ %6d │ %s│\n", 
				field.Offset+j,
				strings.Repeat(" ", 60)))
		}
		
		currentOffset = field.Offset + field.Size
	}
	
	// Check for trailing padding
	if paddingIndex < len(layout.Padding) {
		pad := layout.Padding[paddingIndex]
		if pad.Offset == currentOffset {
			for j := uintptr(0); j < pad.Size; j++ {
				sb.WriteString(fmt.Sprintf("│ %6d │ [padding]                                                   │\n", 
					currentOffset+j))
			}
		}
	}
	
	sb.WriteString("└────────┴─────────────────────────────────────────────────────────────┘\n")
	
	// Summary
	sb.WriteString("\nField Details:\n")
	for _, field := range layout.Fields {
		sb.WriteString(fmt.Sprintf("  %-20s offset=%3d size=%3d align=%d type=%s\n",
			field.Name+":",
			field.Offset,
			field.Size,
			field.Alignment,
			field.Type))
	}
	
	if len(layout.Padding) > 0 {
		sb.WriteString("\nPadding:\n")
		totalPadding := uintptr(0)
		for _, pad := range layout.Padding {
			sb.WriteString(fmt.Sprintf("  %d bytes at offset %d (after %s)\n",
				pad.Size, pad.Offset, pad.After))
			totalPadding += pad.Size
		}
		sb.WriteString(fmt.Sprintf("  Total padding: %d bytes (%.1f%% waste)\n",
			totalPadding,
			float64(totalPadding)/float64(layout.TotalSize)*100.0))
	}
	
	return sb.String()
}

// GenerateCompactDiagram creates a more compact ASCII representation
func GenerateCompactDiagram(layout *StructLayout) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("struct %s {  // %d bytes, align %d\n",
		layout.Name, layout.TotalSize, layout.Alignment))
	
	paddingMap := make(map[uintptr]PaddingInfo)
	for _, pad := range layout.Padding {
		paddingMap[pad.Offset] = pad
	}
	
	for _, field := range layout.Fields {
		sb.WriteString(fmt.Sprintf("  %s %s  // offset %d, size %d, align %d\n",
			field.Name,
			field.Type,
			field.Offset,
			field.Size,
			field.Alignment))
		
		// Check for padding after this field
		nextOffset := field.Offset + field.Size
		if pad, ok := paddingMap[nextOffset]; ok {
			sb.WriteString(fmt.Sprintf("  // [%d bytes padding]\n", pad.Size))
		}
	}
	
	sb.WriteString("}\n")
	
	return sb.String()
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

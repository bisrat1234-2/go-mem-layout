package main

import (
	"reflect"
	"testing"
)

func TestAnalyzeStruct(t *testing.T) {
	type SimpleStruct struct {
		A int8
		B int64
		C int16
	}

	layout, err := AnalyzeStruct(reflect.TypeOf(SimpleStruct{}))
	if err != nil {
		t.Fatalf("AnalyzeStruct failed: %v", err)
	}

	if layout.Name != "SimpleStruct" {
		t.Errorf("Expected name 'SimpleStruct', got '%s'", layout.Name)
	}

	if layout.TotalSize != 24 {
		t.Errorf("Expected total size 24, got %d", layout.TotalSize)
	}

	if layout.Alignment != 8 {
		t.Errorf("Expected alignment 8, got %d", layout.Alignment)
	}

	if len(layout.Fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(layout.Fields))
	}

	// Check field A
	if layout.Fields[0].Name != "A" {
		t.Errorf("Expected field name 'A', got '%s'", layout.Fields[0].Name)
	}
	if layout.Fields[0].Offset != 0 {
		t.Errorf("Expected field A offset 0, got %d", layout.Fields[0].Offset)
	}
	if layout.Fields[0].Size != 1 {
		t.Errorf("Expected field A size 1, got %d", layout.Fields[0].Size)
	}

	// Check field B
	if layout.Fields[1].Name != "B" {
		t.Errorf("Expected field name 'B', got '%s'", layout.Fields[1].Name)
	}
	if layout.Fields[1].Offset != 8 {
		t.Errorf("Expected field B offset 8, got %d", layout.Fields[1].Offset)
	}
	if layout.Fields[1].Size != 8 {
		t.Errorf("Expected field B size 8, got %d", layout.Fields[1].Size)
	}

	// Check padding
	if len(layout.Padding) != 2 {
		t.Fatalf("Expected 2 padding entries, got %d", len(layout.Padding))
	}
}

func TestAnalyzeStructNoPadding(t *testing.T) {
	type NoPaddingStruct struct {
		A int64
		B int64
		C int64
	}

	layout, err := AnalyzeStruct(reflect.TypeOf(NoPaddingStruct{}))
	if err != nil {
		t.Fatalf("AnalyzeStruct failed: %v", err)
	}

	if len(layout.Padding) != 0 {
		t.Errorf("Expected no padding, got %d padding entries", len(layout.Padding))
	}

	if layout.TotalSize != 24 {
		t.Errorf("Expected total size 24, got %d", layout.TotalSize)
	}
}

func TestAnalyzeStructOptimized(t *testing.T) {
	type OptimizedStruct struct {
		A int64
		B int64
		C int32
		D int16
		E int8
		F bool
	}

	layout, err := AnalyzeStruct(reflect.TypeOf(OptimizedStruct{}))
	if err != nil {
		t.Fatalf("AnalyzeStruct failed: %v", err)
	}

	// Should have no internal padding, only potential trailing
	if layout.TotalSize != 24 {
		t.Errorf("Expected total size 24, got %d", layout.TotalSize)
	}
}

func TestAnalyzeStructError(t *testing.T) {
	// Test with non-struct type
	_, err := AnalyzeStruct(reflect.TypeOf(42))
	if err == nil {
		t.Error("Expected error for non-struct type, got nil")
	}
}

func TestCalculatePadding(t *testing.T) {
	fields := []FieldLayout{
		{Name: "A", Offset: 0, Size: 1},
		{Name: "B", Offset: 8, Size: 8},
		{Name: "C", Offset: 16, Size: 2},
	}

	padding := calculatePadding(fields, 24)

	if len(padding) != 2 {
		t.Fatalf("Expected 2 padding entries, got %d", len(padding))
	}

	// First padding after A
	if padding[0].Offset != 1 {
		t.Errorf("Expected first padding offset 1, got %d", padding[0].Offset)
	}
	if padding[0].Size != 7 {
		t.Errorf("Expected first padding size 7, got %d", padding[0].Size)
	}
	if padding[0].After != "A" {
		t.Errorf("Expected first padding after 'A', got '%s'", padding[0].After)
	}

	// Trailing padding after C
	if padding[1].Offset != 18 {
		t.Errorf("Expected second padding offset 18, got %d", padding[1].Offset)
	}
	if padding[1].Size != 6 {
		t.Errorf("Expected second padding size 6, got %d", padding[1].Size)
	}
	if padding[1].After != "C" {
		t.Errorf("Expected second padding after 'C', got '%s'", padding[1].After)
	}
}

func TestGenerateASCIIDiagram(t *testing.T) {
	layout := &StructLayout{
		Name:      "TestStruct",
		TotalSize: 16,
		Alignment: 8,
		Fields: []FieldLayout{
			{Name: "A", Type: "int64", Offset: 0, Size: 8, Alignment: 8},
			{Name: "B", Type: "int64", Offset: 8, Size: 8, Alignment: 8},
		},
		Padding: []PaddingInfo{},
	}

	diagram := GenerateASCIIDiagram(layout)
	
	if diagram == "" {
		t.Error("Expected non-empty diagram")
	}

	// Check for key elements
	if !contains(diagram, "TestStruct") {
		t.Error("Diagram should contain struct name")
	}
	if !contains(diagram, "Total Size: 16") {
		t.Error("Diagram should contain total size")
	}
}

func TestGenerateCompactDiagram(t *testing.T) {
	layout := &StructLayout{
		Name:      "TestStruct",
		TotalSize: 16,
		Alignment: 8,
		Fields: []FieldLayout{
			{Name: "A", Type: "int64", Offset: 0, Size: 8, Alignment: 8},
			{Name: "B", Type: "int64", Offset: 8, Size: 8, Alignment: 8},
		},
		Padding: []PaddingInfo{},
	}

	diagram := GenerateCompactDiagram(layout)
	
	if diagram == "" {
		t.Error("Expected non-empty diagram")
	}

	if !contains(diagram, "struct TestStruct") {
		t.Error("Diagram should contain struct declaration")
	}
}

func TestGenerateJSON(t *testing.T) {
	layout := &StructLayout{
		Name:      "TestStruct",
		TotalSize: 16,
		Alignment: 8,
		Fields: []FieldLayout{
			{Name: "A", Type: "int64", Offset: 0, Size: 8, Alignment: 8},
		},
		Padding: []PaddingInfo{},
	}

	json, err := GenerateJSON(layout, true)
	if err != nil {
		t.Fatalf("GenerateJSON failed: %v", err)
	}

	if json == "" {
		t.Error("Expected non-empty JSON")
	}

	if !contains(json, "TestStruct") {
		t.Error("JSON should contain struct name")
	}
	if !contains(json, "\"total_size\": 16") {
		t.Error("JSON should contain total size")
	}
}

func TestFormatOutput(t *testing.T) {
	layout := &StructLayout{
		Name:      "TestStruct",
		TotalSize: 8,
		Alignment: 8,
		Fields:    []FieldLayout{{Name: "A", Type: "int64", Offset: 0, Size: 8, Alignment: 8}},
		Padding:   []PaddingInfo{},
	}

	tests := []struct {
		format OutputFormat
		want   string
	}{
		{FormatASCII, "TestStruct"},
		{FormatCompact, "struct TestStruct"},
		{FormatJSON, "\"name\": \"TestStruct\""},
	}

	for _, tt := range tests {
		output, err := FormatOutput(layout, tt.format)
		if err != nil {
			t.Errorf("FormatOutput(%v) failed: %v", tt.format, err)
			continue
		}
		if !contains(output, tt.want) {
			t.Errorf("FormatOutput(%v) output should contain %q", tt.format, tt.want)
		}
	}
}

func TestAnalyzeStructFromSource(t *testing.T) {
	source := `package main

type Example struct {
	A int8
	B int64
	C int16
}
`

	layout, err := AnalyzeStructFromSource(source, "Example")
	if err != nil {
		t.Fatalf("AnalyzeStructFromSource failed: %v", err)
	}

	if layout.Name != "Example" {
		t.Errorf("Expected name 'Example', got '%s'", layout.Name)
	}

	if len(layout.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(layout.Fields))
	}
}

func TestAnalyzeStructFromSourceNotFound(t *testing.T) {
	source := `package main

type Example struct {
	A int
}
`

	_, err := AnalyzeStructFromSource(source, "NotFound")
	if err == nil {
		t.Error("Expected error for non-existent struct, got nil")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

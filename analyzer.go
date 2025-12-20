package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"unsafe"
)

// FieldLayout represents the memory layout of a single struct field
type FieldLayout struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Offset    uintptr `json:"offset"`
	Size      uintptr `json:"size"`
	Alignment uintptr `json:"alignment"`
}

// StructLayout represents the complete memory layout of a struct
type StructLayout struct {
	Name       string        `json:"name"`
	TotalSize  uintptr       `json:"total_size"`
	Alignment  uintptr       `json:"alignment"`
	Fields     []FieldLayout `json:"fields"`
	Padding    []PaddingInfo `json:"padding"`
}

// PaddingInfo represents padding bytes in the struct
type PaddingInfo struct {
	Offset uintptr `json:"offset"`
	Size   uintptr `json:"size"`
	After  string  `json:"after"` // Field name after which padding occurs
}

// AnalyzeStruct analyzes a Go struct type and returns its memory layout
func AnalyzeStruct(structType reflect.Type) (*StructLayout, error) {
	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct type, got %v", structType.Kind())
	}

	layout := &StructLayout{
		Name:      structType.Name(),
		TotalSize: structType.Size(),
		Alignment: uintptr(structType.Align()),
		Fields:    make([]FieldLayout, 0, structType.NumField()),
		Padding:   make([]PaddingInfo, 0),
	}

	// Analyze each field
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		fieldLayout := FieldLayout{
			Name:      field.Name,
			Type:      field.Type.String(),
			Offset:    field.Offset,
			Size:      field.Type.Size(),
			Alignment: uintptr(field.Type.Align()),
		}
		
		layout.Fields = append(layout.Fields, fieldLayout)
	}

	// Calculate padding
	layout.Padding = calculatePadding(layout.Fields, layout.TotalSize)

	return layout, nil
}

// calculatePadding identifies padding bytes in the struct
func calculatePadding(fields []FieldLayout, totalSize uintptr) []PaddingInfo {
	padding := make([]PaddingInfo, 0)
	
	for i := 0; i < len(fields); i++ {
		currentField := fields[i]
		nextOffset := currentField.Offset + currentField.Size
		
		// Check for padding between current and next field
		if i < len(fields)-1 {
			nextField := fields[i+1]
			if nextOffset < nextField.Offset {
				padding = append(padding, PaddingInfo{
					Offset: nextOffset,
					Size:   nextField.Offset - nextOffset,
					After:  currentField.Name,
				})
			}
		} else {
			// Check for trailing padding
			if nextOffset < totalSize {
				padding = append(padding, PaddingInfo{
					Offset: nextOffset,
					Size:   totalSize - nextOffset,
					After:  currentField.Name,
				})
			}
		}
	}
	
	return padding
}

// AnalyzeStructFromSource analyzes a struct from Go source code
func AnalyzeStructFromSource(sourceCode, structName string) (*StructLayout, error) {
	fset := token.NewFileSet()
	
	// Parse the source code
	file, err := parser.ParseFile(fset, "", sourceCode, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source: %w", err)
	}

	// Find the struct declaration
	var structType *ast.StructType
	var foundName string
	
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == structName {
				if st, ok := x.Type.(*ast.StructType); ok {
					structType = st
					foundName = x.Name.Name
					return false
				}
			}
		}
		return true
	})

	if structType == nil {
		return nil, fmt.Errorf("struct %s not found in source", structName)
	}

	// Create type checker configuration
	conf := types.Config{Importer: nil}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	// Type check the file
	pkg, err := conf.Check("main", fset, []*ast.File{file}, info)
	if err != nil {
		// Continue even with type errors, as we can still extract basic structure
	}

	// Build layout from AST
	layout := &StructLayout{
		Name:   foundName,
		Fields: make([]FieldLayout, 0),
	}

	// Get the type from the package
	if pkg != nil {
		obj := pkg.Scope().Lookup(structName)
		if obj != nil {
			if named, ok := obj.Type().(*types.Named); ok {
				if st, ok := named.Underlying().(*types.Struct); ok {
					return analyzeStructFromTypes(foundName, st)
				}
			}
		}
	}

	return layout, nil
}

// analyzeStructFromTypes analyzes a types.Struct and calculates layout
func analyzeStructFromTypes(name string, st *types.Struct) (*StructLayout, error) {
	layout := &StructLayout{
		Name:   name,
		Fields: make([]FieldLayout, 0, st.NumFields()),
	}

	var offset uintptr
	var maxAlign uintptr = 1

	// Calculate offsets and sizes for each field
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		fieldType := field.Type()
		
		// Get size and alignment
		size := getSizeOfType(fieldType)
		align := getAlignOfType(fieldType)
		
		if align > maxAlign {
			maxAlign = align
		}
		
		// Align the offset
		if offset%align != 0 {
			offset = (offset + align - 1) / align * align
		}
		
		fieldLayout := FieldLayout{
			Name:      field.Name(),
			Type:      fieldType.String(),
			Offset:    offset,
			Size:      size,
			Alignment: align,
		}
		
		layout.Fields = append(layout.Fields, fieldLayout)
		offset += size
	}

	// Align total size to struct alignment
	if offset%maxAlign != 0 {
		offset = (offset + maxAlign - 1) / maxAlign * maxAlign
	}
	
	layout.TotalSize = offset
	layout.Alignment = maxAlign
	layout.Padding = calculatePadding(layout.Fields, layout.TotalSize)

	return layout, nil
}

// getSizeOfType returns the size of a types.Type
func getSizeOfType(t types.Type) uintptr {
	switch t := t.(type) {
	case *types.Basic:
		info := t.Info()
		switch {
		case info&types.IsBoolean != 0:
			return 1
		case info&types.IsInteger != 0:
			switch t.Kind() {
			case types.Int8, types.Uint8:
				return 1
			case types.Int16, types.Uint16:
				return 2
			case types.Int32, types.Uint32:
				return 4
			case types.Int64, types.Uint64:
				return 8
			case types.Int, types.Uint:
				return unsafe.Sizeof(int(0))
			case types.Uintptr:
				return unsafe.Sizeof(uintptr(0))
			}
		case info&types.IsFloat != 0:
			switch t.Kind() {
			case types.Float32:
				return 4
			case types.Float64:
				return 8
			}
		case info&types.IsComplex != 0:
			switch t.Kind() {
			case types.Complex64:
				return 8
			case types.Complex128:
				return 16
			}
		case info&types.IsString != 0:
			return unsafe.Sizeof("")
		}
	case *types.Pointer, *types.Chan, *types.Map, *types.Signature:
		return unsafe.Sizeof(uintptr(0))
	case *types.Slice:
		return unsafe.Sizeof([]int{})
	case *types.Array:
		elemSize := getSizeOfType(t.Elem())
		return uintptr(t.Len()) * elemSize
	case *types.Struct:
		// Recursively calculate struct size
		var offset uintptr
		var maxAlign uintptr = 1
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			size := getSizeOfType(field.Type())
			align := getAlignOfType(field.Type())
			if align > maxAlign {
				maxAlign = align
			}
			if offset%align != 0 {
				offset = (offset + align - 1) / align * align
			}
			offset += size
		}
		if offset%maxAlign != 0 {
			offset = (offset + maxAlign - 1) / maxAlign * maxAlign
		}
		return offset
	case *types.Interface:
		return unsafe.Sizeof((*interface{})(nil))
	}
	return unsafe.Sizeof(uintptr(0))
}

// getAlignOfType returns the alignment of a types.Type
func getAlignOfType(t types.Type) uintptr {
	switch t := t.(type) {
	case *types.Basic:
		info := t.Info()
		switch {
		case info&types.IsBoolean != 0:
			return 1
		case info&types.IsInteger != 0:
			switch t.Kind() {
			case types.Int8, types.Uint8:
				return 1
			case types.Int16, types.Uint16:
				return 2
			case types.Int32, types.Uint32:
				return 4
			case types.Int64, types.Uint64:
				return 8
			case types.Int, types.Uint:
				return uintptr(unsafe.Alignof(int(0)))
			case types.Uintptr:
				return uintptr(unsafe.Alignof(uintptr(0)))
			}
		case info&types.IsFloat != 0:
			switch t.Kind() {
			case types.Float32:
				return 4
			case types.Float64:
				return 8
			}
		case info&types.IsComplex != 0:
			switch t.Kind() {
			case types.Complex64:
				return 4
			case types.Complex128:
				return 8
			}
		case info&types.IsString != 0:
			return uintptr(unsafe.Alignof(""))
		}
	case *types.Pointer, *types.Chan, *types.Map, *types.Signature:
		return uintptr(unsafe.Alignof(uintptr(0)))
	case *types.Slice:
		return uintptr(unsafe.Alignof([]int{}))
	case *types.Array:
		return getAlignOfType(t.Elem())
	case *types.Struct:
		var maxAlign uintptr = 1
		for i := 0; i < t.NumFields(); i++ {
			align := getAlignOfType(t.Field(i).Type())
			if align > maxAlign {
				maxAlign = align
			}
		}
		return maxAlign
	case *types.Interface:
		return uintptr(unsafe.Alignof((*interface{})(nil)))
	}
	return uintptr(unsafe.Alignof(uintptr(0)))
}

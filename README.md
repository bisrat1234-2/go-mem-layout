# go-mem-layout

A Go tool that analyzes Go structs and visualizes memory alignment, padding, and offsets. Useful for performance engineering and systems research.

## Features

- **Struct Analysis**: Analyze memory layout of Go structs from source code or reflection
- **Multiple Output Formats**: ASCII diagrams, compact format, or JSON
- **Detailed Information**: Shows field offsets, sizes, alignment, and padding
- **Padding Detection**: Identifies and quantifies wasted memory due to alignment
- **Correctness**: Uses Go's reflection and type system for accurate analysis
- **Explainability**: Clear visual output with explanations of memory layout

## Installation

```bash
go install github.com/BaseMax/go-mem-layout@latest
```

Or clone and build:

```bash
git clone https://github.com/BaseMax/go-mem-layout.git
cd go-mem-layout
go build
```

## Usage

### Command Line Options

```
go-mem-layout [flags]

Flags:
  -format string
        Output format: ascii, compact, or json (default: ascii)
  -example string
        Run with example struct: basic, mixed, optimized, or all
  -file string
        Go source file to analyze
  -struct string
        Name of struct to analyze (required with -file)
  -help
        Show help message
```

### Examples

#### Example 1: Built-in Examples

View a basic example with padding issues:

```bash
go-mem-layout -example basic
```

Output:
```
Struct: BasicExample
Total Size: 24 bytes
Alignment: 8 bytes

Memory Layout:
┌────────┬─────────────────────────────────────────────────────────────┐
│ Offset │ Field                                                       │
├────────┼─────────────────────────────────────────────────────────────┤
│      0 │ A                    int8                 (1 bytes)         │
│      1 │ [padding]                                                   │
│      2 │ [padding]                                                   │
│      3 │ [padding]                                                   │
│      4 │ [padding]                                                   │
│      5 │ [padding]                                                   │
│      6 │ [padding]                                                   │
│      7 │ [padding]                                                   │
│      8 │ B                    int64                (8 bytes)         │
...
└────────┴─────────────────────────────────────────────────────────────┘

Padding:
  7 bytes at offset 1 (after A)
  6 bytes at offset 18 (after C)
  Total padding: 13 bytes (54.2% waste)
```

#### Example 2: Compact Format

```bash
go-mem-layout -example optimized -format compact
```

Output:
```
struct OptimizedExample {  // 32 bytes, align 8
  Int64 int64  // offset 0, size 8, align 8
  Pointer *int  // offset 8, size 8, align 8
  Float64 float64  // offset 16, size 8, align 8
  Int32 int32  // offset 24, size 4, align 4
  Int16 int16  // offset 28, size 2, align 2
  Int8 int8  // offset 30, size 1, align 1
  Bool bool  // offset 31, size 1, align 1
}
```

#### Example 3: JSON Output

```bash
go-mem-layout -example basic -format json
```

Output:
```json
{
  "name": "BasicExample",
  "total_size": 24,
  "alignment": 8,
  "fields": [
    {
      "name": "A",
      "type": "int8",
      "offset": 0,
      "size": 1,
      "alignment": 1
    },
    ...
  ],
  "padding": [
    {
      "offset": 1,
      "size": 7,
      "after": "A"
    }
  ]
}
```

#### Example 4: Analyze Source File

Create a Go file with your struct:

```go
// mystruct.go
package main

type Person struct {
    Name    string
    Age     int32
    Active  bool
    Balance float64
}
```

Analyze it:

```bash
go-mem-layout -file mystruct.go -struct Person -format compact
```

Output:
```
struct Person {  // 32 bytes, align 8
  Name string  // offset 0, size 16, align 8
  Age int32  // offset 16, size 4, align 4
  Active bool  // offset 20, size 1, align 1
  // [3 bytes padding]
  Balance float64  // offset 24, size 8, align 8
}
```

## Understanding the Output

### ASCII Format

The ASCII format provides a detailed byte-by-byte visualization:
- Shows exact offset of each byte
- Clearly marks padding bytes
- Includes field details with type information
- Calculates total padding waste percentage

### Compact Format

The compact format shows struct definition with comments:
- Looks like actual Go code
- Inline comments show offset, size, and alignment
- Padding is shown as comments between fields
- Total size and alignment in struct header

### JSON Format

The JSON format provides structured data for programmatic use:
- All numeric values for offsets, sizes, alignments
- Array of fields with complete metadata
- Array of padding information
- Easy to parse for automation

## Performance Tips

Based on the analysis, you can optimize struct layouts:

1. **Order by Size**: Place larger fields first (descending order)
2. **Group by Alignment**: Group fields with similar alignment requirements
3. **Minimize Padding**: Reorder fields to reduce padding bytes

Example optimization:

**Before (54% waste):**
```go
type Bad struct {
    A int8   // 1 byte
    B int64  // 8 bytes
    C int16  // 2 bytes
}
// Total: 24 bytes (13 bytes padding)
```

**After (0% waste):**
```go
type Good struct {
    B int64  // 8 bytes
    C int16  // 2 bytes
    A int8   // 1 byte
}
// Total: 16 bytes (no padding with proper optimization)
```

## How It Works

1. **Reflection Analysis**: For built-in examples, uses Go's `reflect` package to get actual memory layout
2. **Type Analysis**: For source files, parses Go code and uses `go/types` to calculate layout
3. **Platform-Aware**: Respects platform-specific sizes (int, uintptr, pointers)
4. **Accurate Alignment**: Follows Go's alignment rules for accurate padding calculation

## Use Cases

- **Performance Optimization**: Identify and eliminate memory waste in hot-path structs
- **Systems Programming**: Understand memory layout for FFI and unsafe operations
- **Education**: Learn about memory alignment and struct padding
- **Code Review**: Catch inefficient struct layouts during review
- **Research**: Study memory layout patterns in Go programs

## Testing

Run the test suite:

```bash
go test -v
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Max Base

## See Also

- [Go Memory Model](https://go.dev/ref/mem)
- [Go Spec - Size and alignment guarantees](https://go.dev/ref/spec#Size_and_alignment_guarantees)
- [Effective Go - Data alignment](https://go.dev/doc/effective_go)

# Examples

This directory contains example Go files that demonstrate struct memory layout analysis.

## Running Examples

You can analyze any of these files using:

```bash
go run . -file examples/example_name.go -struct StructName
```

Or with different formats:

```bash
# Compact format
go run . -file examples/example_name.go -struct StructName -format compact

# JSON format
go run . -file examples/example_name.go -struct StructName -format json
```

## Example Files

### basic.go
Demonstrates basic struct padding issues and how field ordering affects memory layout.

### optimization.go
Shows before/after examples of struct optimization to reduce memory waste.

### real_world.go
Real-world examples from common use cases like HTTP handlers, database models, etc.

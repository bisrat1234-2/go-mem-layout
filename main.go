package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

func main() {
	// Command line flags
	var (
		formatFlag   = flag.String("format", "ascii", "Output format: ascii, compact, or json")
		exampleFlag  = flag.String("example", "", "Run with example struct: basic, mixed, optimized, or all")
		sourceFile   = flag.String("file", "", "Go source file to analyze")
		structName   = flag.String("struct", "", "Name of struct to analyze (required with -file)")
		showHelp     = flag.Bool("help", false, "Show help message")
	)
	
	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	// Determine output format
	var format OutputFormat
	switch *formatFlag {
	case "ascii":
		format = FormatASCII
	case "compact":
		format = FormatCompact
	case "json":
		format = FormatJSON
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown format '%s'. Use ascii, compact, or json.\n", *formatFlag)
		os.Exit(1)
	}

	// Handle examples
	if *exampleFlag != "" {
		runExamples(*exampleFlag, format)
		return
	}

	// Handle source file analysis
	if *sourceFile != "" {
		if *structName == "" {
			fmt.Fprintf(os.Stderr, "Error: -struct flag is required when using -file\n")
			os.Exit(1)
		}
		
		data, err := os.ReadFile(*sourceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		
		layout, err := AnalyzeStructFromSource(string(data), *structName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing struct: %v\n", err)
			os.Exit(1)
		}
		
		output, err := FormatOutput(layout, format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println(output)
		return
	}

	// If no arguments, show help
	printHelp()
}

func printHelp() {
	fmt.Println("Go Memory Layout Analyzer")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Println("Analyzes Go structs and visualizes memory alignment, padding, and offsets.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go-mem-layout [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -format string")
	fmt.Println("        Output format: ascii, compact, or json (default: ascii)")
	fmt.Println("  -example string")
	fmt.Println("        Run with example struct: basic, mixed, optimized, or all")
	fmt.Println("  -file string")
	fmt.Println("        Go source file to analyze")
	fmt.Println("  -struct string")
	fmt.Println("        Name of struct to analyze (required with -file)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Show basic example in ASCII format")
	fmt.Println("  go-mem-layout -example basic")
	fmt.Println()
	fmt.Println("  # Show all examples in JSON format")
	fmt.Println("  go-mem-layout -example all -format json")
	fmt.Println()
	fmt.Println("  # Analyze a struct from source file")
	fmt.Println("  go-mem-layout -file mystruct.go -struct MyStruct")
	fmt.Println()
	fmt.Println("  # Analyze with compact output")
	fmt.Println("  go-mem-layout -file mystruct.go -struct MyStruct -format compact")
}

func runExamples(exampleName string, format OutputFormat) {
	examples := make(map[string]interface{})
	
	// Define example structs
	type BasicExample struct {
		A int8
		B int64
		C int16
	}
	
	type MixedExample struct {
		Bool    bool
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Float32 float32
		Float64 float64
		Pointer *int
		String  string
	}
	
	type OptimizedExample struct {
		Int64   int64
		Pointer *int
		Float64 float64
		Int32   int32
		Int16   int16
		Int8    int8
		Bool    bool
	}

	examples["basic"] = BasicExample{}
	examples["mixed"] = MixedExample{}
	examples["optimized"] = OptimizedExample{}

	if exampleName == "all" {
		for name, example := range examples {
			analyzeAndPrint(name, example, format)
			fmt.Println()
		}
	} else if example, ok := examples[exampleName]; ok {
		analyzeAndPrint(exampleName, example, format)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Unknown example '%s'. Available: basic, mixed, optimized, all\n", exampleName)
		os.Exit(1)
	}
}

func analyzeAndPrint(name string, example interface{}, format OutputFormat) {
	t := reflect.TypeOf(example)
	layout, err := AnalyzeStruct(t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing %s: %v\n", name, err)
		return
	}
	
	output, err := FormatOutput(layout, format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting %s: %v\n", name, err)
		return
	}
	
	fmt.Println(output)
}

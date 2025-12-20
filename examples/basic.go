package examples

// BadLayout demonstrates a poorly organized struct with significant padding
// Total: 24 bytes, Padding: 13 bytes (54.2% waste)
type BadLayout struct {
	A int8   // 1 byte
	B int64  // 8 bytes - requires 8-byte alignment, so 7 bytes padding after A
	C int16  // 2 bytes - 6 bytes trailing padding to align struct to 8 bytes
}

// GoodLayout demonstrates the same fields optimally ordered
// Total: 16 bytes, Padding: 5 bytes (31.3% waste)
type GoodLayout struct {
	B int64  // 8 bytes
	C int16  // 2 bytes
	A int8   // 1 byte
	// 5 bytes trailing padding to align to 8 bytes
}

// BestLayout demonstrates optimal ordering by size (descending)
type BestLayout struct {
	B int64  // 8 bytes
	C int16  // 2 bytes
	A int8   // 1 byte
}

// WorstCase shows the worst possible ordering
// Total: 40 bytes with maximum padding
type WorstCase struct {
	A bool    // 1 byte
	B int64   // 8 bytes (7 bytes padding before)
	C bool    // 1 byte (7 bytes padding before)
	D int64   // 8 bytes (7 bytes padding before)
	E bool    // 1 byte (7 bytes padding after)
}

// OptimizedCase shows the same fields optimally ordered
// Total: 24 bytes with minimal padding
type OptimizedCase struct {
	B int64  // 8 bytes
	D int64  // 8 bytes
	A bool   // 1 byte
	C bool   // 1 byte
	E bool   // 1 byte
	// 5 bytes trailing padding
}

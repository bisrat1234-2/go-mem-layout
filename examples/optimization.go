package examples

// UnoptimizedUser represents a typical user model before optimization
// Demonstrates common mistakes in field ordering
type UnoptimizedUser struct {
	Active    bool    // 1 byte
	ID        int64   // 8 bytes
	Role      byte    // 1 byte
	CreatedAt int64   // 8 bytes
	IsAdmin   bool    // 1 byte
	Balance   float64 // 8 bytes
	Age       int32   // 4 bytes
}

// OptimizedUser represents the same user model with fields reordered
// Fields are ordered by size (descending) to minimize padding
type OptimizedUser struct {
	ID        int64   // 8 bytes
	CreatedAt int64   // 8 bytes
	Balance   float64 // 8 bytes
	Age       int32   // 4 bytes
	Role      byte    // 1 byte
	Active    bool    // 1 byte
	IsAdmin   bool    // 1 byte
}

// UnoptimizedConfig shows a configuration struct with padding issues
type UnoptimizedConfig struct {
	Enabled   bool    // 1 byte
	Timeout   int64   // 8 bytes
	Debug     bool    // 1 byte
	MaxRetries int32  // 4 bytes
	Verbose   bool    // 1 byte
	Port      int16   // 2 bytes
}

// OptimizedConfig shows the same configuration optimized
type OptimizedConfig struct {
	Timeout    int64  // 8 bytes
	MaxRetries int32  // 4 bytes
	Port       int16  // 2 bytes
	Enabled    bool   // 1 byte
	Debug      bool   // 1 byte
	Verbose    bool   // 1 byte
}

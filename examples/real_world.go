package examples

import "time"

// HTTPRequest represents a typical HTTP request structure
// This is a common pattern in web applications
type HTTPRequest struct {
	Method      string            // 16 bytes (string header)
	URL         string            // 16 bytes
	Headers     map[string]string // 8 bytes (pointer to map)
	Body        []byte            // 24 bytes (slice header)
	Timestamp   time.Time         // 24 bytes (3 int64 fields)
	ContentType string            // 16 bytes
	StatusCode  int               // 8 bytes
	Timeout     int64             // 8 bytes
	KeepAlive   bool              // 1 byte
}

// DatabaseRecord represents a typical database model
// Shows realistic field types and ordering
type DatabaseRecord struct {
	ID        int64     // 8 bytes
	CreatedAt time.Time // 24 bytes
	UpdatedAt time.Time // 24 bytes
	UserID    int64     // 8 bytes
	Name      string    // 16 bytes
	Email     string    // 16 bytes
	Age       int32     // 4 bytes
	Active    bool      // 1 byte
	Verified  bool      // 1 byte
}

// CacheEntry represents a cache entry with metadata
type CacheEntry struct {
	Key        string    // 16 bytes
	Value      []byte    // 24 bytes
	Expiry     time.Time // 24 bytes
	Hits       int64     // 8 bytes
	Size       int32     // 4 bytes
	Compressed bool      // 1 byte
}

// MessageQueue represents a message queue item
type MessageQueue struct {
	ID        string    // 16 bytes
	Payload   []byte    // 24 bytes
	Timestamp time.Time // 24 bytes
	Priority  int32     // 4 bytes
	Retries   int32     // 4 bytes
	Processed bool      // 1 byte
	Failed    bool      // 1 byte
}

// MetricsData represents application metrics
type MetricsData struct {
	Name      string  // 16 bytes
	Value     float64 // 8 bytes
	Timestamp int64   // 8 bytes
	Count     int64   // 8 bytes
	Min       float64 // 8 bytes
	Max       float64 // 8 bytes
	Avg       float64 // 8 bytes
	Tags      string  // 16 bytes
	Source    string  // 16 bytes
}

// SessionData represents user session information
type SessionData struct {
	SessionID string            // 16 bytes
	UserID    int64             // 8 bytes
	Data      map[string]string // 8 bytes (pointer)
	CreatedAt time.Time         // 24 bytes
	ExpiresAt time.Time         // 24 bytes
	IPAddress string            // 16 bytes
	UserAgent string            // 16 bytes
	Active    bool              // 1 byte
	Verified  bool              // 1 byte
}

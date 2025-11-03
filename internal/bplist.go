package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// BinaryPlist represents a binary property list value.
// It stores both the raw binary data and the deserialized value,
// allowing round-trip conversion without data loss.
type BinaryPlist struct {
	// Raw binary plist data
	Raw []byte

	// Deserialized value (can be any valid plist type)
	Value any
}

// NewBinaryPlist creates a new BinaryPlist with the given raw data and deserialized value.
func NewBinaryPlist(raw []byte, value any) *BinaryPlist {
	return &BinaryPlist{
		Raw:   raw,
		Value: value,
	}
}

// String returns a string representation of the binary plist.
func (b *BinaryPlist) String() string {
	return fmt.Sprintf("BinaryPlist{size: %d bytes}", len(b.Raw))
}

// Base64 returns the base64-encoded representation of the raw binary data.
func (b *BinaryPlist) Base64() string {
	return base64.StdEncoding.EncodeToString(b.Raw)
}

// MarshalJSON implements json.Marshaler to allow JSON serialization.
func (b *BinaryPlist) MarshalJSON() ([]byte, error) {
	// When serializing to JSON, we'll include both representations
	data := map[string]any{
		"type":   "binary_plist",
		"base64": b.Base64(),
		"value":  b.Value,
	}
	return json.Marshal(data)
}

package hash

import (
	"crypto/sha256"
)

// Sha256 computes the SHA-256 hash of the input.
// Accepts either a hex string (with 0x prefix) or raw bytes.
// Returns the hash as a hex string by default.
//
// Example:
//
//	hash := Sha256("0x68656c6c6f20776f726c64")
//	// "0xb94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
//
//	hash := Sha256Bytes([]byte("hello world"))
//	// returns raw bytes
func Sha256(value any) string {
	return bytesToHex(Sha256Bytes(value))
}

// Sha256Bytes computes the SHA-256 hash and returns raw bytes.
func Sha256Bytes(value any) []byte {
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		if isHex(v) {
			data = hexToBytes(v)
		} else {
			data = []byte(v)
		}
	default:
		return nil
	}

	h := sha256.Sum256(data)
	return h[:]
}

package hash

import (
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Keccak256 computes the Keccak-256 hash of the input.
// Accepts either a hex string (with 0x prefix) or raw bytes.
// Returns the hash as a hex string by default.
//
// Example:
//
//	hash := Keccak256("0x68656c6c6f20776f726c64")
//	// "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
//
//	hash := Keccak256Bytes([]byte("hello world"))
//	// "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
func Keccak256(value any) string {
	return bytesToHex(Keccak256Bytes(value))
}

// Keccak256Bytes computes the Keccak-256 hash and returns raw bytes.
func Keccak256Bytes(value any) []byte {
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

	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

// Helper functions

func isHex(s string) bool {
	if len(s) < 2 {
		return false
	}
	return s[0:2] == "0x" || s[0:2] == "0X"
}

func hexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b, _ := hex.DecodeString(s)
	return b
}

func bytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

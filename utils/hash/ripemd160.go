package hash

import (
	"golang.org/x/crypto/ripemd160"
)

// Ripemd160 computes the RIPEMD-160 hash of the input.
// Accepts either a hex string (with 0x prefix) or raw bytes.
// Returns the hash as a hex string by default.
//
// Example:
//
//	hash := Ripemd160("0x68656c6c6f20776f726c64")
//	// "0x98c615784ccb5fe5936fbc0cbe9dfdb408d92f0f"
//
//	hash := Ripemd160Bytes([]byte("hello world"))
//	// returns raw bytes
func Ripemd160(value any) string {
	return bytesToHex(Ripemd160Bytes(value))
}

// Ripemd160Bytes computes the RIPEMD-160 hash and returns raw bytes.
func Ripemd160Bytes(value any) []byte {
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

	h := ripemd160.New()
	h.Write(data)
	return h.Sum(nil)
}

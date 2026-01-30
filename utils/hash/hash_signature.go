package hash

// HashSignature computes the Keccak-256 hash of a signature string.
// This is useful for computing function selectors and event topics.
//
// Example:
//
//	hash := HashSignature("transfer(address,uint256)")
//	// "0xa9059cbb2ab09eb219583f4a59a5d0623ade346d962bcd4e46b11da047c9049b"
func HashSignature(sig string) string {
	return Keccak256([]byte(sig))
}

// HashSignatureBytes computes the Keccak-256 hash of a signature and returns raw bytes.
func HashSignatureBytes(sig string) []byte {
	return Keccak256Bytes([]byte(sig))
}

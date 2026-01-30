package hash

// ToEventHash is an alias for ToSignatureHash, specifically for event definitions.
// Returns the full keccak256 hash of the event signature.
//
// Example:
//
//	hash, _ := ToEventHash("event Transfer(address indexed from, address indexed to, uint256 amount)")
//	// "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
func ToEventHash(event string) (string, error) {
	return ToSignatureHash(event)
}

// ToEventHashBytes returns the hash as raw bytes.
func ToEventHashBytes(event string) ([]byte, error) {
	return ToSignatureHashBytes(event)
}

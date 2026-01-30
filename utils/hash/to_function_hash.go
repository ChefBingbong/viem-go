package hash

// ToFunctionHash is an alias for ToSignatureHash, specifically for function definitions.
// Returns the full keccak256 hash of the function signature.
//
// Example:
//
//	hash, _ := ToFunctionHash("function transfer(address to, uint256 amount)")
//	// "0xa9059cbb2ab09eb219583f4a59a5d0623ade346d962bcd4e46b11da047c9049b"
func ToFunctionHash(fn string) (string, error) {
	return ToSignatureHash(fn)
}

// ToFunctionHashBytes returns the hash as raw bytes.
func ToFunctionHashBytes(fn string) ([]byte, error) {
	return ToSignatureHashBytes(fn)
}

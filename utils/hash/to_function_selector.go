package hash

// ToFunctionSelector returns the 4-byte function selector for a given function definition.
// The selector is the first 4 bytes of the keccak256 hash of the function signature.
//
// Example:
//
//	selector, _ := ToFunctionSelector("function ownerOf(uint256 tokenId)")
//	// "0x6352211e"
//
//	selector, _ := ToFunctionSelector("function transfer(address to, uint256 amount)")
//	// "0xa9059cbb"
func ToFunctionSelector(fn string) (string, error) {
	hash, err := ToSignatureHash(fn)
	if err != nil {
		return "", err
	}
	// Take first 4 bytes (8 hex chars + 0x prefix = 10 chars)
	return hash[:10], nil
}

// ToFunctionSelectorBytes returns the selector as raw 4 bytes.
func ToFunctionSelectorBytes(fn string) ([]byte, error) {
	hash, err := ToSignatureHashBytes(fn)
	if err != nil {
		return nil, err
	}
	return hash[:4], nil
}

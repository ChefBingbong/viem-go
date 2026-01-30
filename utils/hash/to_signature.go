package hash

// ToSignature returns the normalized signature for a given function or event definition.
// It accepts either a human-readable signature string or processes it directly.
//
// Example:
//
//	sig, _ := ToSignature("function ownerOf(uint256 tokenId)")
//	// "ownerOf(uint256)"
//
//	sig, _ := ToSignature("event Transfer(address indexed from, address indexed to, uint256 amount)")
//	// "Transfer(address,address,uint256)"
func ToSignature(def string) (string, error) {
	return NormalizeSignature(def)
}

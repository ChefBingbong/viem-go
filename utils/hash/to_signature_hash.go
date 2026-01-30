package hash

// ToSignatureHash returns the hash (keccak256) of the function/event signature.
//
// Example:
//
//	hash, _ := ToSignatureHash("function ownerOf(uint256 tokenId)")
//	// "0x6352211e6566aa027e75ac9dbf2423197fbd9b82b9d981a3ab367d355866aa1c"
//
//	hash, _ := ToSignatureHash("event Transfer(address indexed from, address indexed to, uint256 amount)")
//	// "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
func ToSignatureHash(fn string) (string, error) {
	sig, err := ToSignature(fn)
	if err != nil {
		return "", err
	}
	return HashSignature(sig), nil
}

// ToSignatureHashBytes returns the hash as raw bytes.
func ToSignatureHashBytes(fn string) ([]byte, error) {
	sig, err := ToSignature(fn)
	if err != nil {
		return nil, err
	}
	return HashSignatureBytes(sig), nil
}

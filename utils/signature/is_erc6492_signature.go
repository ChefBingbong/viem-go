package signature

import (
	"strings"
)

// IsErc6492Signature checks whether the signature is an ERC-6492 formatted signature.
// https://eips.ethereum.org/EIPS/eip-6492
//
// Example:
//
//	isErc6492 := IsErc6492Signature("0x...6492649264926492649264926492649264926492649264926492649264926492")
//	// true
func IsErc6492Signature(signature string) bool {
	// Remove 0x prefix and get the last 32 bytes (64 hex chars)
	sig := strings.TrimPrefix(signature, "0x")
	sig = strings.TrimPrefix(sig, "0X")

	// Signature must be at least 32 bytes
	if len(sig) < 64 {
		return false
	}

	// Get the last 32 bytes
	last32Bytes := "0x" + sig[len(sig)-64:]

	// Compare with magic bytes (case insensitive)
	return strings.EqualFold(last32Bytes, Erc6492MagicBytes)
}

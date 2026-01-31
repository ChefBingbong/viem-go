package accounts

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

// GeneratePrivateKey generates a random private key.
//
// Example:
//
//	privateKey := GeneratePrivateKey()
//	// "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
func GeneratePrivateKey() string {
	// Generate 32 random bytes
	privateKeyBytes := make([]byte, 32)
	_, err := rand.Read(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	// Validate that it's a valid secp256k1 private key
	_, err = crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		// Extremely rare, but if the random bytes aren't valid, try again
		return GeneratePrivateKey()
	}

	return "0x" + hex.EncodeToString(privateKeyBytes)
}

// GeneratePrivateKeyWithEntropy generates a private key using the provided entropy source.
// This is useful for deterministic testing.
func GeneratePrivateKeyWithEntropy(entropy []byte) (string, error) {
	if len(entropy) != 32 {
		return "", ErrInvalidPrivateKey
	}

	// Validate that it's a valid secp256k1 private key
	_, err := crypto.ToECDSA(entropy)
	if err != nil {
		return "", ErrInvalidPrivateKey
	}

	return "0x" + hex.EncodeToString(entropy), nil
}

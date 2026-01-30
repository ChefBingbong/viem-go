package signature

import (
	"github.com/ethereum/go-ethereum/crypto"
)

// RecoverAddress recovers the Ethereum address from a hash and signature.
//
// Example:
//
//	address, err := RecoverAddress(
//		"0xd9eba16ed0ecae432b71fe008c98cc872bb4cc214d3220a36f365326cf807d68",
//		"0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c",
//	)
//	// "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
func RecoverAddress(hash string, signature any) (string, error) {
	pubKeyHex, err := RecoverPublicKey(hash, signature)
	if err != nil {
		return "", err
	}

	return publicKeyToAddress(pubKeyHex)
}

// publicKeyToAddress converts a public key to an Ethereum address.
func publicKeyToAddress(pubKeyHex string) (string, error) {
	pubKeyBytes := hexToBytes(pubKeyHex)

	// Parse the public key
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	// Get the address
	address := crypto.PubkeyToAddress(*pubKey)

	return address.Hex(), nil
}

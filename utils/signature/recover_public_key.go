package signature

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ErrInvalidHash is returned when the hash is invalid.
	ErrInvalidHash = errors.New("invalid hash")
	// ErrRecoveryFailed is returned when public key recovery fails.
	ErrRecoveryFailed = errors.New("public key recovery failed")
)

// RecoverPublicKey recovers the public key from a hash and signature.
// Returns the uncompressed public key as a hex string (65 bytes with 04 prefix).
//
// Example:
//
//	pubKey, err := RecoverPublicKey(
//		"0xd9eba16ed0ecae432b71fe008c98cc872bb4cc214d3220a36f365326cf807d68",
//		"0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c",
//	)
func RecoverPublicKey(hash string, signature any) (string, error) {
	// Convert hash to bytes
	hashBytes := hexToBytes(hash)
	if len(hashBytes) != 32 {
		return "", fmt.Errorf("%w: expected 32 bytes, got %d", ErrInvalidHash, len(hashBytes))
	}

	// Convert signature to bytes
	sigBytes, err := signatureToRecoveryBytes(signature)
	if err != nil {
		return "", err
	}

	// Recover public key using go-ethereum's crypto package
	pubKey, err := crypto.SigToPub(hashBytes, sigBytes)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRecoveryFailed, err)
	}

	// Convert to uncompressed public key bytes (65 bytes with 04 prefix)
	pubKeyBytes := crypto.FromECDSAPub(pubKey)

	return bytesToHex(pubKeyBytes), nil
}

// RecoverPublicKeyBytes recovers the public key and returns it as bytes.
func RecoverPublicKeyBytes(hash string, signature any) ([]byte, error) {
	pubKeyHex, err := RecoverPublicKey(hash, signature)
	if err != nil {
		return nil, err
	}
	return hexToBytes(pubKeyHex), nil
}

// signatureToRecoveryBytes converts various signature formats to the 65-byte format
// expected by go-ethereum's SigToPub function.
func signatureToRecoveryBytes(signature any) ([]byte, error) {
	switch sig := signature.(type) {
	case string:
		// Hex signature string
		return parseSignatureToRecoveryBytes(sig)

	case []byte:
		// Raw bytes
		if len(sig) != 65 {
			return nil, fmt.Errorf("%w: expected 65 bytes, got %d", ErrInvalidSignatureLength, len(sig))
		}
		return convertToRecoveryFormat(sig), nil

	case *Signature:
		// Structured signature
		return structuredSignatureToBytes(sig)

	default:
		return nil, errors.New("unsupported signature type")
	}
}

// parseSignatureToRecoveryBytes parses a hex signature string to recovery bytes.
func parseSignatureToRecoveryBytes(signatureHex string) ([]byte, error) {
	sig, err := ParseSignature(signatureHex)
	if err != nil {
		return nil, err
	}
	return structuredSignatureToBytes(sig)
}

// structuredSignatureToBytes converts a Signature struct to 65-byte format.
func structuredSignatureToBytes(sig *Signature) ([]byte, error) {
	r := hexToBytes(sig.R)
	s := hexToBytes(sig.S)

	// Pad r and s to 32 bytes
	if len(r) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(r):], r)
		r = padded
	}
	if len(s) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(s):], s)
		s = padded
	}

	// Get recovery bit
	recoveryBit := toRecoveryBit(sig)

	// Construct 65-byte signature: r (32) + s (32) + recoveryBit (1)
	result := make([]byte, 65)
	copy(result[0:32], r)
	copy(result[32:64], s)
	result[64] = byte(recoveryBit)

	return result, nil
}

// toRecoveryBit extracts the recovery bit (0 or 1) from a signature.
func toRecoveryBit(sig *Signature) int {
	// Use yParity if set
	if sig.YParity == 0 || sig.YParity == 1 {
		return sig.YParity
	}

	// Otherwise derive from v
	if sig.V != nil {
		v := sig.V.Int64()
		if v == 0 || v == 1 {
			return int(v)
		}
		if v == 27 {
			return 0
		}
		if v == 28 {
			return 1
		}
	}

	return 0
}

// convertToRecoveryFormat ensures the signature is in the correct format for recovery.
// go-ethereum expects v to be 0 or 1, not 27 or 28.
func convertToRecoveryFormat(sig []byte) []byte {
	result := make([]byte, 65)
	copy(result, sig)

	// Convert v from 27/28 to 0/1 if needed
	if result[64] >= 27 {
		result[64] -= 27
	}

	return result
}

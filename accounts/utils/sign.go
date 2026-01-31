package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ErrInvalidPrivateKey is returned when the private key is invalid.
	ErrInvalidPrivateKey = errors.New("invalid private key")
	// ErrInvalidHash is returned when the hash is invalid.
	ErrInvalidHash = errors.New("invalid hash")
	// ErrSigningFailed is returned when signing fails.
	ErrSigningFailed = errors.New("signing failed")
)

// SignParameters contains parameters for signing a hash.
type SignParameters struct {
	// Hash is the 32-byte hash to sign (hex string with 0x prefix).
	Hash string
	// PrivateKey is the private key to sign with (hex string with 0x prefix).
	PrivateKey string
	// To specifies the output format (object, hex, or bytes). Default is "object".
	To SignReturnFormat
}

// SignResult can be a Signature, hex string, or byte slice depending on the To parameter.
type SignResult struct {
	Signature *Signature
	Hex       string
	Bytes     []byte
}

// Sign signs a hash with a given private key.
//
// Example:
//
//	result, err := Sign(SignParameters{
//		Hash:       "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
//		PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
//		To:         SignReturnFormatObject,
//	})
func Sign(params SignParameters) (*SignResult, error) {
	// Parse the private key
	privateKey, err := parsePrivateKey(params.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrivateKey, err)
	}

	// Parse the hash
	hashBytes, err := parseHash(params.Hash)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidHash, err)
	}

	// Sign the hash
	sigBytes, err := crypto.Sign(hashBytes, privateKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSigningFailed, err)
	}

	// Extract r, s, v
	r := sigBytes[:32]
	s := sigBytes[32:64]
	recovery := sigBytes[64]

	// Create signature object
	sig := &Signature{
		R:       "0x" + hex.EncodeToString(r),
		S:       "0x" + hex.EncodeToString(s),
		V:       big.NewInt(int64(recovery) + 27),
		YParity: int(recovery),
	}

	// Return based on format
	to := params.To
	if to == "" {
		to = SignReturnFormatObject
	}

	result := &SignResult{}

	switch to {
	case SignReturnFormatObject:
		result.Signature = sig
	case SignReturnFormatHex:
		result.Hex = serializeSignatureToHex(sig)
	case SignReturnFormatBytes:
		result.Bytes = serializeSignatureToBytes(sig)
	default:
		result.Signature = sig
	}

	return result, nil
}

// SignToSignature is a convenience function that signs and returns a Signature.
func SignToSignature(hash, privateKey string) (*Signature, error) {
	result, err := Sign(SignParameters{
		Hash:       hash,
		PrivateKey: privateKey,
		To:         SignReturnFormatObject,
	})
	if err != nil {
		return nil, err
	}
	return result.Signature, nil
}

// SignToHex is a convenience function that signs and returns a hex string.
func SignToHex(hash, privateKey string) (string, error) {
	result, err := Sign(SignParameters{
		Hash:       hash,
		PrivateKey: privateKey,
		To:         SignReturnFormatHex,
	})
	if err != nil {
		return "", err
	}
	return result.Hex, nil
}

// SignToBytes is a convenience function that signs and returns bytes.
func SignToBytes(hash, privateKey string) ([]byte, error) {
	result, err := Sign(SignParameters{
		Hash:       hash,
		PrivateKey: privateKey,
		To:         SignReturnFormatBytes,
	})
	if err != nil {
		return nil, err
	}
	return result.Bytes, nil
}

// parsePrivateKey parses a hex-encoded private key.
func parsePrivateKey(key string) (*ecdsa.PrivateKey, error) {
	key = strings.TrimPrefix(key, "0x")
	key = strings.TrimPrefix(key, "0X")
	return crypto.HexToECDSA(key)
}

// parseHash parses a hex-encoded hash (must be 32 bytes).
func parseHash(hash string) ([]byte, error) {
	hash = strings.TrimPrefix(hash, "0x")
	hash = strings.TrimPrefix(hash, "0X")

	if len(hash) != 64 {
		return nil, fmt.Errorf("hash must be 32 bytes (64 hex chars), got %d", len(hash)/2)
	}

	return hex.DecodeString(hash)
}

// serializeSignatureToHex serializes a signature to hex format.
func serializeSignatureToHex(sig *Signature) string {
	r := strings.TrimPrefix(sig.R, "0x")
	s := strings.TrimPrefix(sig.S, "0x")

	// Pad r and s to 64 characters
	r = padLeft(r, 64)
	s = padLeft(s, 64)

	// Add recovery byte (27 = 0x1b for yParity=0, 28 = 0x1c for yParity=1)
	var vByte string
	if sig.YParity == 0 {
		vByte = "1b"
	} else {
		vByte = "1c"
	}

	return "0x" + r + s + vByte
}

// serializeSignatureToBytes serializes a signature to bytes.
func serializeSignatureToBytes(sig *Signature) []byte {
	hexStr := serializeSignatureToHex(sig)
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, _ := hex.DecodeString(hexStr)
	return b
}

// padLeft pads a string with zeros on the left.
func padLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat("0", length-len(s)) + s
}

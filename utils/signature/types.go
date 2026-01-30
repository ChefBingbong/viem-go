package signature

import "math/big"

// Signature represents a parsed ECDSA signature.
type Signature struct {
	// R component of the signature (32 bytes as hex string).
	R string `json:"r"`
	// S component of the signature (32 bytes as hex string).
	S string `json:"s"`
	// V value (27 or 28 for legacy, optional for EIP-2930/1559).
	V *big.Int `json:"v,omitempty"`
	// YParity is the parity of the y-coordinate of the curve point (0 or 1).
	YParity int `json:"yParity"`
}

// CompactSignature represents an EIP-2098 compact signature.
// https://eips.ethereum.org/EIPS/eip-2098
type CompactSignature struct {
	// R component of the signature (32 bytes as hex string).
	R string `json:"r"`
	// YParityAndS encodes the yParity in the top bit of s.
	YParityAndS string `json:"yParityAndS"`
}

// SignableMessage represents a message that can be signed.
// It can be either a plain string or a raw bytes/hex message.
type SignableMessage struct {
	// Raw is the raw bytes or hex string to sign (takes precedence over Message).
	Raw any // can be []byte or string (hex)
	// Message is a string message to sign.
	Message string
}

// NewSignableMessage creates a SignableMessage from a string.
func NewSignableMessage(message string) SignableMessage {
	return SignableMessage{Message: message}
}

// NewSignableMessageRaw creates a SignableMessage from raw bytes.
func NewSignableMessageRaw(raw []byte) SignableMessage {
	return SignableMessage{Raw: raw}
}

// NewSignableMessageRawHex creates a SignableMessage from a raw hex string.
func NewSignableMessageRawHex(raw string) SignableMessage {
	return SignableMessage{Raw: raw}
}

// Erc6492Signature represents a parsed ERC-6492 signature.
type Erc6492Signature struct {
	// Address is the ERC-4337 Account Factory or preparation address.
	// Empty if the signature is not in ERC-6492 format.
	Address string `json:"address,omitempty"`
	// Data is the calldata to pass to deploy account.
	// Empty if the signature is not in ERC-6492 format.
	Data string `json:"data,omitempty"`
	// Signature is the original signature.
	Signature string `json:"signature"`
}

// TypedDataDomain represents the EIP-712 domain.
type TypedDataDomain struct {
	Name              string   `json:"name,omitempty"`
	Version           string   `json:"version,omitempty"`
	ChainId           *big.Int `json:"chainId,omitempty"`
	VerifyingContract string   `json:"verifyingContract,omitempty"`
	Salt              string   `json:"salt,omitempty"`
}

// TypedDataField represents a single field in a typed data type.
type TypedDataField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// TypedDataDefinition represents the full EIP-712 typed data.
type TypedDataDefinition struct {
	Domain      TypedDataDomain               `json:"domain"`
	Types       map[string][]TypedDataField   `json:"types"`
	PrimaryType string                        `json:"primaryType"`
	Message     map[string]any                `json:"message"`
}

// Constants

// PresignMessagePrefix is the Ethereum Signed Message prefix.
const PresignMessagePrefix = "\x19Ethereum Signed Message:\n"

// Erc6492MagicBytes is the magic suffix for ERC-6492 signatures.
// https://eips.ethereum.org/EIPS/eip-6492
const Erc6492MagicBytes = "0x6492649264926492649264926492649264926492649264926492649264926492"

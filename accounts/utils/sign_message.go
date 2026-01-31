package utils

import (
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// SignMessageParameters contains parameters for signing a message.
type SignMessageParameters struct {
	// Message is the message to sign.
	Message signature.SignableMessage
	// PrivateKey is the private key to sign with (hex string with 0x prefix).
	PrivateKey string
}

// SignMessage calculates an Ethereum-specific signature in EIP-191 format:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Example:
//
//	sig, err := SignMessage(SignMessageParameters{
//		Message:    signature.NewSignableMessage("hello world"),
//		PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
//	})
func SignMessage(params SignMessageParameters) (string, error) {
	// Hash the message using EIP-191 format
	messageHash := signature.HashMessage(params.Message)

	// Sign the hash and return as hex
	return SignToHex(messageHash, params.PrivateKey)
}

// SignMessageToSignature signs a message and returns a Signature struct.
func SignMessageToSignature(params SignMessageParameters) (*Signature, error) {
	// Hash the message using EIP-191 format
	messageHash := signature.HashMessage(params.Message)

	// Sign the hash
	return SignToSignature(messageHash, params.PrivateKey)
}

// MustSignMessage signs a message or panics on error.
func MustSignMessage(params SignMessageParameters) string {
	result, err := SignMessage(params)
	if err != nil {
		panic(err)
	}
	return result
}

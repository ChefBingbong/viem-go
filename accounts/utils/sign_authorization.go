package utils

import (
	"github.com/ChefBingbong/viem-go/utils/authorization"
)

// SignAuthorizationParameters contains parameters for signing an EIP-7702 authorization.
type SignAuthorizationParameters struct {
	// Address is the contract address being authorized.
	Address string
	// ContractAddress is an alias for Address (for compatibility).
	ContractAddress string
	// ChainId is the chain ID for this authorization.
	ChainId int
	// Nonce is the account nonce for this authorization.
	Nonce int
	// PrivateKey is the private key to sign with (hex string with 0x prefix).
	PrivateKey string
	// To specifies the output format (object, hex, or bytes). Default is "object".
	To SignReturnFormat
}

// GetAddress returns the address, preferring ContractAddress if set.
func (p *SignAuthorizationParameters) GetAddress() string {
	if p.ContractAddress != "" {
		return p.ContractAddress
	}
	return p.Address
}

// SignAuthorizationResult can be a SignedAuthorization, hex string, or byte slice.
type SignAuthorizationResult struct {
	SignedAuthorization *SignedAuthorization
	Hex                 string
	Bytes               []byte
}

// SignAuthorization signs an Authorization hash in EIP-7702 format:
// keccak256('0x05' || rlp([chain_id, address, nonce]))
//
// Example:
//
//	result, err := SignAuthorization(SignAuthorizationParameters{
//		Address:    "0x1234567890123456789012345678901234567890",
//		ChainId:    1,
//		Nonce:      0,
//		PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
//	})
func SignAuthorization(params SignAuthorizationParameters) (*SignAuthorizationResult, error) {
	address := params.GetAddress()

	// Create authorization request
	authRequest := authorization.AuthorizationRequest{
		Address: address,
		ChainId: params.ChainId,
		Nonce:   params.Nonce,
	}

	// Hash the authorization
	authHash, err := authorization.HashAuthorizationHex(authRequest)
	if err != nil {
		return nil, err
	}

	// Sign the hash
	sig, err := SignToSignature(authHash, params.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Determine output format
	to := params.To
	if to == "" {
		to = SignReturnFormatObject
	}

	result := &SignAuthorizationResult{}

	switch to {
	case SignReturnFormatObject:
		result.SignedAuthorization = &SignedAuthorization{
			Address: address,
			ChainId: params.ChainId,
			Nonce:   params.Nonce,
			R:       sig.R,
			S:       sig.S,
			V:       sig.V,
			YParity: sig.YParity,
		}
	case SignReturnFormatHex:
		result.Hex = serializeSignatureToHex(sig)
	case SignReturnFormatBytes:
		result.Bytes = serializeSignatureToBytes(sig)
	default:
		result.SignedAuthorization = &SignedAuthorization{
			Address: address,
			ChainId: params.ChainId,
			Nonce:   params.Nonce,
			R:       sig.R,
			S:       sig.S,
			V:       sig.V,
			YParity: sig.YParity,
		}
	}

	return result, nil
}

// SignAuthorizationToObject signs an authorization and returns a SignedAuthorization.
func SignAuthorizationToObject(params SignAuthorizationParameters) (*SignedAuthorization, error) {
	params.To = SignReturnFormatObject
	result, err := SignAuthorization(params)
	if err != nil {
		return nil, err
	}
	return result.SignedAuthorization, nil
}

// SignAuthorizationToHex signs an authorization and returns a hex string.
func SignAuthorizationToHex(params SignAuthorizationParameters) (string, error) {
	params.To = SignReturnFormatHex
	result, err := SignAuthorization(params)
	if err != nil {
		return "", err
	}
	return result.Hex, nil
}

// MustSignAuthorization signs an authorization or panics on error.
func MustSignAuthorization(params SignAuthorizationParameters) *SignAuthorizationResult {
	result, err := SignAuthorization(params)
	if err != nil {
		panic(err)
	}
	return result
}

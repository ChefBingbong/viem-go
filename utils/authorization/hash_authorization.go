package authorization

import (
	"fmt"

	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/hash"
)

// AuthorizationRequest represents an EIP-7702 authorization request.
type AuthorizationRequest struct {
	// Address is the contract address being authorized.
	// Can also use ContractAddress for compatibility.
	Address string `json:"address,omitempty"`
	// ContractAddress is an alias for Address.
	ContractAddress string `json:"contractAddress,omitempty"`
	// ChainId is the chain ID for this authorization.
	ChainId int `json:"chainId"`
	// Nonce is the account nonce for this authorization.
	Nonce int `json:"nonce"`
}

// GetAddress returns the address, preferring ContractAddress if set.
func (a *AuthorizationRequest) GetAddress() string {
	if a.ContractAddress != "" {
		return a.ContractAddress
	}
	return a.Address
}

// HashAuthorizationParameters contains parameters for hashing an authorization.
type HashAuthorizationParameters struct {
	AuthorizationRequest
	// To specifies the output format ("hex" or "bytes"). Default is "hex".
	To string
}

// HashAuthorization computes an Authorization hash in EIP-7702 format:
// keccak256('0x05' || rlp([chain_id, address, nonce]))
//
// Example:
//
//	hash, err := HashAuthorization(HashAuthorizationParameters{
//		AuthorizationRequest: AuthorizationRequest{
//			Address: "0x1234567890123456789012345678901234567890",
//			ChainId: 1,
//			Nonce:   0,
//		},
//	})
func HashAuthorization(params HashAuthorizationParameters) (string, error) {
	return HashAuthorizationHex(params.AuthorizationRequest)
}

// HashAuthorizationHex computes the authorization hash and returns a hex string.
func HashAuthorizationHex(auth AuthorizationRequest) (string, error) {
	encoded, err := encodeAuthorization(auth)
	if err != nil {
		return "", err
	}
	return hash.Keccak256(encoded), nil
}

// HashAuthorizationBytes computes the authorization hash and returns bytes.
func HashAuthorizationBytes(auth AuthorizationRequest) ([]byte, error) {
	encoded, err := encodeAuthorization(auth)
	if err != nil {
		return nil, err
	}
	return hash.Keccak256Bytes(encoded), nil
}

// encodeAuthorization encodes an authorization for hashing.
// Format: 0x05 || rlp([chain_id, address, nonce])
func encodeAuthorization(auth AuthorizationRequest) (string, error) {
	address := auth.GetAddress()
	if address == "" {
		return "", fmt.Errorf("address is required")
	}

	// Build the RLP list: [chainId, address, nonce]
	chainIdHex := "0x"
	if auth.ChainId > 0 {
		chainIdHex = fmt.Sprintf("0x%x", auth.ChainId)
	}

	nonceHex := "0x"
	if auth.Nonce > 0 {
		nonceHex = fmt.Sprintf("0x%x", auth.Nonce)
	}

	// RLP encode the list
	rlpEncoded, err := encoding.RlpEncodeToHex([]any{
		chainIdHex,
		address,
		nonceHex,
	})
	if err != nil {
		return "", fmt.Errorf("failed to RLP encode authorization: %w", err)
	}

	// Concatenate 0x05 with the RLP encoded data
	return concatHex("0x05", rlpEncoded), nil
}

// concatHex concatenates hex strings.
func concatHex(hexStrings ...string) string {
	result := "0x"
	for _, h := range hexStrings {
		// Remove 0x prefix if present
		if len(h) >= 2 && (h[:2] == "0x" || h[:2] == "0X") {
			h = h[2:]
		}
		result += h
	}
	return result
}

package utils

import (
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// SignTypedDataParameters contains parameters for signing EIP-712 typed data.
type SignTypedDataParameters struct {
	// Domain contains the EIP-712 domain parameters.
	Domain signature.TypedDataDomain
	// Types contains the type definitions.
	Types map[string][]signature.TypedDataField
	// PrimaryType is the primary type being signed.
	PrimaryType string
	// Message is the structured message to sign.
	Message map[string]any
	// PrivateKey is the private key to sign with (hex string with 0x prefix).
	PrivateKey string
}

// SignTypedData signs typed data and calculates an Ethereum-specific signature in EIP-712 format:
// sign(keccak256("\x19\x01" ‖ domainSeparator ‖ hashStruct(message)))
//
// Example:
//
//	sig, err := SignTypedData(SignTypedDataParameters{
//		Domain: signature.TypedDataDomain{
//			Name:              "Ether Mail",
//			Version:           "1",
//			ChainId:           big.NewInt(1),
//			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
//		},
//		Types: map[string][]signature.TypedDataField{
//			"Person": {
//				{Name: "name", Type: "string"},
//				{Name: "wallet", Type: "address"},
//			},
//			"Mail": {
//				{Name: "from", Type: "Person"},
//				{Name: "to", Type: "Person"},
//				{Name: "contents", Type: "string"},
//			},
//		},
//		PrimaryType: "Mail",
//		Message: map[string]any{
//			"from": map[string]any{
//				"name":   "Cow",
//				"wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
//			},
//			"to": map[string]any{
//				"name":   "Bob",
//				"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
//			},
//			"contents": "Hello, Bob!",
//		},
//		PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
//	})
func SignTypedData(params SignTypedDataParameters) (string, error) {
	// Create typed data definition
	typedData := signature.TypedDataDefinition{
		Domain:      params.Domain,
		Types:       params.Types,
		PrimaryType: params.PrimaryType,
		Message:     params.Message,
	}

	// Hash the typed data
	dataHash, err := signature.HashTypedData(typedData)
	if err != nil {
		return "", err
	}

	// Sign the hash and return as hex
	return SignToHex(dataHash, params.PrivateKey)
}

// SignTypedDataToSignature signs typed data and returns a Signature struct.
func SignTypedDataToSignature(params SignTypedDataParameters) (*Signature, error) {
	// Create typed data definition
	typedData := signature.TypedDataDefinition{
		Domain:      params.Domain,
		Types:       params.Types,
		PrimaryType: params.PrimaryType,
		Message:     params.Message,
	}

	// Hash the typed data
	dataHash, err := signature.HashTypedData(typedData)
	if err != nil {
		return nil, err
	}

	// Sign the hash
	return SignToSignature(dataHash, params.PrivateKey)
}

// MustSignTypedData signs typed data or panics on error.
func MustSignTypedData(params SignTypedDataParameters) string {
	result, err := SignTypedData(params)
	if err != nil {
		panic(err)
	}
	return result
}

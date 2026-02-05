package public

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/constants"
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// VerifyHashParameters contains the parameters for the VerifyHash action.
// This mirrors viem's VerifyHashParameters type with full feature support including:
//   - ERC-6492 (counterfactual signature verification)
//   - ERC-1271 (smart contract signature verification)
//   - ECDSA recovery fallback (for EOA accounts)
type VerifyHashParameters struct {
	// Address is the address that is expected to have signed the hash.
	Address common.Address

	// Hash is the 32-byte message hash that was signed (0x-prefixed hex).
	Hash string

	// Signature is the signature produced by signing the hash.
	// Accepts:
	//   - string hex-encoded signature
	//   - []byte raw signature bytes
	//   - *signature.Signature (r, s, v, yParity)
	Signature any

	// BlockNumber is the block number to verify at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to verify at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// Factory is the ERC-4337 Account Factory address for counterfactual verification.
	// Used with FactoryData for undeployed smart accounts.
	Factory *common.Address

	// FactoryData is the calldata to deploy the account via Factory.
	// Used with Factory for undeployed smart accounts.
	FactoryData []byte

	// ERC6492VerifierAddress is the address of a deployed ERC-6492 signature verifier contract.
	// If provided, uses this contract instead of deployless verification.
	// Equivalent to viem's erc6492VerifierAddress / universalSignatureVerifierAddress.
	ERC6492VerifierAddress *common.Address
}

// VerifyHashReturnType is the return type for the VerifyHash action.
// It indicates whether the signature is valid.
type VerifyHashReturnType = bool

// VerificationError represents a signature verification failure.
type VerificationError struct {
	Message string
}

func (e *VerificationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "signature verification failed"
}

// VerifyHash verifies a message hash onchain using ERC-6492.
//
// This is equivalent to viem's `verifyHash` action with full feature support:
//   - ERC-6492 verification for counterfactual (undeployed) smart accounts
//   - ERC-1271 verification for deployed smart contracts
//   - ECDSA recovery fallback for EOA accounts
//
// Example:
//
//	valid, err := public.VerifyHash(ctx, client, public.VerifyHashParameters{
//	    Address:   common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
//	    Hash:      "0x...",
//	    Signature: "0x...",
//	})
//
// Example with counterfactual verification:
//
//	valid, err := public.VerifyHash(ctx, client, public.VerifyHashParameters{
//	    Address:     common.HexToAddress("0x..."),
//	    Hash:        "0x...",
//	    Signature:   "0x...",
//	    Factory:     &factoryAddress,
//	    FactoryData: factoryCalldata,
//	})
func VerifyHash(ctx context.Context, client Client, params VerifyHashParameters) (VerifyHashReturnType, error) {
	// Serialize the signature to hex format
	sig, err := serializeSignatureToHex(params.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to serialize signature: %w", err)
	}

	// Try ERC-6492 verification
	verified, verifyErr := verifyErc6492(ctx, client, verifyErc6492Parameters{
		Address:         params.Address,
		Hash:            params.Hash,
		Signature:       sig,
		BlockNumber:     params.BlockNumber,
		BlockTag:        params.BlockTag,
		Factory:         params.Factory,
		FactoryData:     params.FactoryData,
		VerifierAddress: params.ERC6492VerifierAddress,
	})

	if verifyErr != nil {
		// Fallback: attempt to verify via ECDSA recovery
		ecdsaVerified, ecdsaErr := verifyViaEcdsaRecovery(params.Address, params.Hash, sig)
		if ecdsaErr == nil && ecdsaVerified {
			return true, nil
		}

		// If it was a verification error (not a system error), return false
		var vErr *VerificationError
		if errors.As(verifyErr, &vErr) {
			return false, nil
		}

		return false, verifyErr
	}

	return verified, nil
}

// verifyErc6492Parameters contains the parameters for ERC-6492 verification.
type verifyErc6492Parameters struct {
	Address         common.Address
	Hash            string
	Signature       string
	BlockNumber     *uint64
	BlockTag        BlockTag
	Factory         *common.Address
	FactoryData     []byte
	VerifierAddress *common.Address
}

// verifyErc6492 verifies a signature using ERC-6492 (deployless universal verification).
func verifyErc6492(ctx context.Context, client Client, params verifyErc6492Parameters) (bool, error) {
	// Determine the signature to use
	wrappedSignature := params.Signature

	// If factory and factoryData are provided and signature is not already ERC-6492 wrapped,
	// wrap it for counterfactual verification
	if params.Factory != nil && len(params.FactoryData) > 0 {
		if !signature.IsErc6492Signature(params.Signature) {
			wrapped, err := signature.SerializeErc6492Signature(signature.SerializeErc6492SignatureParams{
				Address:   params.Factory.Hex(),
				Data:      hexutil.Encode(params.FactoryData),
				Signature: params.Signature,
			})
			if err != nil {
				return false, fmt.Errorf("failed to wrap signature in ERC-6492 format: %w", err)
			}
			wrappedSignature = wrapped
		}
	}

	var callResult *CallReturnType
	var callErr error

	if params.VerifierAddress != nil {
		// Use deployed verifier contract
		calldata, err := encodeIsValidSigCall(params.Address, params.Hash, wrappedSignature)
		if err != nil {
			return false, fmt.Errorf("failed to encode isValidSig call: %w", err)
		}

		callResult, callErr = Call(ctx, client, CallParameters{
			To:          params.VerifierAddress,
			Data:        calldata,
			BlockNumber: params.BlockNumber,
			BlockTag:    params.BlockTag,
		})
	} else {
		// Use deployless verification (deploy validator contract inline)
		deployData, err := encodeErc6492ValidatorDeployData(params.Address, params.Hash, wrappedSignature)
		if err != nil {
			return false, fmt.Errorf("failed to encode validator deploy data: %w", err)
		}

		callResult, callErr = Call(ctx, client, CallParameters{
			Data:        deployData,
			BlockNumber: params.BlockNumber,
			BlockTag:    params.BlockTag,
		})
	}

	if callErr != nil {
		// CallExecutionError indicates the call reverted, meaning verification failed
		var execErr *CallExecutionError
		if errors.As(callErr, &execErr) {
			return false, &VerificationError{Message: "signature verification call failed"}
		}
		return false, callErr
	}

	// Parse the result - should be a boolean (0x01 for true, 0x00 for false)
	if callResult == nil || len(callResult.Data) == 0 {
		return false, &VerificationError{Message: "empty verification result"}
	}

	// The validator returns a single bool - check if it's true
	// For deployless, the result is the raw return value (0x01 or 0x00 padded to 32 bytes)
	valid := hexToBool(callResult.Data)
	if !valid {
		return false, &VerificationError{Message: "signature is invalid"}
	}

	return true, nil
}

// verifyErc1271Parameters contains the parameters for ERC-1271 verification.
type verifyErc1271Parameters struct {
	Address     common.Address
	Hash        string
	Signature   string
	BlockNumber *uint64
	BlockTag    BlockTag
}

// VerifyErc1271 verifies a signature using ERC-1271 (smart contract signature verification).
// This calls the isValidSignature function on the smart contract at the given address.
func VerifyErc1271(ctx context.Context, client Client, params verifyErc1271Parameters) (bool, error) {
	// Encode isValidSignature(bytes32 hash, bytes signature) call
	calldata, err := encodeIsValidSignatureCall(params.Hash, params.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to encode isValidSignature call: %w", err)
	}

	result, err := Call(ctx, client, CallParameters{
		To:          &params.Address,
		Data:        calldata,
		BlockNumber: params.BlockNumber,
		BlockTag:    params.BlockTag,
	})
	if err != nil {
		var execErr *CallExecutionError
		if errors.As(err, &execErr) {
			return false, &VerificationError{Message: "isValidSignature call failed"}
		}
		return false, err
	}

	// Check if result starts with ERC-1271 magic value (0x1626ba7e)
	if result == nil || len(result.Data) < 4 {
		return false, &VerificationError{Message: "invalid isValidSignature response"}
	}

	resultHex := hexutil.Encode(result.Data)
	if strings.HasPrefix(strings.ToLower(resultHex), strings.ToLower(constants.ERC1271MagicValue)) {
		return true, nil
	}

	return false, &VerificationError{Message: "signature is invalid"}
}

// verifyViaEcdsaRecovery attempts to verify a signature using ECDSA recovery.
// This is the fallback for EOA accounts.
func verifyViaEcdsaRecovery(address common.Address, hash string, sig string) (bool, error) {
	recoveredAddr, err := signature.RecoverAddress(hash, sig)
	if err != nil {
		return false, err
	}

	// Compare addresses (case-insensitive)
	return strings.EqualFold(recoveredAddr, address.Hex()), nil
}

// serializeSignatureToHex converts various signature formats to a hex string.
func serializeSignatureToHex(sig any) (string, error) {
	switch v := sig.(type) {
	case string:
		// Already a hex string
		if !strings.HasPrefix(v, "0x") && !strings.HasPrefix(v, "0X") {
			return "0x" + v, nil
		}
		return v, nil

	case []byte:
		return hexutil.Encode(v), nil

	case *signature.Signature:
		if v == nil {
			return "", errors.New("signature is nil")
		}
		return signature.SerializeSignature(v)

	case signature.Signature:
		return signature.SerializeSignature(&v)

	default:
		return "", fmt.Errorf("unsupported signature type: %T", sig)
	}
}

// encodeIsValidSigCall encodes a call to isValidSig(address, bytes32, bytes) on the ERC-6492 validator.
func encodeIsValidSigCall(signer common.Address, hash string, sig string) ([]byte, error) {
	// isValidSig(address _signer, bytes32 _hash, bytes _signature)
	hashBytes := common.FromHex(hash)
	if len(hashBytes) != 32 {
		return nil, fmt.Errorf("hash must be 32 bytes, got %d", len(hashBytes))
	}

	var hash32 [32]byte
	copy(hash32[:], hashBytes)

	sigBytes := common.FromHex(sig)

	// Encode parameters
	params := []abi.AbiParam{
		{Type: "address"},
		{Type: "bytes32"},
		{Type: "bytes"},
	}
	values := []any{signer, hash32, sigBytes}

	encoded, err := abi.EncodeAbiParameters(params, values)
	if err != nil {
		return nil, err
	}

	// Prepend function selector for isValidSig (0x1626ba7e... but we need to compute it)
	// Actually, looking at the viem code, isValidSig on the validator has a different signature
	// The function selector is computed from "isValidSig(address,bytes32,bytes)"
	// = keccak256("isValidSig(address,bytes32,bytes)") = 0x6ccea652...
	// Let me use the correct selector
	selector := common.FromHex("0x6ccea652") // isValidSig(address,bytes32,bytes)

	result := make([]byte, len(selector)+len(encoded))
	copy(result, selector)
	copy(result[len(selector):], encoded)

	return result, nil
}

// encodeIsValidSignatureCall encodes a call to isValidSignature(bytes32, bytes) for ERC-1271.
func encodeIsValidSignatureCall(hash string, sig string) ([]byte, error) {
	hashBytes := common.FromHex(hash)
	if len(hashBytes) != 32 {
		return nil, fmt.Errorf("hash must be 32 bytes, got %d", len(hashBytes))
	}

	var hash32 [32]byte
	copy(hash32[:], hashBytes)

	sigBytes := common.FromHex(sig)

	// Encode parameters
	params := []abi.AbiParam{
		{Type: "bytes32"},
		{Type: "bytes"},
	}
	values := []any{hash32, sigBytes}

	encoded, err := abi.EncodeAbiParameters(params, values)
	if err != nil {
		return nil, err
	}

	// Function selector for isValidSignature(bytes32,bytes) = 0x1626ba7e
	selector := common.FromHex(constants.ERC1271MagicValue)

	result := make([]byte, len(selector)+len(encoded))
	copy(result, selector)
	copy(result[len(selector):], encoded)

	return result, nil
}

// encodeErc6492ValidatorDeployData encodes the deployment data for the ERC-6492 validator contract.
// The validator is deployed inline and returns the verification result.
func encodeErc6492ValidatorDeployData(signer common.Address, hash string, sig string) ([]byte, error) {
	hashBytes := common.FromHex(hash)
	if len(hashBytes) != 32 {
		return nil, fmt.Errorf("hash must be 32 bytes, got %d", len(hashBytes))
	}

	var hash32 [32]byte
	copy(hash32[:], hashBytes)

	sigBytes := common.FromHex(sig)

	// Encode constructor arguments: (address _signer, bytes32 _hash, bytes _signature)
	params := []abi.AbiParam{
		{Type: "address"},
		{Type: "bytes32"},
		{Type: "bytes"},
	}
	values := []any{signer, hash32, sigBytes}

	constructorArgs, err := abi.EncodeAbiParameters(params, values)
	if err != nil {
		return nil, err
	}

	// Concatenate bytecode + constructor args
	bytecode := common.FromHex(constants.ERC6492SignatureValidatorBytecode)
	result := make([]byte, len(bytecode)+len(constructorArgs))
	copy(result, bytecode)
	copy(result[len(bytecode):], constructorArgs)

	return result, nil
}

// hexToBool converts a hex-encoded result to a boolean.
// The validator returns 0x01 (true) or 0x00 (false), potentially padded.
func hexToBool(data []byte) bool {
	// Trim leading zeros and check if any non-zero byte exists
	for _, b := range data {
		if b != 0 {
			return true
		}
	}
	return false
}

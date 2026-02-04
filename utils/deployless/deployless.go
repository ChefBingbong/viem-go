// Package deployless provides utilities for deployless call encoding.
// Deployless calls allow executing contract code without deploying it.
package deployless

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/constants"
)

// ToDeploylessCallViaBytecodeData creates the calldata for a deployless call
// using bytecode directly.
//
// This encodes: constructor(bytes code, bytes data)
// with the deploylessCallViaBytecodeBytecode as the deployment bytecode.
//
// Parameters:
//   - code: The contract bytecode to execute
//   - data: The calldata to send to the contract
//
// Returns the encoded deployment data that, when sent as a transaction with
// no 'to' address, will deploy a temporary contract, execute the call, and
// return the result.
func ToDeploylessCallViaBytecodeData(code, data []byte) ([]byte, error) {
	// Encode constructor arguments: (bytes, bytes)
	params := []abi.AbiParam{
		{Type: "bytes"},
		{Type: "bytes"},
	}
	values := []any{code, data}

	encodedArgs, err := abi.EncodeAbiParameters(params, values)
	if err != nil {
		return nil, err
	}

	// Decode the bytecode constant
	bytecode := common.FromHex(constants.DeploylessCallViaBytecodeBytecode)

	// Concatenate bytecode + encoded arguments
	result := make([]byte, len(bytecode)+len(encodedArgs))
	copy(result, bytecode)
	copy(result[len(bytecode):], encodedArgs)

	return result, nil
}

// ToDeploylessCallViaFactoryData creates the calldata for a deployless call
// using a factory contract (e.g., Create2 factory, Smart Account factory).
//
// This encodes: constructor(address to, bytes data, address factory, bytes factoryData)
// with the deploylessCallViaFactoryBytecode as the deployment bytecode.
//
// Parameters:
//   - to: The contract address to call (may not exist yet)
//   - data: The calldata to send to the contract
//   - factory: The factory contract address (e.g., Create2)
//   - factoryData: The calldata to send to the factory to deploy the contract
//
// Returns the encoded deployment data that, when sent as a transaction with
// no 'to' address, will:
// 1. Deploy the target contract via the factory (if not already deployed)
// 2. Execute the call to the target contract
// 3. Return the result
func ToDeploylessCallViaFactoryData(to common.Address, data []byte, factory common.Address, factoryData []byte) ([]byte, error) {
	// Encode constructor arguments: (address, bytes, address, bytes)
	params := []abi.AbiParam{
		{Type: "address"},
		{Type: "bytes"},
		{Type: "address"},
		{Type: "bytes"},
	}
	values := []any{to, data, factory, factoryData}

	encodedArgs, err := abi.EncodeAbiParameters(params, values)
	if err != nil {
		return nil, err
	}

	// Decode the bytecode constant
	bytecode := common.FromHex(constants.DeploylessCallViaFactoryBytecode)

	// Concatenate bytecode + encoded arguments
	result := make([]byte, len(bytecode)+len(encodedArgs))
	copy(result, bytecode)
	copy(result[len(bytecode):], encodedArgs)

	return result, nil
}

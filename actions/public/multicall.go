package public

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/sync/errgroup"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/constants"
	"github.com/ChefBingbong/viem-go/utils/deployless"
)

// MulticallContract defines a contract call for multicall.
// This mirrors viem's ContractFunctionParameters type.
type MulticallContract struct {
	// Address is the contract address to call.
	Address common.Address

	// ABI is the contract ABI as JSON bytes, string, or *abi.ABI.
	ABI any

	// FunctionName is the name of the function to call.
	FunctionName string

	// Args are the function arguments.
	Args []any
}

// MulticallParameters contains the parameters for the Multicall action.
// This mirrors viem's MulticallParameters type.
type MulticallParameters struct {
	// Contracts is the list of contract calls to execute.
	Contracts []MulticallContract

	// AllowFailure determines whether to continue if individual calls fail.
	// If true, failed calls will be marked with status "failure" but won't
	// stop the entire multicall. Default is true.
	AllowFailure *bool

	// BatchSize is the maximum size in bytes for each batch of calls.
	// Calls are chunked into batches based on their calldata size.
	// Default is 1024 bytes.
	BatchSize int

	// Deployless enables deployless multicall using bytecode execution.
	// This allows multicall on chains without a deployed multicall3 contract.
	Deployless bool

	// MulticallAddress overrides the default multicall3 contract address.
	MulticallAddress *common.Address

	// BlockNumber is the block number to execute the calls at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to execute the calls at.
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// MaxConcurrentChunks limits the number of concurrent chunk executions.
	// This prevents overwhelming RPC endpoints. Default is 4.
	// Set to 0 or negative for unlimited concurrency.
	MaxConcurrentChunks int
}

// MulticallResult represents the result of a single contract call in a multicall.
type MulticallResult struct {
	// Status is either "success" or "failure".
	Status string

	// Result contains the decoded return value(s) if Status is "success".
	Result any

	// Error contains the error if Status is "failure".
	Error error
}

// MulticallReturnType is the return type for the Multicall action.
type MulticallReturnType = []MulticallResult

// call3 represents a single call in the aggregate3 function.
type call3 struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

// aggregate3Result represents the result from aggregate3.
type aggregate3Result struct {
	Success    bool
	ReturnData []byte
}

// chunkResult holds the result of executing a chunk.
type chunkResult struct {
	Results []aggregate3Result
	Err     error
}

// Multicall batches multiple contract function calls into a single RPC call
// using the multicall3 contract.
//
// This is equivalent to viem's `multicall` action.
//
// Example:
//
//	results, err := public.Multicall(ctx, client, public.MulticallParameters{
//	    Contracts: []public.MulticallContract{
//	        {
//	            Address:      tokenAddress,
//	            ABI:          erc20ABI,
//	            FunctionName: "balanceOf",
//	            Args:         []any{userAddress},
//	        },
//	        {
//	            Address:      tokenAddress,
//	            ABI:          erc20ABI,
//	            FunctionName: "totalSupply",
//	        },
//	    },
//	})
func Multicall(ctx context.Context, client Client, params MulticallParameters) (MulticallReturnType, error) {
	// Set defaults
	allowFailure := true
	if params.AllowFailure != nil {
		allowFailure = *params.AllowFailure
	}

	batchSize := params.BatchSize
	if batchSize <= 0 {
		batchSize = 1024
	}

	maxConcurrent := params.MaxConcurrentChunks
	if maxConcurrent <= 0 {
		maxConcurrent = 4
	}

	// Resolve multicall address
	multicallAddress, err := resolveMulticallAddress(client, params)
	if err != nil && !params.Deployless {
		return nil, err
	}

	// Parse ABIs and encode calls
	contracts := params.Contracts
	parsedABIs := make([]*abi.ABI, len(contracts))
	encodedCalls := make([]call3, len(contracts))
	encodeErrors := make([]error, len(contracts))

	for i, contract := range contracts {
		// Parse ABI
		parsedABI, parseErr := parseABIParam(contract.ABI)
		if parseErr != nil {
			encodeErrors[i] = fmt.Errorf("failed to parse ABI for contract %d: %w", i, parseErr)
			encodedCalls[i] = call3{
				Target:       contract.Address,
				AllowFailure: true,
				CallData:     nil,
			}
			continue
		}
		parsedABIs[i] = parsedABI

		// Encode function call
		callData, encodeErr := parsedABI.EncodeFunctionData(contract.FunctionName, contract.Args...)
		if encodeErr != nil {
			encodeErrors[i] = fmt.Errorf("failed to encode call for %q: %w", contract.FunctionName, encodeErr)
			encodedCalls[i] = call3{
				Target:       contract.Address,
				AllowFailure: true,
				CallData:     nil,
			}
			continue
		}

		encodedCalls[i] = call3{
			Target:       contract.Address,
			AllowFailure: true,
			CallData:     callData,
		}
	}

	// Chunk calls by batch size
	chunkedCalls := chunkCalls(encodedCalls, batchSize)

	// Execute chunks in parallel
	chunkResults := make([]*chunkResult, len(chunkedCalls))

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrent)

	for i, chunk := range chunkedCalls {
		i, chunk := i, chunk // capture for goroutine
		g.Go(func() error {
			result, execErr := executeChunk(gCtx, client, chunk, multicallAddress, params)
			chunkResults[i] = &chunkResult{
				Results: result,
				Err:     execErr,
			}
			// Don't return error to continue like Promise.allSettled
			return nil
		})
	}

	if waitErr := g.Wait(); waitErr != nil {
		return nil, fmt.Errorf("multicall execution failed: %w", waitErr)
	}

	// Flatten results and decode
	results := make(MulticallReturnType, len(contracts))
	resultIndex := 0

	for chunkIdx, chunkRes := range chunkResults {
		// Handle chunk-level errors
		if chunkRes.Err != nil {
			if !allowFailure {
				return nil, chunkRes.Err
			}
			// Mark all calls in this chunk as failed
			for j := 0; j < len(chunkedCalls[chunkIdx]); j++ {
				results[resultIndex] = MulticallResult{
					Status: "failure",
					Error:  chunkRes.Err,
				}
				resultIndex++
			}
			continue
		}

		// Process individual results
		for j, aggResult := range chunkRes.Results {
			// Check for encode errors first
			if encodeErrors[resultIndex] != nil {
				if !allowFailure {
					return nil, encodeErrors[resultIndex]
				}
				results[resultIndex] = MulticallResult{
					Status: "failure",
					Error:  encodeErrors[resultIndex],
				}
				resultIndex++
				continue
			}

			contract := contracts[resultIndex]
			parsedABI := parsedABIs[resultIndex]

			// Check if the call succeeded
			if !aggResult.Success {
				err := &RawContractError{Data: aggResult.ReturnData}
				if !allowFailure {
					return nil, err
				}
				results[resultIndex] = MulticallResult{
					Status: "failure",
					Error:  err,
				}
				resultIndex++
				continue
			}

			// Check for empty calldata (encode failed earlier)
			if len(chunkedCalls[chunkIdx][j].CallData) == 0 {
				if !allowFailure {
					return nil, &AbiDecodingZeroDataError{}
				}
				results[resultIndex] = MulticallResult{
					Status: "failure",
					Error:  &AbiDecodingZeroDataError{},
				}
				resultIndex++
				continue
			}

			// Decode the result
			decoded, decodeErr := parsedABI.DecodeFunctionResult(contract.FunctionName, aggResult.ReturnData)
			if decodeErr != nil {
				if !allowFailure {
					return nil, fmt.Errorf("failed to decode result for %q: %w", contract.FunctionName, decodeErr)
				}
				results[resultIndex] = MulticallResult{
					Status: "failure",
					Error:  fmt.Errorf("failed to decode result for %q: %w", contract.FunctionName, decodeErr),
				}
				resultIndex++
				continue
			}

			// Unwrap single return value
			var result any
			if len(decoded) == 1 {
				result = decoded[0]
			} else {
				result = decoded
			}

			results[resultIndex] = MulticallResult{
				Status: "success",
				Result: result,
			}
			resultIndex++
		}
	}

	// Sanity check
	if resultIndex != len(contracts) {
		return nil, fmt.Errorf("multicall results mismatch: got %d, expected %d", resultIndex, len(contracts))
	}

	return results, nil
}

// chunkCalls splits calls into chunks based on batch size.
func chunkCalls(calls []call3, batchSize int) [][]call3 {
	if len(calls) == 0 {
		return nil
	}

	var chunks [][]call3
	currentChunk := []call3{}
	currentSize := 0

	for _, call := range calls {
		callSize := len(call.CallData)
		if callSize == 0 {
			callSize = 2 // "0x" placeholder
		}

		// Check if we need a new chunk
		if batchSize > 0 && currentSize+callSize > batchSize && len(currentChunk) > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = []call3{}
			currentSize = 0
		}

		currentChunk = append(currentChunk, call)
		currentSize += callSize
	}

	// Add final chunk
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// executeChunk executes a single chunk of calls via multicall3.
func executeChunk(ctx context.Context, client Client, calls []call3, multicallAddress *common.Address, params MulticallParameters) ([]aggregate3Result, error) {
	// Encode aggregate3 call
	calldata, err := encodeAggregate3(calls)
	if err != nil {
		return nil, fmt.Errorf("failed to encode aggregate3: %w", err)
	}

	// Build call request
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	var req callRequest
	var rpcParams []any

	if params.Deployless || multicallAddress == nil {
		// Deployless multicall - wrap in deployless bytecode
		deploylessData, deploylessErr := deployless.ToDeploylessCallViaBytecodeData(
			common.FromHex(constants.Multicall3Bytecode),
			calldata,
		)
		if deploylessErr != nil {
			return nil, fmt.Errorf("failed to encode deployless multicall: %w", deploylessErr)
		}
		req = callRequest{Data: hexutil.Encode(deploylessData)}
	} else {
		req = callRequest{
			To:   multicallAddress.Hex(),
			Data: hexutil.Encode(calldata),
		}
	}

	rpcParams = []any{req, blockTag}

	// Execute call
	resp, requestErr := client.Request(ctx, "eth_call", rpcParams...)
	if requestErr != nil {
		return nil, fmt.Errorf("eth_call failed: %w", requestErr)
	}

	var hexResult string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", unmarshalErr)
	}

	// Decode aggregate3 result
	resultData := common.FromHex(hexResult)
	return decodeAggregate3Result(resultData)
}

// encodeAggregate3 encodes calls for the aggregate3 function.
func encodeAggregate3(calls []call3) ([]byte, error) {
	// Build tuple array for encoding
	callTuples := make([]any, len(calls))
	for i, call := range calls {
		callTuples[i] = []any{call.Target, call.AllowFailure, call.CallData}
	}

	// Encode using ABI parameters
	encoded, err := abi.EncodeAbiParameters(
		[]abi.AbiParam{
			{
				Type: "tuple[]",
				Components: []abi.AbiParam{
					{Name: "target", Type: "address"},
					{Name: "allowFailure", Type: "bool"},
					{Name: "callData", Type: "bytes"},
				},
			},
		},
		[]any{callTuples},
	)
	if err != nil {
		return nil, err
	}

	// Prepend aggregate3 selector (0x82ad56cb)
	selector := common.FromHex(constants.Aggregate3Signature)
	result := make([]byte, len(selector)+len(encoded))
	copy(result, selector)
	copy(result[len(selector):], encoded)

	return result, nil
}

// decodeAggregate3Result decodes the result from aggregate3.
func decodeAggregate3Result(data []byte) ([]aggregate3Result, error) {
	// aggregate3 returns: (bool success, bytes returnData)[]
	decoded, err := abi.DecodeAbiParameters(
		[]abi.AbiParam{
			{
				Type: "tuple[]",
				Components: []abi.AbiParam{
					{Name: "success", Type: "bool"},
					{Name: "returnData", Type: "bytes"},
				},
			},
		},
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decode aggregate3 result: %w", err)
	}

	if len(decoded) == 0 {
		return nil, fmt.Errorf("empty aggregate3 result")
	}

	// Extract results
	resultsRaw, ok := decoded[0].([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected aggregate3 result type: %T", decoded[0])
	}

	results := make([]aggregate3Result, len(resultsRaw))
	for i, r := range resultsRaw {
		tuple, ok := r.([]any)
		if !ok || len(tuple) < 2 {
			return nil, fmt.Errorf("invalid aggregate3 result tuple at index %d", i)
		}

		success, ok := tuple[0].(bool)
		if !ok {
			return nil, fmt.Errorf("invalid success value at index %d", i)
		}

		returnData, ok := tuple[1].([]byte)
		if !ok {
			// Try to handle nil or empty
			if tuple[1] == nil {
				returnData = nil
			} else {
				return nil, fmt.Errorf("invalid returnData at index %d: %T", i, tuple[1])
			}
		}

		results[i] = aggregate3Result{
			Success:    success,
			ReturnData: returnData,
		}
	}

	return results, nil
}

// resolveMulticallAddress determines the multicall3 contract address.
func resolveMulticallAddress(client Client, params MulticallParameters) (*common.Address, error) {
	// Use provided address if specified
	if params.MulticallAddress != nil {
		return params.MulticallAddress, nil
	}

	// Deployless doesn't need an address
	if params.Deployless {
		return nil, nil
	}

	// Get from chain config
	chain := client.Chain()
	if chain == nil {
		return nil, &ChainNotConfiguredError{}
	}

	if chain.Contracts == nil || chain.Contracts.Multicall3 == nil {
		return nil, &ChainDoesNotSupportContractError{
			ChainID:      chain.ID,
			ContractName: "multicall3",
		}
	}

	// Check block number constraint
	if params.BlockNumber != nil && chain.Contracts.Multicall3.BlockCreated != nil {
		if *params.BlockNumber < *chain.Contracts.Multicall3.BlockCreated {
			return nil, &ChainDoesNotSupportContractError{
				ChainID:      chain.ID,
				ContractName: "multicall3",
				BlockNumber:  params.BlockNumber,
			}
		}
	}

	return &chain.Contracts.Multicall3.Address, nil
}

// parseABIParam parses the ABI parameter which can be []byte, string, or *abi.ABI.
func parseABIParam(abiParam any) (*abi.ABI, error) {
	switch v := abiParam.(type) {
	case *abi.ABI:
		return v, nil
	case []byte:
		return abi.Parse(v)
	case string:
		return abi.Parse([]byte(v))
	default:
		return nil, fmt.Errorf("ABI must be []byte, string, or *abi.ABI, got %T", abiParam)
	}
}

// AbiDecodingZeroDataError is returned when trying to decode zero data.
type AbiDecodingZeroDataError struct{}

func (e *AbiDecodingZeroDataError) Error() string {
	return "cannot decode zero data (0x) - the function may have reverted"
}

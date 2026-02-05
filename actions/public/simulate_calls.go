package public

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// SimulateCall represents a single call to simulate.
type SimulateCall struct {
	// From is the sender address.
	From *common.Address

	// To is the contract address to call.
	To *common.Address

	// Data is the calldata to send.
	Data []byte

	// Value is the amount of wei to send with the call.
	Value *big.Int

	// Gas is the gas limit for the call.
	Gas *uint64

	// GasPrice is the gas price for the call (legacy).
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int

	// Nonce is the nonce for the call.
	Nonce *uint64

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// ABI is the contract ABI for encoding/decoding.
	ABI *abi.ABI

	// FunctionName is the function name for ABI encoding/decoding.
	FunctionName string

	// Args are the function arguments.
	Args []any
}

// SimulateCallsParameters contains the parameters for the SimulateCalls action.
type SimulateCallsParameters struct {
	// Account is the account attached to the calls (msg.sender).
	Account *common.Address

	// Calls are the calls to simulate.
	Calls []SimulateCall

	// BlockNumber is the block number to simulate at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to simulate at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// StateOverrides contains state overrides for the simulation.
	StateOverrides types.StateOverride

	// TraceAssetChanges enables tracing of asset (ETH/ERC20) balance changes.
	TraceAssetChanges bool

	// TraceTransfers enables transfer tracing.
	TraceTransfers bool

	// Validation enables validation mode.
	Validation bool
}

// AssetChange represents a change in an asset balance.
type AssetChange struct {
	// Token contains information about the token.
	Token TokenInfo

	// Value contains the pre/post/diff balance values.
	Value BalanceValue
}

// TokenInfo contains information about a token.
type TokenInfo struct {
	// Address is the token contract address (ETH uses a special address).
	Address common.Address

	// Decimals is the token decimals (if available).
	Decimals *int

	// Symbol is the token symbol (if available).
	Symbol string
}

// BalanceValue contains pre/post balance values.
type BalanceValue struct {
	// Pre is the balance before the simulation.
	Pre *big.Int

	// Post is the balance after the simulation.
	Post *big.Int

	// Diff is the difference (Post - Pre).
	Diff *big.Int
}

// SimulateCallResult represents the result of a single call.
type SimulateCallResult struct {
	// Status is the result status ("success" or "failure").
	Status string

	// Data is the return data from the call.
	Data []byte

	// GasUsed is the amount of gas used by the call.
	GasUsed *big.Int

	// Logs are the logs emitted by the call.
	Logs []formatters.Log

	// Result is the decoded result (if ABI was provided).
	Result any

	// Error is the error if the call failed.
	Error error
}

// SimulateCallsReturnType is the return type for the SimulateCalls action.
type SimulateCallsReturnType struct {
	// AssetChanges contains the asset balance changes (if TraceAssetChanges was enabled).
	AssetChanges []AssetChange

	// Block is the simulated block data.
	Block formatters.Block

	// Results are the results of each call.
	Results []SimulateCallResult
}

// ETH address constant (used for representing native ETH in asset changes).
var ETHAddress = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")

// zeroAddress is the zero address.
var zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// getBalanceCode is bytecode for a contract that returns ETH balance.
// This is used to get ETH balance via contract call in simulation.
const getBalanceCode = "0x6080604052348015600e575f80fd5b5061016d8061001c5f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c8063f8b2cb4f1461002d575b5f80fd5b610047600480360381019061004291906100db565b61005d565b604051610054919061011e565b60405180910390f35b5f8173ffffffffffffffffffffffffffffffffffffffff16319050919050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100aa82610081565b9050919050565b6100ba816100a0565b81146100c4575f80fd5b50565b5f813590506100d5816100b1565b92915050565b5f602082840312156100f0576100ef61007d565b5b5f6100fd848285016100c7565b91505092915050565b5f819050919050565b61011881610106565b82525050565b5f6020820190506101315f83018461010f565b9291505056fea26469706673582212203b9fe929fe995c7cf9887f0bdba8a36dd78e8b73f149b17d2d9ad7cd09d2dc6264736f6c634300081a0033"

// SimulateCalls simulates execution of a batch of calls.
//
// This is equivalent to viem's `simulateCalls` action.
// It internally uses eth_simulateV1 to execute the calls.
//
// When TraceAssetChanges is enabled, it also tracks ETH and ERC20 balance changes
// by making additional simulation calls.
//
// Example:
//
//	result, err := public.SimulateCalls(ctx, client, public.SimulateCallsParameters{
//	    Account: &senderAddr,
//	    Calls: []public.SimulateCall{
//	        {To: &contractAddr, Data: calldata1},
//	        {To: &contractAddr, Data: calldata2},
//	    },
//	})
func SimulateCalls(ctx context.Context, client Client, params SimulateCallsParameters) (*SimulateCallsReturnType, error) {
	if params.TraceAssetChanges && params.Account == nil {
		return nil, fmt.Errorf("`account` is required when `traceAssetChanges` is true")
	}

	// If not tracing asset changes, use simple simulation
	if !params.TraceAssetChanges {
		return simulateCallsSimple(ctx, client, params)
	}

	// With asset tracing, we need to do more complex simulation
	return simulateCallsWithAssetTracing(ctx, client, params)
}

// simulateCallsSimple performs a simple simulation without asset tracing.
func simulateCallsSimple(ctx context.Context, client Client, params SimulateCallsParameters) (*SimulateCallsReturnType, error) {
	// Build calls for simulation
	simCalls := make([]SimulateBlockCall, len(params.Calls))
	for i, call := range params.Calls {
		from := call.From
		if from == nil && params.Account != nil {
			from = params.Account
		}

		// Handle data encoding
		data := call.Data
		if len(data) == 0 && call.ABI != nil && call.FunctionName != "" {
			encoded, err := call.ABI.EncodeFunctionData(call.FunctionName, call.Args...)
			if err != nil {
				return nil, fmt.Errorf("failed to encode function data for call %d: %w", i, err)
			}
			data = encoded
		}

		simCalls[i] = SimulateBlockCall{
			From:                 from,
			To:                   call.To,
			Data:                 data,
			Value:                call.Value,
			Gas:                  call.Gas,
			GasPrice:             call.GasPrice,
			MaxFeePerGas:         call.MaxFeePerGas,
			MaxPriorityFeePerGas: call.MaxPriorityFeePerGas,
			Nonce:                call.Nonce,
			AccessList:           call.AccessList,
			ABI:                  call.ABI,
			FunctionName:         call.FunctionName,
			Args:                 call.Args,
		}
	}

	// Add an empty call at the end (matches viem behavior)
	simCalls = append(simCalls, SimulateBlockCall{})

	// Execute simulation
	blocks, err := SimulateBlocks(ctx, client, SimulateBlocksParameters{
		Blocks: []SimulateBlock{{
			Calls:          simCalls,
			StateOverrides: params.StateOverrides,
		}},
		BlockNumber:    params.BlockNumber,
		BlockTag:       params.BlockTag,
		TraceTransfers: params.TraceTransfers,
		Validation:     params.Validation,
	})
	if err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("no block results returned from simulation")
	}

	blockResult := blocks[0]

	// Extract results (excluding the last empty call)
	results := make([]SimulateCallResult, 0, len(blockResult.Calls)-1)
	for i := 0; i < len(blockResult.Calls)-1 && i < len(params.Calls); i++ {
		callResult := blockResult.Calls[i]
		results = append(results, SimulateCallResult(callResult))
	}

	return &SimulateCallsReturnType{
		AssetChanges: nil,
		Block:        blockResult.Block,
		Results:      results,
	}, nil
}

// simulateCallsWithAssetTracing performs simulation with ETH/ERC20 balance tracking.
func simulateCallsWithAssetTracing(ctx context.Context, client Client, params SimulateCallsParameters) (*SimulateCallsReturnType, error) {
	account := params.Account

	// First, fetch ERC20/721 addresses that were "touched" from the calls via createAccessList
	assetAddresses := []common.Address{}
	for _, call := range params.Calls {
		if len(call.Data) == 0 && call.ABI == nil {
			continue
		}

		data := call.Data
		if len(data) == 0 && call.ABI != nil && call.FunctionName != "" {
			encoded, err := call.ABI.EncodeFunctionData(call.FunctionName, call.Args...)
			if err != nil {
				continue // Skip if encoding fails
			}
			data = encoded
		}

		accessListResult, err := CreateAccessList(ctx, client, CreateAccessListParameters{
			Account: account,
			To:      call.To,
			Data:    data,
		})
		if err != nil {
			continue // Skip if access list creation fails
		}

		for _, item := range accessListResult.AccessList {
			if len(item.StorageKeys) > 0 {
				assetAddresses = append(assetAddresses, item.Address)
			}
		}
	}

	// Remove duplicates
	assetAddresses = uniqueAddresses(assetAddresses)

	// Build the getBalance ABI for ETH balance queries
	getBalanceABI, err := abi.ParseFromString(`[{"inputs":[{"type":"address"}],"name":"getBalance","outputs":[{"type":"uint256"}],"type":"function"}]`)
	if err != nil {
		return nil, fmt.Errorf("failed to parse getBalance ABI: %w", err)
	}

	// Build balanceOf ABI for ERC20 balance queries
	balanceOfABI, err := abi.ParseFromString(`[{"inputs":[{"type":"address"}],"name":"balanceOf","outputs":[{"type":"uint256"}],"type":"function"}]`)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balanceOf ABI: %w", err)
	}

	// Build decimals/symbol ABIs
	decimalsABI, err := abi.ParseFromString(`[{"inputs":[],"name":"decimals","outputs":[{"type":"uint256"}],"type":"function"}]`)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decimals ABI: %w", err)
	}

	symbolABI, err := abi.ParseFromString(`[{"inputs":[],"name":"symbol","outputs":[{"type":"string"}],"type":"function"}]`)
	if err != nil {
		return nil, fmt.Errorf("failed to parse symbol ABI: %w", err)
	}

	// Build all simulation blocks
	blocks := []SimulateBlock{}

	// Block 1: ETH pre balance
	getBalanceData, err := getBalanceABI.EncodeFunctionData("getBalance", *account)
	if err != nil {
		return nil, fmt.Errorf("failed to encode getBalance call: %w", err)
	}
	blocks = append(blocks, SimulateBlock{
		Calls: []SimulateBlockCall{{
			Data: common.FromHex(getBalanceCode + common.Bytes2Hex(getBalanceData)),
		}},
		StateOverrides: params.StateOverrides,
	})

	// Block 2: Asset pre balances
	if len(assetAddresses) > 0 {
		assetPreCalls := make([]SimulateBlockCall, len(assetAddresses))
		for i, addr := range assetAddresses {
			nonce := uint64(i)
			addrCopy := addr
			assetPreCalls[i] = SimulateBlockCall{
				To:           &addrCopy,
				ABI:          balanceOfABI,
				FunctionName: "balanceOf",
				Args:         []any{*account},
				From:         &zeroAddress,
				Nonce:        &nonce,
			}
		}
		blocks = append(blocks, SimulateBlock{
			Calls: assetPreCalls,
			StateOverrides: types.StateOverride{
				zeroAddress: types.StateOverrideAccount{
					Nonce: ptr(uint64(0)),
				},
			},
		})
	}

	// Block 3: Main calls
	mainCalls := make([]SimulateBlockCall, len(params.Calls)+1)
	for i, call := range params.Calls {
		from := call.From
		if from == nil && params.Account != nil {
			from = params.Account
		}

		data := call.Data
		if len(data) == 0 && call.ABI != nil && call.FunctionName != "" {
			encoded, _ := call.ABI.EncodeFunctionData(call.FunctionName, call.Args...)
			data = encoded
		}

		mainCalls[i] = SimulateBlockCall{
			From:                 from,
			To:                   call.To,
			Data:                 data,
			Value:                call.Value,
			Gas:                  call.Gas,
			GasPrice:             call.GasPrice,
			MaxFeePerGas:         call.MaxFeePerGas,
			MaxPriorityFeePerGas: call.MaxPriorityFeePerGas,
			Nonce:                call.Nonce,
			AccessList:           call.AccessList,
			ABI:                  call.ABI,
			FunctionName:         call.FunctionName,
			Args:                 call.Args,
		}
	}
	mainCalls[len(params.Calls)] = SimulateBlockCall{} // Empty call at end
	blocks = append(blocks, SimulateBlock{
		Calls:          mainCalls,
		StateOverrides: params.StateOverrides,
	})

	// Block 4: ETH post balance
	blocks = append(blocks, SimulateBlock{
		Calls: []SimulateBlockCall{{
			Data: common.FromHex(getBalanceCode + common.Bytes2Hex(getBalanceData)),
		}},
	})

	// Block 5: Asset post balances
	if len(assetAddresses) > 0 {
		assetPostCalls := make([]SimulateBlockCall, len(assetAddresses))
		for i, addr := range assetAddresses {
			nonce := uint64(i)
			addrCopy := addr
			assetPostCalls[i] = SimulateBlockCall{
				To:           &addrCopy,
				ABI:          balanceOfABI,
				FunctionName: "balanceOf",
				Args:         []any{*account},
				From:         &zeroAddress,
				Nonce:        &nonce,
			}
		}
		blocks = append(blocks, SimulateBlock{
			Calls: assetPostCalls,
			StateOverrides: types.StateOverride{
				zeroAddress: types.StateOverrideAccount{
					Nonce: ptr(uint64(0)),
				},
			},
		})
	}

	// Block 6: Decimals
	if len(assetAddresses) > 0 {
		decimalsCalls := make([]SimulateBlockCall, len(assetAddresses))
		for i, addr := range assetAddresses {
			nonce := uint64(i)
			addrCopy := addr
			decimalsCalls[i] = SimulateBlockCall{
				To:           &addrCopy,
				ABI:          decimalsABI,
				FunctionName: "decimals",
				From:         &zeroAddress,
				Nonce:        &nonce,
			}
		}
		blocks = append(blocks, SimulateBlock{
			Calls: decimalsCalls,
			StateOverrides: types.StateOverride{
				zeroAddress: types.StateOverrideAccount{
					Nonce: ptr(uint64(0)),
				},
			},
		})
	}

	// Block 7: Symbols
	if len(assetAddresses) > 0 {
		symbolCalls := make([]SimulateBlockCall, len(assetAddresses))
		for i, addr := range assetAddresses {
			nonce := uint64(i)
			addrCopy := addr
			symbolCalls[i] = SimulateBlockCall{
				To:           &addrCopy,
				ABI:          symbolABI,
				FunctionName: "symbol",
				From:         &zeroAddress,
				Nonce:        &nonce,
			}
		}
		blocks = append(blocks, SimulateBlock{
			Calls: symbolCalls,
			StateOverrides: types.StateOverride{
				zeroAddress: types.StateOverrideAccount{
					Nonce: ptr(uint64(0)),
				},
			},
		})
	}

	// Execute simulation
	blockResults, err := SimulateBlocks(ctx, client, SimulateBlocksParameters{
		Blocks:         blocks,
		BlockNumber:    params.BlockNumber,
		BlockTag:       params.BlockTag,
		TraceTransfers: params.TraceTransfers,
		Validation:     params.Validation,
	})
	if err != nil {
		return nil, err
	}

	// Parse results
	// Block indices depend on whether we have asset addresses
	var (
		ethPreBlock    BlockResult
		assetPreBlock  *BlockResult
		mainBlock      BlockResult
		ethPostBlock   BlockResult
		assetPostBlock *BlockResult
		decimalsBlock  *BlockResult
		symbolsBlock   *BlockResult
	)

	blockIdx := 0
	ethPreBlock = blockResults[blockIdx]
	blockIdx++

	if len(assetAddresses) > 0 {
		assetPreBlock = &blockResults[blockIdx]
		blockIdx++
	}

	mainBlock = blockResults[blockIdx]
	blockIdx++

	ethPostBlock = blockResults[blockIdx]
	blockIdx++

	if len(assetAddresses) > 0 {
		assetPostBlock = &blockResults[blockIdx]
		blockIdx++
		decimalsBlock = &blockResults[blockIdx]
		blockIdx++
		symbolsBlock = &blockResults[blockIdx]
	}

	// Extract main results (excluding the last empty call)
	results := make([]SimulateCallResult, 0, len(mainBlock.Calls)-1)
	for i := 0; i < len(mainBlock.Calls)-1 && i < len(params.Calls); i++ {
		callResult := mainBlock.Calls[i]
		results = append(results, SimulateCallResult(callResult))
	}

	// Build asset changes
	assetChanges := []AssetChange{}

	// ETH balance change
	if len(ethPreBlock.Calls) > 0 && len(ethPostBlock.Calls) > 0 {
		ethPre := extractBalance(ethPreBlock.Calls[0])
		ethPost := extractBalance(ethPostBlock.Calls[0])

		if ethPre != nil && ethPost != nil {
			diff := new(big.Int).Sub(ethPost, ethPre)
			decimals := 18
			assetChanges = append(assetChanges, AssetChange{
				Token: TokenInfo{
					Address:  ETHAddress,
					Decimals: &decimals,
					Symbol:   "ETH",
				},
				Value: BalanceValue{
					Pre:  ethPre,
					Post: ethPost,
					Diff: diff,
				},
			})
		}
	}

	// Asset balance changes
	if assetPreBlock != nil && assetPostBlock != nil {
		for i, addr := range assetAddresses {
			if i >= len(assetPreBlock.Calls) || i >= len(assetPostBlock.Calls) {
				continue
			}

			pre := extractBalance(assetPreBlock.Calls[i])
			post := extractBalance(assetPostBlock.Calls[i])

			if pre == nil || post == nil {
				continue
			}

			// Check if this address is already in asset changes
			alreadyAdded := false
			for _, change := range assetChanges {
				if change.Token.Address == addr {
					alreadyAdded = true
					break
				}
			}
			if alreadyAdded {
				continue
			}

			diff := new(big.Int).Sub(post, pre)

			// Get decimals and symbol
			var decimals *int
			var symbol string

			if decimalsBlock != nil && i < len(decimalsBlock.Calls) {
				if decimalsBlock.Calls[i].Status == "success" && decimalsBlock.Calls[i].Result != nil {
					if d, ok := decimalsBlock.Calls[i].Result.(*big.Int); ok {
						dec := int(d.Int64())
						decimals = &dec
					}
				}
			}

			if symbolsBlock != nil && i < len(symbolsBlock.Calls) {
				if symbolsBlock.Calls[i].Status == "success" && symbolsBlock.Calls[i].Result != nil {
					if s, ok := symbolsBlock.Calls[i].Result.(string); ok {
						symbol = s
					}
				}
			}

			assetChanges = append(assetChanges, AssetChange{
				Token: TokenInfo{
					Address:  addr,
					Decimals: decimals,
					Symbol:   symbol,
				},
				Value: BalanceValue{
					Pre:  pre,
					Post: post,
					Diff: diff,
				},
			})
		}
	}

	return &SimulateCallsReturnType{
		AssetChanges: assetChanges,
		Block:        mainBlock.Block,
		Results:      results,
	}, nil
}

// extractBalance extracts a balance value from a call result.
func extractBalance(result CallResult) *big.Int {
	if result.Status != "success" {
		return nil
	}

	// Try to get from decoded result
	if result.Result != nil {
		if b, ok := result.Result.(*big.Int); ok {
			return b
		}
	}

	// Try to decode from raw data
	if len(result.Data) >= 32 {
		return new(big.Int).SetBytes(result.Data[:32])
	}

	return nil
}

// uniqueAddresses returns a slice with duplicate addresses removed.
func uniqueAddresses(addrs []common.Address) []common.Address {
	seen := make(map[common.Address]struct{})
	result := make([]common.Address, 0, len(addrs))

	for _, addr := range addrs {
		if _, ok := seen[addr]; !ok {
			seen[addr] = struct{}{}
			result = append(result, addr)
		}
	}

	return result
}

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

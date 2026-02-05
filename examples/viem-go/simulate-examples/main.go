// Simulate Examples - viem-go
//
// Comprehensive examples demonstrating simulation actions:
// - SimulateBlocks: Simulate multiple blocks with state/block overrides
// - CreateAccessList: Generate EIP-2930 access lists
// - SimulateCalls: Simulate batch calls with asset change tracking
// - SimulateContract: Simulate contract write functions
//
// NOTE: eth_simulateV1 (used by SimulateBlocks and SimulateCalls) is a relatively
// new RPC method that requires specific node support. Not all RPC providers support it.
// For full support, use:
// - Alchemy (mainnet, sepolia, etc.)
// - Your own node with recent geth/nethermind
// - Tenderly
//
// CreateAccessList and SimulateContract use standard eth_createAccessList and eth_call
// which are supported by most RPC providers.
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// Example addresses (Ethereum Mainnet)
var (
	// USDC contract on Ethereum Mainnet
	usdcAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	// WETH contract on Ethereum Mainnet
	wethAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	// Vitalik's address
	vitalikAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// Test addresses
	testAddress   = common.HexToAddress("0x1234567890123456789012345678901234567890")
	testRecipient = common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
)

// ERC20 ABI for token interactions
const erc20ABIJson = `[
	{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"type":"bool"}],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[],"name":"name","outputs":[{"type":"string"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"symbol","outputs":[{"type":"string"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"decimals","outputs":[{"type":"uint8"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"totalSupply","outputs":[{"type":"uint256"}],"stateMutability":"view","type":"function"}
]`

func main() {
	ctx := context.Background()

	printHeader("Simulation Action Examples (viem-go)")

	// Get RPC URL from environment or use default
	// NOTE: eth_simulateV1 requires specific RPC support (Alchemy, Tenderly, or recent geth/nethermind)
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://eth.llamarpc.com"
		fmt.Println("TIP: Set ETH_RPC_URL env var for a specific RPC (Alchemy recommended for eth_simulateV1)")
	}

	// Create Public Client
	printSection("1. Creating Public Client")
	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Mainnet,
		Transport: transport.HTTP(rpcURL),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()
	fmt.Printf("Connected to Ethereum Mainnet via %s\n", truncateURL(rpcURL))

	// Parse ERC20 ABI
	erc20ABI, err := abi.ParseFromString(erc20ABIJson)
	if err != nil {
		fmt.Printf("Error parsing ABI: %v\n", err)
		return
	}

	// // ============================================================================
	// // CreateAccessList Examples
	// // ============================================================================
	// printHeader("CreateAccessList Examples")

	// // Example 1: Basic CreateAccessList
	// printSection("2. Basic CreateAccessList")
	balanceOfData, _ := erc20ABI.EncodeFunctionData("balanceOf", vitalikAddress)
	accessListResult, err := public.CreateAccessList(ctx, publicClient, public.CreateAccessListParameters{
		Account: &vitalikAddress,
		To:      &usdcAddress,
		Data:    balanceOfData,
	})
	if err != nil {
		fmt.Printf("Error creating access list: %v\n", err)
	} else {
		fmt.Printf("Access list created successfully!\n")
		fmt.Printf("  Gas used estimate: %s\n", accessListResult.GasUsed.String())
		fmt.Printf("  Access list entries: %d\n", len(accessListResult.AccessList))
		for i, entry := range accessListResult.AccessList {
			fmt.Printf("    [%d] Address: %s\n", i, truncateAddress(entry.Address))
			fmt.Printf("        Storage keys: %d\n", len(entry.StorageKeys))
		}
	}

	// Example 2: CreateAccessList for transfer
	printSection("3. CreateAccessList for Transfer")
	transferData, _ := erc20ABI.EncodeFunctionData("transfer", testRecipient, big.NewInt(1000000))
	accessListResult, err = public.CreateAccessList(ctx, publicClient, public.CreateAccessListParameters{
		Account: &vitalikAddress,
		To:      &usdcAddress,
		Data:    transferData,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Transfer access list created!\n")
		fmt.Printf("  Gas estimate: %s\n", accessListResult.GasUsed.String())
		fmt.Printf("  Touched addresses: %d\n", len(accessListResult.AccessList))
	}

	// ============================================================================
	// SimulateContract Examples
	// ============================================================================
	printHeader("SimulateContract Examples")

	// Example 3: Simulate ERC20 Transfer
	printSection("4. Simulate ERC20 Transfer")
	simResult, err := public.SimulateContract(ctx, publicClient, public.SimulateContractParameters{
		Account:      &vitalikAddress,
		Address:      usdcAddress,
		ABI:          erc20ABI,
		FunctionName: "transfer",
		Args:         []any{testRecipient, big.NewInt(1000000)}, // 1 USDC (6 decimals)
	})
	if err != nil {
		fmt.Printf("Simulation failed: %v\n", err)
	} else {
		fmt.Printf("Transfer simulation successful!\n")
		if success, ok := simResult.Result.(bool); ok {
			fmt.Printf("  Would succeed: %v\n", success)
		}
		fmt.Printf("  Request ready for writeContract:\n")
		fmt.Printf("    Address: %s\n", simResult.Request.Address.Hex())
		fmt.Printf("    Function: %s\n", simResult.Request.FunctionName)
	}

	// Example 4: Simulate with State Override (give test address tokens)
	printSection("5. Simulate with State Override")
	// Override testAddress to have a large token balance
	tokenBalanceSlot := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	largeBalance := common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240") // 1M
	stateOverride := types.StateOverride{
		usdcAddress: types.StateOverrideAccount{
			StateDiff: types.StateMapping{
				{Slot: tokenBalanceSlot, Value: largeBalance},
			},
		},
	}
	simResult, err = public.SimulateContract(ctx, publicClient, public.SimulateContractParameters{
		Account:       &testAddress,
		Address:       usdcAddress,
		ABI:           erc20ABI,
		FunctionName:  "balanceOf",
		Args:          []any{testAddress},
		StateOverride: stateOverride,
	})
	if err != nil {
		fmt.Printf("Simulation with override failed: %v\n", err)
	} else {
		fmt.Printf("State override simulation successful!\n")
		if balance, ok := simResult.Result.(*big.Int); ok {
			fmt.Printf("  Overridden balance: %s\n", balance.String())
		}
	}

	// Example 5: Read contract name via SimulateContract
	printSection("6. Simulate Contract Read (name)")
	simResult, err = public.SimulateContract(ctx, publicClient, public.SimulateContractParameters{
		Address:      usdcAddress,
		ABI:          erc20ABI,
		FunctionName: "name",
		Args:         []any{},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		if name, ok := simResult.Result.(string); ok {
			fmt.Printf("Contract name: %s\n", name)
		}
	}

	// ============================================================================
	// SimulateBlocks Examples
	// ============================================================================
	printHeader("SimulateBlocks Examples")

	// Example 6: Basic SimulateBlocks
	printSection("7. Basic SimulateBlocks - Single Block")
	blocksResult, err := public.SimulateBlocks(ctx, publicClient, public.SimulateBlocksParameters{
		Blocks: []public.SimulateBlock{
			{
				Calls: []public.SimulateBlockCall{
					{
						From:         &vitalikAddress,
						To:           &usdcAddress,
						Data:         balanceOfData,
						ABI:          erc20ABI,
						FunctionName: "balanceOf",
						Args:         []any{vitalikAddress},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("SimulateBlocks error: %v\n", err)
	} else {
		fmt.Printf("Block simulation successful!\n")
		fmt.Printf("  Blocks simulated: %d\n", len(blocksResult))
		if len(blocksResult) > 0 && len(blocksResult[0].Calls) > 0 {
			call := blocksResult[0].Calls[0]
			fmt.Printf("  Call status: %s\n", call.Status)
			fmt.Printf("  Gas used: %s\n", call.GasUsed.String())
			if call.Result != nil {
				if balance, ok := call.Result.(*big.Int); ok {
					fmt.Printf("  Balance: %s\n", balance.String())
				}
			}
		}
	}

	// Example 7: SimulateBlocks with Block Overrides
	printSection("8. SimulateBlocks with Block Overrides")
	futureBlockNum := uint64(82584814) // Far future block
	blocksResult, err = public.SimulateBlocks(ctx, publicClient, public.SimulateBlocksParameters{
		Blocks: []public.SimulateBlock{
			{
				BlockOverrides: &types.BlockOverrides{
					Number: &futureBlockNum,
				},
				Calls: []public.SimulateBlockCall{
					{
						To:   &usdcAddress,
						Data: balanceOfData,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Block override simulation error: %v\n", err)
	} else {
		fmt.Printf("Block override simulation successful!\n")
		if len(blocksResult) > 0 {
			fmt.Printf("  Simulated block number: %s\n", blocksResult[0].Number.String())
			fmt.Printf("  Simulated timestamp: %s\n", blocksResult[0].Timestamp.String())
		}
	}

	// Example 8: SimulateBlocks with Multiple Blocks
	printSection("9. SimulateBlocks - Multiple Sequential Blocks")
	blocksResult, err = public.SimulateBlocks(ctx, publicClient, public.SimulateBlocksParameters{
		Blocks: []public.SimulateBlock{
			{
				Calls: []public.SimulateBlockCall{
					{To: &usdcAddress, Data: balanceOfData},
				},
			},
			{
				Calls: []public.SimulateBlockCall{
					{To: &wethAddress, Data: balanceOfData},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("Multi-block simulation error: %v\n", err)
	} else {
		fmt.Printf("Multi-block simulation successful!\n")
		fmt.Printf("  Blocks simulated: %d\n", len(blocksResult))
		for i, block := range blocksResult {
			fmt.Printf("  Block %d: %d calls\n", i+1, len(block.Calls))
		}
	}

	// Example 9: SimulateBlocks with State Override
	printSection("10. SimulateBlocks with State Override")
	overrideBalance := mustParseEther("1000")
	blocksResult, err = public.SimulateBlocks(ctx, publicClient, public.SimulateBlocksParameters{
		Blocks: []public.SimulateBlock{
			{
				StateOverrides: types.StateOverride{
					testAddress: types.StateOverrideAccount{
						Balance: overrideBalance,
					},
				},
				Calls: []public.SimulateBlockCall{
					{
						From:  &testAddress,
						To:    &testRecipient,
						Value: mustParseEther("100"),
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Printf("State override simulation error: %v\n", err)
	} else {
		fmt.Printf("State override block simulation successful!\n")
		if len(blocksResult) > 0 && len(blocksResult[0].Calls) > 0 {
			fmt.Printf("  ETH transfer status: %s\n", blocksResult[0].Calls[0].Status)
		}
	}

	// ============================================================================
	// SimulateCalls Examples
	// ============================================================================
	printHeader("SimulateCalls Examples")

	// Example 10: Basic SimulateCalls
	printSection("11. Basic SimulateCalls - Batch Operations")
	nameData, _ := erc20ABI.EncodeFunctionData("name")
	symbolData, _ := erc20ABI.EncodeFunctionData("symbol")
	decimalsData, _ := erc20ABI.EncodeFunctionData("decimals")

	callsResult, err := public.SimulateCalls(ctx, publicClient, public.SimulateCallsParameters{
		Account: &vitalikAddress,
		Calls: []public.SimulateCall{
			{To: &usdcAddress, Data: nameData},
			{To: &usdcAddress, Data: symbolData},
			{To: &usdcAddress, Data: decimalsData},
		},
	})
	if err != nil {
		fmt.Printf("SimulateCalls error: %v\n", err)
	} else {
		fmt.Printf("Batch simulation successful!\n")
		fmt.Printf("  Calls simulated: %d\n", len(callsResult.Results))
		for i, result := range callsResult.Results {
			fmt.Printf("  [%d] Status: %s, Gas: %s\n", i, result.Status, result.GasUsed.String())
		}
	}

	// Example 11: SimulateCalls with ABI Decoding
	printSection("12. SimulateCalls with ABI Decoding")
	callsResult, err = public.SimulateCalls(ctx, publicClient, public.SimulateCallsParameters{
		Account: &vitalikAddress,
		Calls: []public.SimulateCall{
			{
				To:           &usdcAddress,
				ABI:          erc20ABI,
				FunctionName: "balanceOf",
				Args:         []any{vitalikAddress},
			},
			{
				To:           &usdcAddress,
				ABI:          erc20ABI,
				FunctionName: "totalSupply",
				Args:         []any{},
			},
		},
	})
	if err != nil {
		fmt.Printf("SimulateCalls error: %v\n", err)
	} else {
		fmt.Printf("Batch simulation with decoding successful!\n")
		for i, result := range callsResult.Results {
			fmt.Printf("  [%d] Status: %s\n", i, result.Status)
			if result.Result != nil {
				if val, ok := result.Result.(*big.Int); ok {
					fmt.Printf("      Value: %s\n", val.String())
				}
			}
		}
	}

	// Example 12: SimulateCalls with Value
	printSection("13. SimulateCalls - Simulate ETH Transfers")
	callsResult, err = public.SimulateCalls(ctx, publicClient, public.SimulateCallsParameters{
		Account: &vitalikAddress,
		Calls: []public.SimulateCall{
			{
				To:    &testAddress,
				Value: mustParseEther("0.1"),
			},
			{
				To:    &testRecipient,
				Value: mustParseEther("0.2"),
			},
		},
		StateOverrides: types.StateOverride{
			vitalikAddress: types.StateOverrideAccount{
				Balance: mustParseEther("1000"),
			},
		},
	})
	if err != nil {
		fmt.Printf("ETH transfer simulation error: %v\n", err)
	} else {
		fmt.Printf("ETH transfer simulation successful!\n")
		for i, result := range callsResult.Results {
			fmt.Printf("  Transfer %d: %s\n", i+1, result.Status)
		}
	}

	// Summary
	printHeader("Examples Complete")
	fmt.Println("Demonstrated Simulation features:")
	fmt.Println("  - CreateAccessList: Generate EIP-2930 access lists")
	fmt.Println("  - SimulateContract: Simulate contract write functions")
	fmt.Println("  - SimulateBlocks: Simulate multiple blocks with overrides")
	fmt.Println("  - SimulateCalls: Batch simulation with asset tracking")
	fmt.Println()
	fmt.Println("Key use cases:")
	fmt.Println("  - Pre-flight transaction validation")
	fmt.Println("  - Gas estimation with access lists")
	fmt.Println("  - Testing contract interactions")
	fmt.Println("  - Simulating state changes")
	fmt.Println()
}

// Helper functions

func mustParseEther(s string) *big.Int {
	v, err := unit.ParseEther(s)
	if err != nil {
		panic(fmt.Sprintf("invalid ether value %q: %v", s, err))
	}
	return v
}

func printHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 60))
}

func printSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}

func truncateAddress(addr common.Address) string {
	hex := addr.Hex()
	return hex[:10] + "..." + hex[len(hex)-4:]
}

func truncateURL(url string) string {
	if len(url) <= 40 {
		return url
	}
	return url[:30] + "..." + url[len(url)-7:]
}

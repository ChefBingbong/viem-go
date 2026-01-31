package erc20_test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
)

// This example demonstrates how to use the ERC20 contract bindings
// to interact with an ERC20 token contract.
//
// NOTE: This example requires a running Ethereum node to actually execute.
// It's provided as documentation of the API usage pattern.
func Example_readTokenInfo() {
	// Connect to an Ethereum node (e.g., Infura, Alchemy, or local node)
	rpcClient, err := client.NewClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer rpcClient.Close()

	// USDC token address on Ethereum mainnet
	tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	// Create ERC20 contract binding
	token, err := erc20.New(tokenAddress, rpcClient)
	if err != nil {
		fmt.Printf("Failed to create binding: %v\n", err)
		return
	}

	ctx := context.Background()

	// Read token metadata
	name, _ := token.Name(ctx)
	symbol, _ := token.Symbol(ctx)
	decimals, _ := token.Decimals(ctx)
	totalSupply, _ := token.TotalSupply(ctx)

	fmt.Printf("Token: %s (%s)\n", name, symbol)
	fmt.Printf("Decimals: %d\n", decimals)
	fmt.Printf("Total Supply: %s\n", totalSupply.String())

	// Check balance of an address
	userAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	balance, _ := token.BalanceOf(ctx, userAddress)
	fmt.Printf("Balance: %s\n", balance.String())

	// Check allowance
	spenderAddress := common.HexToAddress("0xabcdef1234567890123456789012345678901234")
	allowance, _ := token.Allowance(ctx, userAddress, spenderAddress)
	fmt.Printf("Allowance: %s\n", allowance.String())
}

// This example demonstrates how to send ERC20 tokens.
//
// NOTE: This requires a wallet/signer to sign the transaction.
// The example shows the API pattern but won't execute without proper setup.
func Example_transferTokens() {
	rpcClient, _ := client.NewClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
	defer rpcClient.Close()

	tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	token, _ := erc20.New(tokenAddress, rpcClient)

	ctx := context.Background()

	// Prepare transfer parameters
	from := common.HexToAddress("0xYOUR_ADDRESS")
	to := common.HexToAddress("0xRECIPIENT_ADDRESS")
	amount := big.NewInt(1000000) // 1 USDC (6 decimals)

	// Configure write options
	opts := contract.WriteOptions{
		From: from,
		// Gas and GasPrice will be estimated if not provided
	}

	// Send the transfer transaction
	txHash, err := token.Transfer(ctx, opts, to, amount)
	if err != nil {
		fmt.Printf("Transfer failed: %v\n", err)
		return
	}
	fmt.Printf("Transaction sent: %s\n", txHash.Hex())

	// Or use TransferAndWait to wait for confirmation
	receipt, err := token.TransferAndWait(ctx, opts, to, amount)
	if err != nil {
		fmt.Printf("Transfer failed: %v\n", err)
		return
	}
	fmt.Printf("Transaction confirmed in block %d\n", receipt.BlockNumber)
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)
}

// This example demonstrates how to approve a spender.
func Example_approveSpender() {
	rpcClient, _ := client.NewClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
	defer rpcClient.Close()

	tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	token, _ := erc20.New(tokenAddress, rpcClient)

	ctx := context.Background()

	owner := common.HexToAddress("0xYOUR_ADDRESS")
	spender := common.HexToAddress("0xSPENDER_CONTRACT") // e.g., Uniswap router

	// Approve unlimited spending (max uint256)
	maxAmount := new(big.Int).Sub(
		new(big.Int).Lsh(big.NewInt(1), 256),
		big.NewInt(1),
	)

	opts := contract.WriteOptions{From: owner}

	// Send approval and wait for confirmation
	receipt, err := token.ApproveAndWait(ctx, opts, spender, maxAmount)
	if err != nil {
		fmt.Printf("Approval failed: %v\n", err)
		return
	}

	if receipt.IsSuccess() {
		fmt.Println("Approval successful!")
	} else {
		fmt.Println("Approval transaction failed")
	}
}

// This example shows how to use the underlying contract for advanced usage.
func Example_advancedUsage() {
	rpcClient, _ := client.NewClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
	defer rpcClient.Close()

	tokenAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	token, _ := erc20.New(tokenAddress, rpcClient)

	ctx := context.Background()

	// Access the underlying contract for dynamic method calls
	underlying := token.Contract()

	// Call any method by name (useful for non-standard ERC20 methods)
	result, err := underlying.Read(ctx, "name")
	if err != nil {
		fmt.Printf("Read failed: %v\n", err)
		return
	}
	fmt.Printf("Name: %v\n", result[0])

	// Get raw calldata for a method (useful for multicall or batching)
	calldata, err := underlying.Calldata("balanceOf", common.HexToAddress("0x0"))
	if err != nil {
		fmt.Printf("Encoding failed: %v\n", err)
		return
	}
	fmt.Printf("Calldata: 0x%x\n", calldata)

	// Check what functions are available
	functions := underlying.FunctionNames()
	fmt.Printf("Available functions: %v\n", functions)
}

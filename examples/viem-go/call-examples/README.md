# Call Action Examples (viem-go)

Comprehensive examples demonstrating the `Call` action in viem-go, showcasing all supported features for simulating contract calls.

## Features Demonstrated

1. **Basic Contract Calls** - Reading ERC20 token name
2. **Calls with Parameters** - Encoding function arguments (balanceOf)
3. **Account (From Address)** - Specifying the sender
4. **Block Number Queries** - Reading historical state
5. **Block Tag Queries** - Using "latest", "pending", etc.
6. **Gas Parameters** - Setting gas limit and gas price
7. **EIP-1559 Fees** - Using maxFeePerGas and maxPriorityFeePerGas
8. **Value Transfers** - Simulating ETH transfers
9. **State Overrides** - Modifying account balances/nonces for simulation
10. **Block Overrides** - Modifying block context (gasLimit, baseFee, timestamp)
11. **Access Lists (EIP-2930)** - Pre-warming storage slots
12. **Deployless Calls** - Executing bytecode without deployment
13. **Combined Overrides** - Using state and block overrides together
14. **Error Handling** - Handling invalid parameter combinations

## Running the Example

```bash
cd examples/viem-go/call-examples
go run main.go
```

## Prerequisites

- Go 1.21+
- Internet connection (connects to Ethereum Mainnet via public RPC)

## Code Structure

```go
// Basic call
result, err := public.Call(ctx, client, public.CallParameters{
    To:   &contractAddress,
    Data: calldata,
})

// Call with state override
result, err := public.Call(ctx, client, public.CallParameters{
    To:   &contractAddress,
    Data: calldata,
    StateOverride: types.StateOverride{
        someAddress: types.StateOverrideAccount{
            Balance: big.NewInt(1000000),
        },
    },
})

// Deployless call
result, err := public.Call(ctx, client, public.CallParameters{
    Code: bytecode,
    Data: calldata,
})
```

## Notes

- CCIP-Read (EIP-3668) is supported but not demonstrated in this example
- Multicall batching is available when configured on the client
- State and block overrides require node support (most modern nodes support this)

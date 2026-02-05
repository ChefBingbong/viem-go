# Simulation Examples

This example demonstrates the simulation actions in viem-go:

- **CreateAccessList** - Generate EIP-2930 access lists for gas optimization
- **SimulateContract** - Simulate contract write functions before execution
- **SimulateBlocks** - Simulate multiple blocks with state/block overrides
- **SimulateCalls** - Batch simulation with optional asset change tracking

## Running the Example

```bash
cd examples/viem-go/simulate-examples
go run main.go
```

## Features Demonstrated

### CreateAccessList

```go
result, err := public.CreateAccessList(ctx, client, public.CreateAccessListParameters{
    Account: &account,
    To:      &contractAddr,
    Data:    calldata,
})
// result.AccessList - list of addresses and storage keys
// result.GasUsed - estimated gas with access list
```

### SimulateContract

```go
result, err := public.SimulateContract(ctx, client, public.SimulateContractParameters{
    Account:      &account,
    Address:      contractAddr,
    ABI:          contractABI,
    FunctionName: "transfer",
    Args:         []any{recipient, amount},
})
// result.Result - decoded return value
// result.Request - ready for writeContract
```

### SimulateBlocks

```go
results, err := public.SimulateBlocks(ctx, client, public.SimulateBlocksParameters{
    Blocks: []public.SimulateBlock{{
        BlockOverrides: &types.BlockOverrides{
            Number: &blockNum,
        },
        StateOverrides: types.StateOverride{
            addr: types.StateOverrideAccount{Balance: balance},
        },
        Calls: []public.SimulateBlockCall{{
            To:   &contractAddr,
            Data: calldata,
        }},
    }},
})
// results[0].Block - simulated block data
// results[0].Calls - call results with status, data, gas, logs
```

### SimulateCalls

```go
result, err := public.SimulateCalls(ctx, client, public.SimulateCallsParameters{
    Account: &account,
    Calls: []public.SimulateCall{
        {To: &addr1, Data: data1},
        {To: &addr2, Data: data2},
    },
    TraceAssetChanges: true, // Optional: track ETH/ERC20 changes
})
// result.Results - individual call results
// result.AssetChanges - ETH/token balance changes (if enabled)
// result.Block - block context
```

## Use Cases

1. **Pre-flight validation** - Check if a transaction will succeed before sending
2. **Gas estimation** - Get accurate gas estimates with access lists
3. **State simulation** - Test contract behavior with modified state
4. **Batch operations** - Simulate multiple calls in sequence
5. **Asset tracking** - Monitor balance changes from transactions

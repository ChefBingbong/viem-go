# Call Action Examples (viem TypeScript)

Comprehensive examples demonstrating the `call` action in viem, showcasing all supported features for simulating contract calls.

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
13. **Factory Pattern** - Deployless calls via factory contract
14. **Storage Overrides** - Modifying specific storage slots
15. **Error Handling** - Handling various error cases

## Running the Example

```bash
cd examples/viem/call-examples
bun install
bun run start
```

Or with npm:

```bash
npm install
npx tsx index.ts
```

## Prerequisites

- Node.js 18+ or Bun
- Internet connection (connects to Ethereum Mainnet via public RPC)

## Code Examples

### Basic Call

```typescript
const result = await client.call({
  to: contractAddress,
  data: calldata,
})
```

### Call with State Override

```typescript
const result = await client.call({
  to: contractAddress,
  data: calldata,
  stateOverride: [
    {
      address: someAddress,
      balance: parseEther('1000'),
    },
  ],
})
```

### Deployless Call

```typescript
const result = await client.call({
  code: bytecode,
  data: calldata,
})
```

### Call with Block Override

```typescript
const result = await client.call({
  to: contractAddress,
  data: calldata,
  blockOverrides: {
    gasLimit: 50000000n,
    baseFeePerGas: parseGwei('1'),
  },
})
```

## Notes

- CCIP-Read (EIP-3668) is supported but not demonstrated in this example
- Multicall batching is available via `client.multicall()`
- State and block overrides require node support (most modern nodes support this)

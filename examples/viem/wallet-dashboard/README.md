# Wallet Dashboard - viem TypeScript Example

A simple CLI-based wallet dashboard demonstrating key features of the [viem](https://viem.sh) TypeScript library for Ethereum.

## Features

- **Public Client Creation** - Connect to Ethereum mainnet via HTTP transport
- **Network Information** - Fetch current block number and chain ID
- **Gas Price Estimation** - Display current gas prices and fee history
- **Address Balance** - Check ETH balance and transaction count for any address
- **Message Signing** - Sign and verify messages using a local account
- **Block Information** - Display latest block details

## Prerequisites

- [Bun](https://bun.sh) runtime (recommended) or Node.js 18+

## Running the Example

```bash
# Install dependencies
bun install

# Run the dashboard
bun run start
```

## Output

The dashboard displays:

```
╔══════════════════════════════════════════════════╗
║           Wallet Dashboard (viem)                ║
║        Ethereum Network Information              ║
╚══════════════════════════════════════════════════╝

--- Creating Public Client ---
Connected to: Ethereum (Chain ID: 1)

--- Network Information ---
Current Block Number: 19,500,000
Chain ID: 1

--- Gas Prices ---
Current Gas Price: 25.5 Gwei
Avg Base Fee (last 4 blocks): 24.8 Gwei
Max Priority Fee: 1.5 Gwei

--- Address Balance ---
Address: 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
Balance: 1234.56 ETH
Transaction Count: 1,500

--- Message Signing Demo ---
Demo Account Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Message: "Hello from Wallet Dashboard!"
Signature: 0x...
Signature Valid: Yes

--- Latest Block Info ---
Block Hash: 0x...
Timestamp: 2024-03-15T12:00:00.000Z
Transactions: 150
Gas Used: 15,000,000
Gas Limit: 30,000,000
```

## Code Structure

- `main.ts` - Main application showcasing viem features

## Key viem APIs Used

| API | Description |
|-----|-------------|
| `createPublicClient` | Create a client for reading blockchain data |
| `getBlockNumber` | Get the current block number |
| `getChainId` | Get the chain ID |
| `getGasPrice` | Get the current gas price |
| `getFeeHistory` | Get historical fee data |
| `getBalance` | Get ETH balance of an address |
| `getTransactionCount` | Get nonce/transaction count |
| `getBlock` | Get block information |
| `privateKeyToAccount` | Create account from private key |
| `signMessage` | Sign a message with account |
| `verifyMessage` | Verify a signed message |

## Related

- [viem-go Wallet Dashboard](../../viem-go/wallet-dashboard/) - Go equivalent of this example

# Wallet Dashboard - viem-go Example

A simple CLI-based wallet dashboard demonstrating key features of the [viem-go](https://github.com/ChefBingbong/viem-go) library for Ethereum.

## Features

- **Public Client Creation** - Connect to Ethereum mainnet via HTTP transport
- **Network Information** - Fetch current block number and chain ID
- **Gas Price Estimation** - Display current gas prices and max priority fee
- **Address Balance** - Check ETH balance and transaction count for any address
- **Message Signing** - Sign and verify messages using a local account
- **Block Information** - Display latest block details

## Prerequisites

- Go 1.21+

## Running the Example

```bash
# From the wallet-dashboard directory
go run main.go

# Or from the repository root
go run ./examples/viem-go/wallet-dashboard
```

## Output

The dashboard displays:

```
╔══════════════════════════════════════════════════╗
║          Wallet Dashboard (viem-go)              ║
║        Ethereum Network Information              ║
╚══════════════════════════════════════════════════╝

--- Creating Public Client ---
Connected to: Ethereum (Chain ID: 1)

--- Network Information ---
Current Block Number: 19,500,000
Chain ID: 1

--- Gas Prices ---
Current Gas Price: 25.5 Gwei
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
Timestamp: 2024-03-15T12:00:00Z
Transactions: 150
Gas Used: 15,000,000
Gas Limit: 30,000,000
```

## Code Structure

- `main.go` - Main application showcasing viem-go features

## Key viem-go APIs Used

| API | Package | Description |
|-----|---------|-------------|
| `CreatePublicClient` | `client` | Create a client for reading blockchain data |
| `GetBlockNumber` | `client` | Get the current block number |
| `GetChainID` | `client` | Get the chain ID |
| `GetGasPrice` | `client` | Get the current gas price |
| `GetMaxPriorityFeePerGas` | `client` | Get the max priority fee |
| `GetBalance` | `client` | Get ETH balance of an address |
| `GetTransactionCount` | `client` | Get nonce/transaction count |
| `GetBlock` | `client` | Get block information |
| `PrivateKeyToAccount` | `accounts` | Create account from private key |
| `SignMessage` | `accounts` | Sign a message with account |
| `VerifyMessage` | `utils/signature` | Verify a signed message |
| `FormatEther` | `utils/unit` | Format wei as ETH string |
| `FormatGwei` | `utils/unit` | Format wei as Gwei string |

## Related

- [viem Wallet Dashboard](../../viem/wallet-dashboard/) - TypeScript equivalent of this example

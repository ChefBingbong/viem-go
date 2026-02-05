# Verify Examples - viem-go

Comprehensive examples demonstrating signature verification actions in viem-go.

## Features Demonstrated

### VerifyMessage
Verify that a message was signed by the provided address.

```go
valid, err := public.VerifyMessage(ctx, client, public.VerifyMessageParameters{
    Address:   common.HexToAddress("0x..."),
    Message:   signature.NewSignableMessage("Hello, world!"),
    Signature: "0x...",
})
```

### VerifyHash
Verify a raw 32-byte hash signature.

```go
valid, err := public.VerifyHash(ctx, client, public.VerifyHashParameters{
    Address:   common.HexToAddress("0x..."),
    Hash:      "0x...", // 32-byte hash
    Signature: "0x...",
})
```

### VerifyTypedData
Verify EIP-712 typed data signatures.

```go
valid, err := public.VerifyTypedData(ctx, client, public.VerifyTypedDataParameters{
    Address:   common.HexToAddress("0x..."),
    TypedData: signature.TypedDataDefinition{
        Domain: signature.TypedDataDomain{
            Name:    "Example DApp",
            Version: "1",
            ChainId: big.NewInt(1),
        },
        Types: map[string][]signature.TypedDataField{
            "Mail": {
                {Name: "contents", Type: "string"},
            },
        },
        PrimaryType: "Mail",
        Message: map[string]any{
            "contents": "Hello!",
        },
    },
    Signature: "0x...",
})
```

## Verification Standards Supported

| Standard | Description | Use Case |
|----------|-------------|----------|
| **ECDSA Recovery** | Traditional signature verification | EOA (Externally Owned Account) |
| **ERC-1271** | Smart contract signature verification | Deployed smart contract wallets |
| **ERC-6492** | Counterfactual signature verification | Undeployed smart accounts |

## Local vs On-Chain Verification

### On-Chain Verification (Default)
Uses ERC-6492 deployless verification to support all account types:

```go
valid, err := public.VerifyMessage(ctx, client, params)
```

### Local Verification (EOA-only)
For simple EOA verification without network calls:

```go
valid, err := public.VerifyMessageLocal(params)
valid, err := public.VerifyTypedDataLocal(params)
```

## Smart Account Verification

### Counterfactual Verification
For smart accounts that haven't been deployed yet:

```go
valid, err := public.VerifyHash(ctx, client, public.VerifyHashParameters{
    Address:     smartAccountAddress,
    Hash:        messageHash,
    Signature:   signature,
    Factory:     &factoryAddress,     // ERC-4337 factory
    FactoryData: factoryCalldata,     // Deploy calldata
})
```

### Using a Deployed Verifier
For gas efficiency with a pre-deployed ERC-6492 verifier:

```go
valid, err := public.VerifyHash(ctx, client, public.VerifyHashParameters{
    Address:                address,
    Hash:                   hash,
    Signature:              signature,
    ERC6492VerifierAddress: &verifierAddress,
})
```

## Signature Formats

All verify functions accept multiple signature formats:

```go
// Hex string
Signature: "0x..."

// Raw bytes
Signature: []byte{...}

// Signature struct
Signature: &signature.Signature{
    R:       "0x...",
    S:       "0x...",
    V:       big.NewInt(27),
    YParity: 0,
}
```

## Running the Example

```bash
cd examples/viem-go/verify-examples
go run main.go
```

## Example Output

```
============================================================
  Signature Verification Examples (viem-go)
============================================================
Test Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

--- 1. Creating Public Client ---
Connected to Ethereum Mainnet

--- 2. Sign and Verify Simple Message ---
Message: "Hello, viem-go!"
Message Hash: 0xd9eba16ed0ecae432b71fe008c98cc872bb4cc214d3220a36f365326cf807d68
Signature: 0x1234567890abcdef...12345678
Signature valid: true
...
```

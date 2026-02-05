// Verify Examples - viem-go
//
// Comprehensive examples demonstrating signature verification actions:
// - VerifyHash: Verify a message hash (supports EOA, ERC-1271, ERC-6492)
// - VerifyMessage: Verify a signed message
// - VerifyTypedData: Verify EIP-712 typed data signatures
//
// These actions support:
// - EOA (Externally Owned Account) signatures via ECDSA recovery
// - Smart Contract Account signatures via ERC-1271
// - Counterfactual (undeployed) smart account signatures via ERC-6492
//
// For local-only EOA verification (no RPC calls), use:
// - VerifyMessageLocal / VerifyTypedDataLocal
package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// Test private key for examples (DO NOT use in production!)
// This is a well-known test private key from Hardhat/Foundry
const testPrivateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func main() {
	ctx := context.Background()

	printHeader("Signature Verification Examples (viem-go)")

	// Derive test address from private key
	privateKey, err := crypto.HexToECDSA(testPrivateKeyHex)
	if err != nil {
		fmt.Printf("Error parsing private key: %v\n", err)
		return
	}
	testAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("Test Address: %s\n", testAddress.Hex())

	// Create Public Client
	printSection("1. Creating Public Client")
	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Mainnet,
		Transport: transport.HTTP("https://eth.llamarpc.com"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()
	fmt.Println("Connected to Ethereum Mainnet")

	// ============================================================================
	// VerifyMessage Examples
	// ============================================================================
	printHeader("VerifyMessage Examples")

	// Example 1: Sign and Verify a Simple Message
	printSection("2. Sign and Verify Simple Message")
	message := "Hello, viem-go!"
	messageHash := signature.HashMessageBytes(signature.NewSignableMessage(message))
	sig, err := crypto.Sign(messageHash, privateKey)
	if err != nil {
		fmt.Printf("Error signing message: %v\n", err)
		return
	}
	// Adjust v value for Ethereum (add 27)
	if sig[64] < 27 {
		sig[64] += 27
	}
	sigHex := fmt.Sprintf("0x%x", sig)

	fmt.Printf("Message: %q\n", message)
	fmt.Printf("Message Hash: 0x%x\n", messageHash)
	fmt.Printf("Signature: %s...%s\n", sigHex[:20], sigHex[len(sigHex)-8:])

	// Verify the message (on-chain verification via ERC-6492)
	valid, err := public.VerifyMessage(ctx, publicClient, public.VerifyMessageParameters{
		Address:   testAddress,
		Message:   signature.NewSignableMessage(message),
		Signature: sigHex,
	})
	if err != nil {
		fmt.Printf("Verification error: %v\n", err)
	} else {
		fmt.Printf("Signature valid: %v\n", valid)
	}

	// Example 2: Local Message Verification (no RPC call)
	printSection("3. Local Message Verification (EOA-only, no network)")
	validLocal, err := public.VerifyMessageLocal(public.VerifyMessageParameters{
		Address:   testAddress,
		Message:   signature.NewSignableMessage(message),
		Signature: sigHex,
	})
	if err != nil {
		fmt.Printf("Local verification error: %v\n", err)
	} else {
		fmt.Printf("Local verification result: %v\n", validLocal)
	}

	// Example 3: Verify with Raw Message Bytes
	printSection("4. Verify Message with Raw Bytes")
	rawMessage := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" in bytes
	rawMessageHash := signature.HashMessageBytes(signature.NewSignableMessageRaw(rawMessage))
	rawSig, _ := crypto.Sign(rawMessageHash, privateKey)
	if rawSig[64] < 27 {
		rawSig[64] += 27
	}

	validRaw, err := public.VerifyMessage(ctx, publicClient, public.VerifyMessageParameters{
		Address:   testAddress,
		Message:   signature.NewSignableMessageRaw(rawMessage),
		Signature: fmt.Sprintf("0x%x", rawSig),
	})
	if err != nil {
		fmt.Printf("Raw message verification error: %v\n", err)
	} else {
		fmt.Printf("Raw message signature valid: %v\n", validRaw)
	}

	// Example 4: Verify with Wrong Address (should fail)
	printSection("5. Verify with Wrong Address (Expected: false)")
	wrongAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	validWrong, err := public.VerifyMessage(ctx, publicClient, public.VerifyMessageParameters{
		Address:   wrongAddress,
		Message:   signature.NewSignableMessage(message),
		Signature: sigHex,
	})
	if err != nil {
		fmt.Printf("Expected failure: %v\n", err)
	} else {
		fmt.Printf("Verification result (expected false): %v\n", validWrong)
	}

	// ============================================================================
	// VerifyHash Examples
	// ============================================================================
	printHeader("VerifyHash Examples")

	// Example 5: Verify Raw Hash
	printSection("6. Verify Raw Hash")
	// Create a raw 32-byte hash (not Ethereum prefixed)
	rawHash := crypto.Keccak256([]byte("some data to hash"))
	hashHex := fmt.Sprintf("0x%x", rawHash)

	// Sign the raw hash directly
	hashSig, _ := crypto.Sign(rawHash, privateKey)
	if hashSig[64] < 27 {
		hashSig[64] += 27
	}

	fmt.Printf("Raw Hash: %s\n", hashHex)
	validHash, err := public.VerifyHash(ctx, publicClient, public.VerifyHashParameters{
		Address:   testAddress,
		Hash:      hashHex,
		Signature: fmt.Sprintf("0x%x", hashSig),
	})
	if err != nil {
		fmt.Printf("Hash verification error: %v\n", err)
	} else {
		fmt.Printf("Hash signature valid: %v\n", validHash)
	}

	// Example 6: Verify with Signature Struct
	printSection("7. Verify with Signature Struct")
	// Parse signature components
	r := fmt.Sprintf("0x%x", hashSig[:32])
	s := fmt.Sprintf("0x%x", hashSig[32:64])
	v := big.NewInt(int64(hashSig[64]))

	sigStruct := &signature.Signature{
		R:       r,
		S:       s,
		V:       v,
		YParity: int(hashSig[64] - 27),
	}

	fmt.Printf("Signature components:\n")
	fmt.Printf("  R: %s...%s\n", r[:18], r[len(r)-8:])
	fmt.Printf("  S: %s...%s\n", s[:18], s[len(s)-8:])
	fmt.Printf("  V: %d\n", v.Uint64())

	validStruct, err := public.VerifyHash(ctx, publicClient, public.VerifyHashParameters{
		Address:   testAddress,
		Hash:      hashHex,
		Signature: sigStruct,
	})
	if err != nil {
		fmt.Printf("Struct signature verification error: %v\n", err)
	} else {
		fmt.Printf("Struct signature valid: %v\n", validStruct)
	}

	// Example 7: Verify at Historical Block
	printSection("8. Verify at Historical Block")
	blockNum := uint64(19000000) // Historical block
	validAtBlock, err := public.VerifyHash(ctx, publicClient, public.VerifyHashParameters{
		Address:     testAddress,
		Hash:        hashHex,
		Signature:   fmt.Sprintf("0x%x", hashSig),
		BlockNumber: &blockNum,
	})
	if err != nil {
		fmt.Printf("Historical verification error: %v\n", err)
	} else {
		fmt.Printf("Signature valid at block %d: %v\n", blockNum, validAtBlock)
	}

	// ============================================================================
	// VerifyTypedData Examples
	// ============================================================================
	printHeader("VerifyTypedData Examples (EIP-712)")

	// Example 8: Verify EIP-712 Typed Data
	printSection("9. Verify EIP-712 Typed Data")

	// Define typed data (common EIP-712 structure)
	typedData := signature.TypedDataDefinition{
		Domain: signature.TypedDataDomain{
			Name:              "Example DApp",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		},
		Types: map[string][]signature.TypedDataField{
			"Person": {
				{Name: "name", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
			"Mail": {
				{Name: "from", Type: "Person"},
				{Name: "to", Type: "Person"},
				{Name: "contents", Type: "string"},
			},
		},
		PrimaryType: "Mail",
		Message: map[string]any{
			"from": map[string]any{
				"name":   "Alice",
				"wallet": testAddress.Hex(),
			},
			"to": map[string]any{
				"name":   "Bob",
				"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
			},
			"contents": "Hello, Bob!",
		},
	}

	fmt.Printf("Typed Data Domain:\n")
	fmt.Printf("  Name: %s\n", typedData.Domain.Name)
	fmt.Printf("  Version: %s\n", typedData.Domain.Version)
	fmt.Printf("  ChainId: %s\n", typedData.Domain.ChainId.String())
	fmt.Printf("Primary Type: %s\n", typedData.PrimaryType)

	// Hash the typed data
	typedDataHash, err := signature.HashTypedData(typedData)
	if err != nil {
		fmt.Printf("Error hashing typed data: %v\n", err)
		return
	}
	fmt.Printf("Typed Data Hash: %s\n", typedDataHash)

	// Sign the typed data hash
	typedDataHashBytes := common.FromHex(typedDataHash)
	typedDataSig, _ := crypto.Sign(typedDataHashBytes, privateKey)
	if typedDataSig[64] < 27 {
		typedDataSig[64] += 27
	}

	// Verify typed data
	validTyped, err := public.VerifyTypedData(ctx, publicClient, public.VerifyTypedDataParameters{
		Address:   testAddress,
		TypedData: typedData,
		Signature: fmt.Sprintf("0x%x", typedDataSig),
	})
	if err != nil {
		fmt.Printf("Typed data verification error: %v\n", err)
	} else {
		fmt.Printf("Typed data signature valid: %v\n", validTyped)
	}

	// Example 9: Local Typed Data Verification
	printSection("10. Local Typed Data Verification (EOA-only)")
	validTypedLocal, err := public.VerifyTypedDataLocal(public.VerifyTypedDataParameters{
		Address:   testAddress,
		TypedData: typedData,
		Signature: fmt.Sprintf("0x%x", typedDataSig),
	})
	if err != nil {
		fmt.Printf("Local typed data verification error: %v\n", err)
	} else {
		fmt.Printf("Local typed data verification result: %v\n", validTypedLocal)
	}

	// Example 10: Permit-style Typed Data
	printSection("11. ERC20 Permit-style Typed Data")
	permitData := signature.TypedDataDefinition{
		Domain: signature.TypedDataDomain{
			Name:              "USD Coin",
			Version:           "2",
			ChainId:           big.NewInt(1),
			VerifyingContract: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
		},
		Types: map[string][]signature.TypedDataField{
			"Permit": {
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "Permit",
		Message: map[string]any{
			"owner":    testAddress.Hex(),
			"spender":  "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
			"value":    big.NewInt(1000000), // 1 USDC
			"nonce":    big.NewInt(0),
			"deadline": big.NewInt(1893456000), // Far future
		},
	}

	permitHash, _ := signature.HashTypedData(permitData)
	permitHashBytes := common.FromHex(permitHash)
	permitSig, _ := crypto.Sign(permitHashBytes, privateKey)
	if permitSig[64] < 27 {
		permitSig[64] += 27
	}

	fmt.Printf("Permit Domain: %s (USDC)\n", permitData.Domain.Name)
	fmt.Printf("Permit Owner: %s\n", testAddress.Hex())
	fmt.Printf("Permit Hash: %s\n", permitHash)

	validPermit, err := public.VerifyTypedData(ctx, publicClient, public.VerifyTypedDataParameters{
		Address:   testAddress,
		TypedData: permitData,
		Signature: fmt.Sprintf("0x%x", permitSig),
	})
	if err != nil {
		fmt.Printf("Permit verification error: %v\n", err)
	} else {
		fmt.Printf("Permit signature valid: %v\n", validPermit)
	}

	// ============================================================================
	// Smart Account Examples (ERC-6492)
	// ============================================================================
	printHeader("Smart Account Verification (ERC-6492)")

	// Example 11: Counterfactual Verification (for undeployed smart accounts)
	printSection("12. Counterfactual Verification Setup")
	fmt.Println("ERC-6492 enables signature verification for smart accounts")
	fmt.Println("that haven't been deployed yet (counterfactual addresses).")
	fmt.Println()
	fmt.Println("Usage with factory/factoryData:")
	fmt.Println("  public.VerifyHash(ctx, client, public.VerifyHashParameters{")
	fmt.Println("      Address:     smartAccountAddress,")
	fmt.Println("      Hash:        messageHash,")
	fmt.Println("      Signature:   signature,")
	fmt.Println("      Factory:     &factoryAddress,     // ERC-4337 factory")
	fmt.Println("      FactoryData: factoryCalldata,     // Deploy calldata")
	fmt.Println("  })")
	fmt.Println()
	fmt.Println("The ERC-6492 validator will:")
	fmt.Println("  1. Deploy the smart account using factory/factoryData")
	fmt.Println("  2. Call isValidSignature on the deployed account")
	fmt.Println("  3. Return the verification result")

	// Example 12: Using a deployed ERC-6492 verifier
	printSection("13. Using Deployed ERC-6492 Verifier")
	fmt.Println("For gas efficiency, you can use a pre-deployed verifier:")
	fmt.Println()
	fmt.Println("  verifierAddr := common.HexToAddress(\"0x...\")")
	fmt.Println("  public.VerifyHash(ctx, client, public.VerifyHashParameters{")
	fmt.Println("      Address:                address,")
	fmt.Println("      Hash:                   hash,")
	fmt.Println("      Signature:              signature,")
	fmt.Println("      ERC6492VerifierAddress: &verifierAddr,")
	fmt.Println("  })")

	// Summary
	printHeader("Examples Complete")
	fmt.Println("Demonstrated Verification features:")
	fmt.Println("  - VerifyMessage: Verify Ethereum signed messages")
	fmt.Println("  - VerifyHash: Verify raw 32-byte hashes")
	fmt.Println("  - VerifyTypedData: Verify EIP-712 typed data")
	fmt.Println("  - Local verification (no RPC calls)")
	fmt.Println("  - Signature struct support (r, s, v components)")
	fmt.Println("  - Historical block verification")
	fmt.Println("  - ERC-6492 smart account verification")
	fmt.Println("  - ERC20 Permit-style typed data")
	fmt.Println()
	fmt.Println("Verification Standards Supported:")
	fmt.Println("  - ECDSA recovery (EOA accounts)")
	fmt.Println("  - ERC-1271 (deployed smart contracts)")
	fmt.Println("  - ERC-6492 (counterfactual smart accounts)")
	fmt.Println()
}

// Helper functions

func printHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 60))
}

func printSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}

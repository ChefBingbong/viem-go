package utils

import (
	"github.com/ChefBingbong/viem-go/utils/hash"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// SignTransactionParameters contains parameters for signing a transaction.
type SignTransactionParameters struct {
	// PrivateKey is the private key to sign with (hex string with 0x prefix).
	PrivateKey string
	// Transaction is the transaction to sign.
	Transaction *transaction.Transaction
	// Serializer is an optional custom serializer function.
	// If not provided, the default SerializeTransaction is used.
	Serializer func(tx *transaction.Transaction, sig *transaction.Signature) (string, error)
}

// SignTransaction signs a transaction with a private key.
//
// The function:
// 1. Serializes the transaction to RLP
// 2. Hashes the serialized transaction with keccak256
// 3. Signs the hash with the private key
// 4. Re-serializes the transaction with the signature attached
//
// Example:
//
//	signedTx, err := SignTransaction(SignTransactionParameters{
//		PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
//		Transaction: &transaction.Transaction{
//			Type:                 transaction.TransactionTypeEIP1559,
//			ChainId:              1,
//			Nonce:                0,
//			MaxPriorityFeePerGas: big.NewInt(1000000000),
//			MaxFeePerGas:         big.NewInt(2000000000),
//			Gas:                  big.NewInt(21000),
//			To:                   "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
//			Value:                big.NewInt(1000000000000000000),
//		},
//	})
func SignTransaction(params SignTransactionParameters) (string, error) {
	// Use default serializer if not provided
	serializer := params.Serializer
	if serializer == nil {
		serializer = transaction.SerializeTransaction
	}

	// Prepare transaction for signing
	// For EIP-4844 transactions, we sign without sidecars (network wrapper)
	signableTx := params.Transaction
	if params.Transaction.Type == transaction.TransactionTypeEIP4844 {
		// Create a copy without sidecars for signing
		txCopy := *params.Transaction
		txCopy.Sidecars = nil
		signableTx = &txCopy
	}

	// Serialize the transaction (without signature)
	serialized, err := serializer(signableTx, nil)
	if err != nil {
		return "", err
	}

	// Hash the serialized transaction
	txHash := hash.Keccak256(serialized)

	// Sign the hash
	sig, err := SignToSignature(txHash, params.PrivateKey)
	if err != nil {
		return "", err
	}

	// Convert to transaction.Signature
	txSig := &transaction.Signature{
		R:       sig.R,
		S:       sig.S,
		V:       sig.V,
		YParity: sig.YParity,
	}

	// Serialize the transaction with signature
	return serializer(params.Transaction, txSig)
}

// MustSignTransaction signs a transaction or panics on error.
func MustSignTransaction(params SignTransactionParameters) string {
	result, err := SignTransaction(params)
	if err != nil {
		panic(err)
	}
	return result
}

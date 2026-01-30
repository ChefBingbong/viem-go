package transaction

import (
	"fmt"
)

// GetTransactionType determines the transaction type based on the transaction fields.
//
// Example:
//
//	txType, err := GetTransactionType(&Transaction{
//		MaxFeePerGas: big.NewInt(1000000000),
//	})
//	// txType = "eip1559"
func GetTransactionType(tx *Transaction) (TransactionType, error) {
	// If type is explicitly set, use it
	if tx.Type != "" {
		return tx.Type, nil
	}

	// EIP-7702: has authorizationList
	if len(tx.AuthorizationList) > 0 {
		return TransactionTypeEIP7702, nil
	}

	// EIP-4844: has blob-related fields
	if len(tx.Blobs) > 0 || len(tx.BlobVersionedHashes) > 0 ||
		tx.MaxFeePerBlobGas != nil || len(tx.Sidecars) > 0 {
		return TransactionTypeEIP4844, nil
	}

	// EIP-1559: has maxFeePerGas or maxPriorityFeePerGas
	if tx.MaxFeePerGas != nil || tx.MaxPriorityFeePerGas != nil {
		return TransactionTypeEIP1559, nil
	}

	// Legacy or EIP-2930: has gasPrice
	if tx.GasPrice != nil {
		// EIP-2930: has accessList
		if len(tx.AccessList) > 0 {
			return TransactionTypeEIP2930, nil
		}
		return TransactionTypeLegacy, nil
	}

	return "", fmt.Errorf("%w: cannot determine transaction type from fields", ErrInvalidSerializableTransaction)
}

// MustGetTransactionType returns the transaction type or panics.
func MustGetTransactionType(tx *Transaction) TransactionType {
	txType, err := GetTransactionType(tx)
	if err != nil {
		panic(err)
	}
	return txType
}

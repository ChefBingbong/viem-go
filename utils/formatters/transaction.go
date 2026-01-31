package formatters

import (
	"math/big"
)

// FormatTransaction formats an RPC transaction into a Transaction struct.
//
// Example:
//
//	rpcTx := RpcTransaction{
//		Hash:        "0x...",
//		BlockNumber: "0x1234",
//		Value:       "0xde0b6b3a7640000",
//	}
//	tx := FormatTransaction(rpcTx)
func FormatTransaction(tx RpcTransaction) Transaction {
	result := Transaction{
		From:    tx.From,
		Hash:    tx.Hash,
		Input:   tx.Input,
		R:       tx.R,
		S:       tx.S,
	}

	// Block hash
	if tx.BlockHash != "" {
		result.BlockHash = &tx.BlockHash
	}

	// Block number
	if tx.BlockNumber != "" {
		result.BlockNumber = hexToBigInt(tx.BlockNumber)
	}

	// Chain ID
	if tx.ChainID != "" {
		chainID := hexToInt(tx.ChainID)
		result.ChainID = &chainID
	}

	// Gas
	if tx.Gas != "" {
		result.Gas = hexToBigInt(tx.Gas)
	}

	// Gas price
	if tx.GasPrice != "" {
		result.GasPrice = hexToBigInt(tx.GasPrice)
	}

	// Max fee per blob gas
	if tx.MaxFeePerBlobGas != "" {
		result.MaxFeePerBlobGas = hexToBigInt(tx.MaxFeePerBlobGas)
	}

	// Max fee per gas
	if tx.MaxFeePerGas != "" {
		result.MaxFeePerGas = hexToBigInt(tx.MaxFeePerGas)
	}

	// Max priority fee per gas
	if tx.MaxPriorityFeePerGas != "" {
		result.MaxPriorityFeePerGas = hexToBigInt(tx.MaxPriorityFeePerGas)
	}

	// Nonce
	if tx.Nonce != "" {
		nonce := hexToInt(tx.Nonce)
		result.Nonce = &nonce
	}

	// To
	if tx.To != "" {
		result.To = &tx.To
	}

	// Transaction index
	if tx.TransactionIndex != "" {
		idx := hexToInt(tx.TransactionIndex)
		result.TransactionIndex = &idx
	}

	// Type
	if tx.Type != "" {
		if txType, ok := TransactionTypeFromHex[tx.Type]; ok {
			result.Type = txType
		}
		result.TypeHex = tx.Type
	}

	// Value
	if tx.Value != "" {
		result.Value = hexToBigInt(tx.Value)
	}

	// V
	if tx.V != "" {
		result.V = hexToBigInt(tx.V)
	}

	// Access list
	if len(tx.AccessList) > 0 {
		result.AccessList = tx.AccessList
	}

	// Authorization list
	if len(tx.AuthorizationList) > 0 {
		result.AuthorizationList = formatAuthorizationList(tx.AuthorizationList)
	}

	// Derive yParity
	result.YParity = deriveYParity(tx.YParity, result.V)

	// Clean up based on transaction type
	cleanupTransactionByType(&result)

	return result
}

// FormatTransactions formats multiple RPC transactions.
func FormatTransactions(txs []RpcTransaction) []Transaction {
	result := make([]Transaction, len(txs))
	for i, tx := range txs {
		result[i] = FormatTransaction(tx)
	}
	return result
}

// deriveYParity derives yParity from the provided value or v.
func deriveYParity(yParityHex string, v *big.Int) *int {
	// If yParity is provided, use it
	if yParityHex != "" {
		yParity := hexToInt(yParityHex)
		return &yParity
	}

	// Try to derive from v
	if v != nil {
		var yParity int
		if v.Cmp(big.NewInt(0)) == 0 || v.Cmp(big.NewInt(27)) == 0 {
			yParity = 0
			return &yParity
		}
		if v.Cmp(big.NewInt(1)) == 0 || v.Cmp(big.NewInt(28)) == 0 {
			yParity = 1
			return &yParity
		}
		if v.Cmp(big.NewInt(35)) >= 0 {
			// EIP-155: v = chainId * 2 + 35 + yParity
			mod := new(big.Int).Mod(v, big.NewInt(2))
			if mod.Cmp(big.NewInt(0)) == 0 {
				yParity = 1
			} else {
				yParity = 0
			}
			return &yParity
		}
	}

	return nil
}

// cleanupTransactionByType removes fields not applicable to the transaction type.
func cleanupTransactionByType(tx *Transaction) {
	switch tx.Type {
	case TransactionTypeLegacy:
		tx.AccessList = nil
		tx.MaxFeePerBlobGas = nil
		tx.MaxFeePerGas = nil
		tx.MaxPriorityFeePerGas = nil
		tx.YParity = nil
	case TransactionTypeEIP2930:
		tx.MaxFeePerBlobGas = nil
		tx.MaxFeePerGas = nil
		tx.MaxPriorityFeePerGas = nil
	case TransactionTypeEIP1559:
		tx.MaxFeePerBlobGas = nil
	}
}

// formatAuthorizationList formats the authorization list from RPC format.
func formatAuthorizationList(authList []any) []SignedAuthorization {
	result := make([]SignedAuthorization, 0, len(authList))
	for _, auth := range authList {
		if authMap, ok := auth.(map[string]any); ok {
			sa := SignedAuthorization{}
			
			if addr, ok := authMap["address"].(string); ok {
				sa.Address = addr
			}
			if chainID, ok := authMap["chainId"].(string); ok {
				sa.ChainID = hexToInt(chainID)
			}
			if nonce, ok := authMap["nonce"].(string); ok {
				sa.Nonce = hexToInt(nonce)
			}
			if r, ok := authMap["r"].(string); ok {
				sa.R = r
			}
			if s, ok := authMap["s"].(string); ok {
				sa.S = s
			}
			if yParity, ok := authMap["yParity"].(string); ok {
				sa.YParity = hexToInt(yParity)
			}
			
			result = append(result, sa)
		}
	}
	return result
}

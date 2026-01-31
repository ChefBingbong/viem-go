package formatters

// FormatTransactionReceipt formats an RPC transaction receipt into a TransactionReceipt struct.
//
// Example:
//
//	rpcReceipt := RpcTransactionReceipt{
//		BlockNumber:       "0x1234",
//		GasUsed:           "0x5208",
//		Status:            "0x1",
//	}
//	receipt := FormatTransactionReceipt(rpcReceipt)
func FormatTransactionReceipt(receipt RpcTransactionReceipt) TransactionReceipt {
	result := TransactionReceipt{
		BlockHash:       receipt.BlockHash,
		From:            receipt.From,
		LogsBloom:       receipt.LogsBloom,
		Root:            receipt.Root,
		TransactionHash: receipt.TransactionHash,
	}

	// Block number
	if receipt.BlockNumber != "" {
		result.BlockNumber = hexToBigInt(receipt.BlockNumber)
	}

	// Contract address
	if receipt.ContractAddress != "" {
		result.ContractAddress = &receipt.ContractAddress
	}

	// Cumulative gas used
	if receipt.CumulativeGasUsed != "" {
		result.CumulativeGasUsed = hexToBigInt(receipt.CumulativeGasUsed)
	}

	// Effective gas price
	if receipt.EffectiveGasPrice != "" {
		result.EffectiveGasPrice = hexToBigInt(receipt.EffectiveGasPrice)
	}

	// Gas used
	if receipt.GasUsed != "" {
		result.GasUsed = hexToBigInt(receipt.GasUsed)
	}

	// Logs
	if len(receipt.Logs) > 0 {
		result.Logs = FormatLogs(receipt.Logs)
	}

	// To
	if receipt.To != "" {
		result.To = &receipt.To
	}

	// Transaction index
	if receipt.TransactionIndex != "" {
		idx := hexToInt(receipt.TransactionIndex)
		result.TransactionIndex = &idx
	}

	// Status
	if receipt.Status != "" {
		if status, ok := ReceiptStatusFromHex[receipt.Status]; ok {
			result.Status = status
		}
	}

	// Type
	if receipt.Type != "" {
		if txType, ok := TransactionTypeFromHex[receipt.Type]; ok {
			result.Type = txType
		} else {
			// Use the raw type if not recognized
			result.Type = TransactionType(receipt.Type)
		}
	}

	// Blob gas price (EIP-4844)
	if receipt.BlobGasPrice != "" {
		result.BlobGasPrice = hexToBigInt(receipt.BlobGasPrice)
	}

	// Blob gas used (EIP-4844)
	if receipt.BlobGasUsed != "" {
		result.BlobGasUsed = hexToBigInt(receipt.BlobGasUsed)
	}

	return result
}

// FormatTransactionReceipts formats multiple RPC transaction receipts.
func FormatTransactionReceipts(receipts []RpcTransactionReceipt) []TransactionReceipt {
	result := make([]TransactionReceipt, len(receipts))
	for i, receipt := range receipts {
		result[i] = FormatTransactionReceipt(receipt)
	}
	return result
}

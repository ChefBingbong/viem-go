package formatters

// FormatBlock formats an RPC block into a Block struct.
//
// Example:
//
//	rpcBlock := RpcBlock{
//		Number:    "0x1234",
//		Timestamp: "0x5f5e100",
//		GasUsed:   "0x5208",
//	}
//	block := FormatBlock(rpcBlock)
func FormatBlock(block RpcBlock) Block {
	result := Block{
		ExtraData:        block.ExtraData,
		Miner:            block.Miner,
		MixHash:          block.MixHash,
		ParentHash:       block.ParentHash,
		ReceiptsRoot:     block.ReceiptsRoot,
		Sha3Uncles:       block.Sha3Uncles,
		StateRoot:        block.StateRoot,
		TransactionsRoot: block.TransactionsRoot,
		Uncles:           block.Uncles,
	}

	// Base fee per gas
	if block.BaseFeePerGas != "" {
		result.BaseFeePerGas = hexToBigInt(block.BaseFeePerGas)
	}

	// Blob gas used
	if block.BlobGasUsed != "" {
		result.BlobGasUsed = hexToBigInt(block.BlobGasUsed)
	}

	// Difficulty
	if block.Difficulty != "" {
		result.Difficulty = hexToBigInt(block.Difficulty)
	}

	// Excess blob gas
	if block.ExcessBlobGas != "" {
		result.ExcessBlobGas = hexToBigInt(block.ExcessBlobGas)
	}

	// Gas limit
	if block.GasLimit != "" {
		result.GasLimit = hexToBigInt(block.GasLimit)
	}

	// Gas used
	if block.GasUsed != "" {
		result.GasUsed = hexToBigInt(block.GasUsed)
	}

	// Hash
	if block.Hash != "" {
		result.Hash = &block.Hash
	}

	// Logs bloom
	if block.LogsBloom != "" {
		result.LogsBloom = &block.LogsBloom
	}

	// Nonce
	if block.Nonce != "" {
		result.Nonce = &block.Nonce
	}

	// Number
	if block.Number != "" {
		result.Number = hexToBigInt(block.Number)
	}

	// Size
	if block.Size != "" {
		result.Size = hexToBigInt(block.Size)
	}

	// Timestamp
	if block.Timestamp != "" {
		result.Timestamp = hexToBigInt(block.Timestamp)
	}

	// Total difficulty
	if block.TotalDifficulty != "" {
		result.TotalDifficulty = hexToBigInt(block.TotalDifficulty)
	}

	// Transactions - can be either hashes (strings) or full transaction objects
	if len(block.Transactions) > 0 {
		result.Transactions = formatBlockTransactions(block.Transactions)
	}

	return result
}

// formatBlockTransactions formats block transactions.
// Transactions can be either transaction hashes (strings) or full transaction objects.
func formatBlockTransactions(txs []any) []any {
	result := make([]any, len(txs))
	for i, tx := range txs {
		switch v := tx.(type) {
		case string:
			// Transaction hash
			result[i] = v
		case map[string]any:
			// Full transaction object - convert to RpcTransaction and format
			rpcTx := mapToRpcTransaction(v)
			result[i] = FormatTransaction(rpcTx)
		default:
			result[i] = tx
		}
	}
	return result
}

// mapToRpcTransaction converts a map to RpcTransaction.
func mapToRpcTransaction(m map[string]any) RpcTransaction {
	tx := RpcTransaction{}
	
	if v, ok := m["blockHash"].(string); ok {
		tx.BlockHash = v
	}
	if v, ok := m["blockNumber"].(string); ok {
		tx.BlockNumber = v
	}
	if v, ok := m["chainId"].(string); ok {
		tx.ChainID = v
	}
	if v, ok := m["from"].(string); ok {
		tx.From = v
	}
	if v, ok := m["gas"].(string); ok {
		tx.Gas = v
	}
	if v, ok := m["gasPrice"].(string); ok {
		tx.GasPrice = v
	}
	if v, ok := m["hash"].(string); ok {
		tx.Hash = v
	}
	if v, ok := m["input"].(string); ok {
		tx.Input = v
	}
	if v, ok := m["maxFeePerBlobGas"].(string); ok {
		tx.MaxFeePerBlobGas = v
	}
	if v, ok := m["maxFeePerGas"].(string); ok {
		tx.MaxFeePerGas = v
	}
	if v, ok := m["maxPriorityFeePerGas"].(string); ok {
		tx.MaxPriorityFeePerGas = v
	}
	if v, ok := m["nonce"].(string); ok {
		tx.Nonce = v
	}
	if v, ok := m["r"].(string); ok {
		tx.R = v
	}
	if v, ok := m["s"].(string); ok {
		tx.S = v
	}
	if v, ok := m["to"].(string); ok {
		tx.To = v
	}
	if v, ok := m["transactionIndex"].(string); ok {
		tx.TransactionIndex = v
	}
	if v, ok := m["type"].(string); ok {
		tx.Type = v
	}
	if v, ok := m["v"].(string); ok {
		tx.V = v
	}
	if v, ok := m["value"].(string); ok {
		tx.Value = v
	}
	if v, ok := m["yParity"].(string); ok {
		tx.YParity = v
	}
	
	return tx
}

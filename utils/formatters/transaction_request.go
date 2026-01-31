package formatters

import (
	"encoding/hex"
	"math/big"
)

// FormatTransactionRequest formats a TransactionRequest into RPC format.
//
// Example:
//
//	request := TransactionRequest{
//		To:    "0x...",
//		Value: big.NewInt(1000000000000000000),
//		Gas:   big.NewInt(21000),
//	}
//	rpcRequest := FormatTransactionRequest(request)
func FormatTransactionRequest(request TransactionRequest) RpcTransactionRequest {
	result := RpcTransactionRequest{
		AccessList: request.AccessList,
		Data:       request.Data,
		From:       request.From,
		To:         request.To,
	}

	// Authorization list
	if len(request.AuthorizationList) > 0 {
		result.AuthorizationList = formatRequestAuthorizationList(request.AuthorizationList)
	}

	// Blob versioned hashes
	if len(request.BlobVersionedHashes) > 0 {
		result.BlobVersionedHashes = request.BlobVersionedHashes
	}

	// Blobs - can be []byte or []string
	if len(request.Blobs) > 0 {
		result.Blobs = formatBlobs(request.Blobs)
	}

	// Gas
	if request.Gas != nil {
		result.Gas = bigIntToHex(request.Gas)
	}

	// Gas price
	if request.GasPrice != nil {
		result.GasPrice = bigIntToHex(request.GasPrice)
	}

	// Max fee per blob gas
	if request.MaxFeePerBlobGas != nil {
		result.MaxFeePerBlobGas = bigIntToHex(request.MaxFeePerBlobGas)
	}

	// Max fee per gas
	if request.MaxFeePerGas != nil {
		result.MaxFeePerGas = bigIntToHex(request.MaxFeePerGas)
	}

	// Max priority fee per gas
	if request.MaxPriorityFeePerGas != nil {
		result.MaxPriorityFeePerGas = bigIntToHex(request.MaxPriorityFeePerGas)
	}

	// Nonce
	if request.Nonce != nil {
		result.Nonce = intToHex(*request.Nonce)
	}

	// Type
	if request.Type != "" {
		if txType, ok := RpcTransactionType[request.Type]; ok {
			result.Type = txType
		}
	}

	// Value
	if request.Value != nil {
		result.Value = bigIntToHex(request.Value)
	}

	return result
}

// formatBlobs converts blobs to hex strings.
func formatBlobs(blobs []any) []string {
	result := make([]string, len(blobs))
	for i, blob := range blobs {
		switch v := blob.(type) {
		case string:
			result[i] = v
		case []byte:
			result[i] = "0x" + hex.EncodeToString(v)
		default:
			// Skip unknown types
		}
	}
	return result
}

// formatRequestAuthorizationList formats authorization list for RPC requests.
func formatRequestAuthorizationList(authList []any) []any {
	result := make([]any, 0, len(authList))
	for _, auth := range authList {
		if authMap, ok := auth.(map[string]any); ok {
			rpcAuth := make(map[string]any)

			if addr, ok := authMap["address"].(string); ok {
				rpcAuth["address"] = addr
			}
			if chainID, ok := authMap["chainId"].(int); ok {
				rpcAuth["chainId"] = intToHex(chainID)
			}
			if nonce, ok := authMap["nonce"].(int); ok {
				rpcAuth["nonce"] = intToHex(nonce)
			}
			if r, ok := authMap["r"].(string); ok {
				rpcAuth["r"] = r
			} else if r, ok := authMap["r"].(*big.Int); ok {
				rpcAuth["r"] = bigIntToHex(r)
			}
			if s, ok := authMap["s"].(string); ok {
				rpcAuth["s"] = s
			} else if s, ok := authMap["s"].(*big.Int); ok {
				rpcAuth["s"] = bigIntToHex(s)
			}
			if yParity, ok := authMap["yParity"].(int); ok {
				rpcAuth["yParity"] = intToHex(yParity)
			}
			if v, ok := authMap["v"].(int); ok {
				if _, hasYParity := rpcAuth["yParity"]; !hasYParity {
					rpcAuth["v"] = intToHex(v)
				}
			}

			result = append(result, rpcAuth)
		}
	}
	return result
}

// bigIntToHex converts a big.Int to a hex string.
func bigIntToHex(n *big.Int) string {
	if n == nil {
		return "0x0"
	}
	return "0x" + n.Text(16)
}

// intToHex converts an int to a hex string.
func intToHex(n int) string {
	return "0x" + big.NewInt(int64(n)).Text(16)
}

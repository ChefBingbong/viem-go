// Package blockoverride provides utilities for block override serialization.
package blockoverride

import (
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/types"
)

// SerializeBlockOverrides converts BlockOverrides to the RPC format.
// All numeric values are converted to hex strings.
func SerializeBlockOverrides(overrides *types.BlockOverrides) *types.RpcBlockOverrides {
	if overrides == nil {
		return nil
	}

	result := &types.RpcBlockOverrides{}

	if overrides.Number != nil {
		result.Number = hexutil.EncodeUint64(*overrides.Number)
	}

	if overrides.Difficulty != nil {
		result.Difficulty = hexutil.EncodeBig(overrides.Difficulty)
	}

	if overrides.Time != nil {
		result.Time = hexutil.EncodeUint64(*overrides.Time)
	}

	if overrides.GasLimit != nil {
		result.GasLimit = hexutil.EncodeUint64(*overrides.GasLimit)
	}

	if overrides.Coinbase != nil {
		result.Coinbase = overrides.Coinbase.Hex()
	}

	if overrides.Random != nil {
		result.Random = overrides.Random.Hex()
	}

	if overrides.BaseFeePerGas != nil {
		result.BaseFeePerGas = hexutil.EncodeBig(overrides.BaseFeePerGas)
	}

	if overrides.BlobBaseFee != nil {
		result.BlobBaseFee = hexutil.EncodeBig(overrides.BlobBaseFee)
	}

	return result
}

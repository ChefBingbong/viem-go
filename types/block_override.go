package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// BlockOverrides contains block-level overrides for eth_call.
// These allow simulating calls under different block conditions.
//
// Example:
//
//	overrides := types.BlockOverrides{
//	    BaseFeePerGas: big.NewInt(1e9),
//	    GasLimit:      ptr(uint64(30000000)),
//	}
type BlockOverrides struct {
	// Number overrides the block number.
	Number *uint64 `json:"number,omitempty"`

	// Difficulty overrides the block difficulty (pre-merge).
	Difficulty *big.Int `json:"difficulty,omitempty"`

	// Time overrides the block timestamp.
	Time *uint64 `json:"time,omitempty"`

	// GasLimit overrides the block gas limit.
	GasLimit *uint64 `json:"gasLimit,omitempty"`

	// Coinbase overrides the block coinbase/miner address.
	Coinbase *common.Address `json:"coinbase,omitempty"`

	// Random overrides the block random value (post-merge prevrandao).
	Random *common.Hash `json:"random,omitempty"`

	// BaseFeePerGas overrides the block base fee (EIP-1559).
	BaseFeePerGas *big.Int `json:"baseFeePerGas,omitempty"`

	// BlobBaseFee overrides the blob base fee (EIP-4844).
	BlobBaseFee *big.Int `json:"blobBaseFee,omitempty"`
}

// RpcBlockOverrides is the RPC format for block overrides.
// All numeric values are hex-encoded strings.
type RpcBlockOverrides struct {
	Number        string `json:"number,omitempty"`
	Difficulty    string `json:"difficulty,omitempty"`
	Time          string `json:"time,omitempty"`
	GasLimit      string `json:"gasLimit,omitempty"`
	Coinbase      string `json:"coinbase,omitempty"`
	Random        string `json:"random,omitempty"`
	BaseFeePerGas string `json:"baseFeePerGas,omitempty"`
	BlobBaseFee   string `json:"blobBaseFee,omitempty"`
}

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// StateMapping represents a mapping of storage slots to values.
// Each entry maps a 32-byte slot to a 32-byte value.
type StateMapping []StateMappingEntry

// StateMappingEntry represents a single slot-value pair in a state mapping.
type StateMappingEntry struct {
	// Slot is the storage slot (must be 32 bytes / 66 hex chars with 0x prefix).
	Slot common.Hash `json:"slot"`
	// Value is the storage value (must be 32 bytes / 66 hex chars with 0x prefix).
	Value common.Hash `json:"value"`
}

// StateOverrideAccount represents state overrides for a single account.
// This is used with eth_call to simulate different account states.
type StateOverrideAccount struct {
	// Balance overrides the account balance.
	Balance *big.Int `json:"balance,omitempty"`

	// Nonce overrides the account nonce.
	Nonce *uint64 `json:"nonce,omitempty"`

	// Code overrides the account bytecode.
	Code []byte `json:"code,omitempty"`

	// State replaces the entire account storage.
	// Mutually exclusive with StateDiff.
	State StateMapping `json:"state,omitempty"`

	// StateDiff overrides specific storage slots.
	// Mutually exclusive with State.
	StateDiff StateMapping `json:"stateDiff,omitempty"`
}

// StateOverride maps addresses to their state overrides.
// This allows simulating calls with modified account states.
//
// Example:
//
//	override := types.StateOverride{
//	    contractAddr: {
//	        Balance: big.NewInt(1e18),
//	        StateDiff: types.StateMapping{
//	            {Slot: slot, Value: value},
//	        },
//	    },
//	}
type StateOverride map[common.Address]StateOverrideAccount

// RpcStateMapping is the RPC format for state mapping (slot -> value map).
type RpcStateMapping map[string]string

// RpcAccountStateOverride is the RPC format for account state overrides.
type RpcAccountStateOverride struct {
	Balance   string          `json:"balance,omitempty"`
	Nonce     string          `json:"nonce,omitempty"`
	Code      string          `json:"code,omitempty"`
	State     RpcStateMapping `json:"state,omitempty"`
	StateDiff RpcStateMapping `json:"stateDiff,omitempty"`
}

// RpcStateOverride is the RPC format for state overrides (address -> account overrides).
type RpcStateOverride map[string]RpcAccountStateOverride

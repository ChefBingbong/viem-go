// Package stateoverride provides utilities for state override serialization.
package stateoverride

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/address"
)

// ErrInvalidBytesLength is returned when a hex value has an invalid length.
type ErrInvalidBytesLength struct {
	Size       int
	TargetSize int
	Type       string
}

func (e *ErrInvalidBytesLength) Error() string {
	return fmt.Sprintf("invalid %s length: got %d, expected %d", e.Type, e.Size, e.TargetSize)
}

// ErrAccountStateConflict is returned when the same address appears multiple times in state override.
type ErrAccountStateConflict struct {
	Address common.Address
}

func (e *ErrAccountStateConflict) Error() string {
	return fmt.Sprintf("state override conflict: address %s specified multiple times", e.Address.Hex())
}

// ErrStateAssignmentConflict is returned when both state and stateDiff are specified.
type ErrStateAssignmentConflict struct{}

func (e *ErrStateAssignmentConflict) Error() string {
	return "state override conflict: cannot specify both 'state' and 'stateDiff'"
}

// ErrInvalidAddress is returned when an address is invalid.
type ErrInvalidAddress struct {
	Address string
}

func (e *ErrInvalidAddress) Error() string {
	return fmt.Sprintf("invalid address: %s", e.Address)
}

// SerializeStateMapping converts a StateMapping to the RPC format.
// Each slot and value must be exactly 66 characters (including 0x prefix).
func SerializeStateMapping(stateMapping types.StateMapping) (types.RpcStateMapping, error) {
	if stateMapping == nil || len(stateMapping) == 0 {
		return nil, nil
	}

	result := make(types.RpcStateMapping, len(stateMapping))
	for _, entry := range stateMapping {
		slotHex := entry.Slot.Hex()
		valueHex := entry.Value.Hex()

		// Validate slot length (32 bytes = 64 hex chars + 0x prefix = 66)
		if len(slotHex) != 66 {
			return nil, &ErrInvalidBytesLength{
				Size:       len(slotHex),
				TargetSize: 66,
				Type:       "hex slot",
			}
		}

		// Validate value length
		if len(valueHex) != 66 {
			return nil, &ErrInvalidBytesLength{
				Size:       len(valueHex),
				TargetSize: 66,
				Type:       "hex value",
			}
		}

		result[slotHex] = valueHex
	}

	return result, nil
}

// SerializeAccountStateOverride converts a StateOverrideAccount to the RPC format.
func SerializeAccountStateOverride(account types.StateOverrideAccount) (types.RpcAccountStateOverride, error) {
	result := types.RpcAccountStateOverride{}

	if account.Code != nil {
		result.Code = hexutil.Encode(account.Code)
	}

	if account.Balance != nil {
		result.Balance = hexutil.EncodeBig(account.Balance)
	}

	if account.Nonce != nil {
		result.Nonce = hexutil.EncodeUint64(*account.Nonce)
	}

	if account.State != nil {
		state, err := SerializeStateMapping(account.State)
		if err != nil {
			return types.RpcAccountStateOverride{}, err
		}
		result.State = state
	}

	if account.StateDiff != nil {
		if result.State != nil {
			return types.RpcAccountStateOverride{}, &ErrStateAssignmentConflict{}
		}
		stateDiff, err := SerializeStateMapping(account.StateDiff)
		if err != nil {
			return types.RpcAccountStateOverride{}, err
		}
		result.StateDiff = stateDiff
	}

	return result, nil
}

// SerializeStateOverride converts a StateOverride to the RPC format.
// It validates addresses and checks for duplicate addresses.
func SerializeStateOverride(stateOverride types.StateOverride) (types.RpcStateOverride, error) {
	if stateOverride == nil {
		return nil, nil
	}

	result := make(types.RpcStateOverride, len(stateOverride))

	for addr, accountState := range stateOverride {
		addrHex := addr.Hex()

		// Validate address
		if !address.IsAddress(addrHex) {
			return nil, &ErrInvalidAddress{Address: addrHex}
		}

		// Check for duplicate (shouldn't happen with map, but defensive)
		if _, exists := result[addrHex]; exists {
			return nil, &ErrAccountStateConflict{Address: addr}
		}

		// Serialize account state
		serialized, err := SerializeAccountStateOverride(accountState)
		if err != nil {
			return nil, err
		}

		result[addrHex] = serialized
	}

	return result, nil
}

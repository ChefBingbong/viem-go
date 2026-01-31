package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// StateMutability represents the state mutability of a function.
type StateMutability uint8

const (
	// StateMutabilityPure functions promise not to read from or modify the state.
	StateMutabilityPure StateMutability = iota
	// StateMutabilityView functions promise not to modify the state.
	StateMutabilityView
	// StateMutabilityNonPayable functions may read and modify the state but cannot receive Ether.
	StateMutabilityNonPayable
	// StateMutabilityPayable functions may read and modify the state and can receive Ether.
	StateMutabilityPayable
)

// String returns the string representation of StateMutability.
func (s StateMutability) String() string {
	switch s {
	case StateMutabilityPure:
		return "pure"
	case StateMutabilityView:
		return "view"
	case StateMutabilityNonPayable:
		return "nonpayable"
	case StateMutabilityPayable:
		return "payable"
	default:
		return "unknown"
	}
}

// ParseStateMutability parses a string into a StateMutability.
func ParseStateMutability(s string) StateMutability {
	switch s {
	case "pure":
		return StateMutabilityPure
	case "view":
		return StateMutabilityView
	case "payable":
		return StateMutabilityPayable
	default:
		return StateMutabilityNonPayable
	}
}

// IsReadOnly returns true if the function is read-only (pure or view).
func (s StateMutability) IsReadOnly() bool {
	return s == StateMutabilityPure || s == StateMutabilityView
}

// ABIFunction represents an ABI function definition.
type ABIFunction struct {
	Name            string
	Inputs          []ABIParameter
	Outputs         []ABIParameter
	StateMutability StateMutability
	Selector        [4]byte
	Signature       string
}

// IsReadOnly returns true if the function is read-only (pure or view).
func (f *ABIFunction) IsReadOnly() bool {
	return f.StateMutability.IsReadOnly()
}

// ABIEvent represents an ABI event definition.
type ABIEvent struct {
	Name      string
	Inputs    []ABIParameter
	Anonymous bool
	Topic     common.Hash
	Signature string
}

// ABIError represents an ABI error definition.
type ABIError struct {
	Name      string
	Inputs    []ABIParameter
	Selector  [4]byte
	Signature string
}

// ABIParameter represents a function/event parameter.
type ABIParameter struct {
	Name       string
	Type       string
	Indexed    bool
	Components []ABIParameter
}

// ABIConstructor represents the contract constructor.
type ABIConstructor struct {
	Inputs          []ABIParameter
	StateMutability StateMutability
}

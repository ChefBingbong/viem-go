package abi

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// StateMutability represents the state mutability of a function.
type StateMutability uint8

const (
	// Pure functions promise not to read from or modify the state.
	Pure StateMutability = iota
	// View functions promise not to modify the state.
	View
	// NonPayable functions may read and modify the state but cannot receive Ether.
	NonPayable
	// Payable functions may read and modify the state and can receive Ether.
	Payable
)

// String returns the string representation of StateMutability.
func (s StateMutability) String() string {
	switch s {
	case Pure:
		return "pure"
	case View:
		return "view"
	case NonPayable:
		return "nonpayable"
	case Payable:
		return "payable"
	default:
		return "unknown"
	}
}

// ParseStateMutability parses a string into a StateMutability.
func ParseStateMutability(s string) StateMutability {
	switch s {
	case "pure":
		return Pure
	case "view":
		return View
	case "payable":
		return Payable
	default:
		return NonPayable
	}
}

// IsReadOnly returns true if the function is read-only (pure or view).
func (s StateMutability) IsReadOnly() bool {
	return s == Pure || s == View
}

// ABI represents a parsed Ethereum ABI that wraps go-ethereum's ABI.
type ABI struct {
	gethABI   abi.ABI
	raw       []byte
	Functions map[string]Function
	Events    map[string]Event
	Errors    map[string]Error
}

// GethABI returns the underlying go-ethereum ABI.
func (a *ABI) GethABI() *abi.ABI {
	return &a.gethABI
}

// Raw returns the original JSON ABI bytes.
func (a *ABI) Raw() []byte {
	return a.raw
}

// Function represents an ABI function definition.
type Function struct {
	Name            string
	Inputs          []Parameter
	Outputs         []Parameter
	StateMutability StateMutability
	Selector        [4]byte
	Signature       string
}

// IsReadOnly returns true if the function is read-only (pure or view).
func (f *Function) IsReadOnly() bool {
	return f.StateMutability.IsReadOnly()
}

// Event represents an ABI event definition.
type Event struct {
	Name      string
	Inputs    []Parameter
	Anonymous bool
	Topic     common.Hash
	Signature string
}

// Error represents an ABI error definition.
type Error struct {
	Name      string
	Inputs    []Parameter
	Selector  [4]byte
	Signature string
}

// Parameter represents a function/event parameter.
type Parameter struct {
	Name       string
	Type       string
	Indexed    bool
	Components []Parameter
}

// Constructor represents the contract constructor.
type Constructor struct {
	Inputs          []Parameter
	StateMutability StateMutability
}

// convertGethArgumentsToParameters converts go-ethereum Arguments to our Parameter type.
func convertGethArgumentsToParameters(args abi.Arguments) []Parameter {
	params := make([]Parameter, len(args))
	for i, arg := range args {
		params[i] = Parameter{
			Name:       arg.Name,
			Type:       arg.Type.String(),
			Indexed:    arg.Indexed,
			Components: convertGethTupleComponents(arg.Type),
		}
	}
	return params
}

// convertGethTupleComponents converts tuple components recursively.
func convertGethTupleComponents(typ abi.Type) []Parameter {
	if typ.T != abi.TupleTy {
		return nil
	}
	components := make([]Parameter, len(typ.TupleElems))
	for i, elem := range typ.TupleElems {
		name := ""
		if i < len(typ.TupleRawNames) {
			name = typ.TupleRawNames[i]
		}
		components[i] = Parameter{
			Name:       name,
			Type:       elem.String(),
			Components: convertGethTupleComponents(*elem),
		}
	}
	return components
}

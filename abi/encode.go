package abi

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EncodeCall encodes a function call with the given method name and arguments.
// Returns the full calldata including the 4-byte function selector.
func (a *ABI) EncodeCall(method string, args ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}

	packed, err := a.gethABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode call for method %q: %w", method, err)
	}

	// gethABI.Pack already includes the selector
	_ = m // method found
	return packed, nil
}

// Pack is an alias for EncodeCall for compatibility with go-ethereum naming.
func (a *ABI) Pack(method string, args ...any) ([]byte, error) {
	return a.EncodeCall(method, args...)
}

// EncodeCallWithSelector encodes function arguments and prepends the given selector.
// Useful when you have a custom selector or are encoding for a different function.
func (a *ABI) EncodeCallWithSelector(selector [4]byte, method string, args ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}

	// Pack arguments only (without selector)
	packed, err := m.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode arguments for method %q: %w", method, err)
	}

	// Prepend custom selector
	result := make([]byte, 4+len(packed))
	copy(result[:4], selector[:])
	copy(result[4:], packed)
	return result, nil
}

// EncodeArgs encodes only the arguments without the function selector.
// Useful for encoding constructor arguments or when you need raw argument encoding.
func (a *ABI) EncodeArgs(method string, args ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}

	packed, err := m.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode arguments for method %q: %w", method, err)
	}

	return packed, nil
}

// EncodeConstructor encodes constructor arguments.
func (a *ABI) EncodeConstructor(args ...any) ([]byte, error) {
	if a.gethABI.Constructor.Inputs == nil {
		if len(args) > 0 {
			return nil, fmt.Errorf("constructor takes no arguments but %d provided", len(args))
		}
		return nil, nil
	}

	packed, err := a.gethABI.Constructor.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode constructor arguments: %w", err)
	}

	return packed, nil
}

// EncodeError encodes an error with the given name and arguments.
func (a *ABI) EncodeError(name string, args ...any) ([]byte, error) {
	e, ok := a.gethABI.Errors[name]
	if !ok {
		return nil, fmt.Errorf("error %q not found in ABI", name)
	}

	packed, err := e.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode error %q: %w", name, err)
	}

	// Prepend error selector
	result := make([]byte, 4+len(packed))
	copy(result[:4], e.ID[:4])
	copy(result[4:], packed)
	return result, nil
}

// EncodeEventTopics encodes indexed event parameters as topics.
// Non-indexed parameters should be passed as nil.
// Note: This is a simplified implementation that handles common cases.
func (a *ABI) EncodeEventTopics(name string, args ...any) ([][]byte, error) {
	e, ok := a.gethABI.Events[name]
	if !ok {
		return nil, fmt.Errorf("event %q not found in ABI", name)
	}

	// First topic is always the event signature (unless anonymous)
	var topics [][]byte
	if !e.Anonymous {
		topics = append(topics, e.ID.Bytes())
	}

	// Encode indexed parameters
	indexedCount := 0
	for _, input := range e.Inputs {
		if input.Indexed {
			indexedCount++
		}
	}

	if len(args) != indexedCount {
		return nil, fmt.Errorf("event %q expects %d indexed arguments, got %d", name, indexedCount, len(args))
	}

	argIndex := 0
	for _, input := range e.Inputs {
		if input.Indexed {
			if args[argIndex] != nil {
				// Encode the indexed argument as a topic
				topic, err := encodeIndexedArg(input.Type.String(), args[argIndex])
				if err != nil {
					return nil, fmt.Errorf("failed to encode indexed argument %q: %w", input.Name, err)
				}
				topics = append(topics, topic)
			} else {
				// nil means "match any" - add empty topic
				topics = append(topics, nil)
			}
			argIndex++
		}
	}

	return topics, nil
}

// encodeIndexedArg encodes a value as an indexed topic (32 bytes).
func encodeIndexedArg(typeStr string, value any) ([]byte, error) {
	topic := make([]byte, 32)

	switch v := value.(type) {
	case common.Address:
		copy(topic[12:], v.Bytes())
	case *big.Int:
		b := v.Bytes()
		copy(topic[32-len(b):], b)
	case bool:
		if v {
			topic[31] = 1
		}
	case []byte:
		// For bytes/string, use keccak256 hash
		hash := crypto.Keccak256(v)
		copy(topic, hash)
	case string:
		hash := crypto.Keccak256([]byte(v))
		copy(topic, hash)
	default:
		return nil, fmt.Errorf("unsupported type for indexed argument: %T", value)
	}

	return topic, nil
}

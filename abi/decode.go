package abi

import (
	"fmt"
	"math/big"

	gethABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// DecodeReturn decodes the return data from a function call.
// Returns a slice of decoded values matching the function's output parameters.
func (a *ABI) DecodeReturn(method string, data []byte) ([]any, error) {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}

	if len(data) == 0 {
		if len(m.Outputs) > 0 {
			return nil, fmt.Errorf("expected return data for method %q but got empty", method)
		}
		return nil, nil
	}

	unpacked, err := m.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode return for method %q: %w", method, err)
	}

	return unpacked, nil
}

// Unpack is an alias for DecodeReturn for compatibility with go-ethereum naming.
func (a *ABI) Unpack(method string, data []byte) ([]any, error) {
	return a.DecodeReturn(method, data)
}

// DecodeReturnInto decodes the return data into the provided struct or variables.
// The output parameter must be a pointer to a struct or a slice of pointers.
func (a *ABI) DecodeReturnInto(method string, data []byte, output any) error {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return fmt.Errorf("method %q not found in ABI", method)
	}

	if len(data) == 0 {
		if len(m.Outputs) > 0 {
			return fmt.Errorf("expected return data for method %q but got empty", method)
		}
		return nil
	}

	return a.gethABI.UnpackIntoInterface(output, method, data)
}

// UnpackIntoInterface is an alias for DecodeReturnInto.
func (a *ABI) UnpackIntoInterface(method string, data []byte, output any) error {
	return a.DecodeReturnInto(method, data, output)
}

// DecodeArgs decodes the input arguments from calldata (including the 4-byte selector).
func (a *ABI) DecodeArgs(method string, data []byte) ([]any, error) {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("calldata too short: expected at least 4 bytes, got %d", len(data))
	}

	// Skip the 4-byte selector
	unpacked, err := m.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode args for method %q: %w", method, err)
	}

	return unpacked, nil
}

// DecodeArgsFromData decodes the input arguments from raw data (without selector).
func (a *ABI) DecodeArgsFromData(method string, data []byte) ([]any, error) {
	m, ok := a.gethABI.Methods[method]
	if !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}

	unpacked, err := m.Inputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode args for method %q: %w", method, err)
	}

	return unpacked, nil
}

// DecodeEvent decodes event data and topics into a map of parameter names to values.
func (a *ABI) DecodeEvent(name string, topics []common.Hash, data []byte) (map[string]any, error) {
	e, ok := a.gethABI.Events[name]
	if !ok {
		return nil, fmt.Errorf("event %q not found in ABI", name)
	}

	result := make(map[string]any)

	// Separate indexed and non-indexed inputs
	var indexedInputs, nonIndexedInputs []int
	for i, input := range e.Inputs {
		if input.Indexed {
			indexedInputs = append(indexedInputs, i)
		} else {
			nonIndexedInputs = append(nonIndexedInputs, i)
		}
	}

	// Decode indexed topics
	// Skip the first topic if not anonymous (it's the event signature)
	topicOffset := 0
	if !e.Anonymous {
		topicOffset = 1
	}

	for i, idx := range indexedInputs {
		topicIdx := topicOffset + i
		if topicIdx >= len(topics) {
			return nil, fmt.Errorf("not enough topics for event %q: expected %d, got %d", name, len(indexedInputs)+topicOffset, len(topics))
		}

		input := e.Inputs[idx]
		// Indexed dynamic types (string, bytes, arrays) are hashed, so we can only return the hash
		// Check for dynamic types by string comparison
		typeStr := input.Type.String()
		if typeStr == "string" || typeStr == "bytes" || input.Type.T == 4 { // 4 is SliceTy
			result[input.Name] = topics[topicIdx]
		} else {
			// For fixed-size types, decode the topic as the value
			// The topic contains the value directly for basic types
			result[input.Name] = decodeIndexedTopic(input.Type, topics[topicIdx])
		}
	}

	// Decode non-indexed data
	if len(nonIndexedInputs) > 0 && len(data) > 0 {
		// Build arguments for non-indexed inputs only
		unpacked, err := e.Inputs.UnpackValues(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode event data for %q: %w", name, err)
		}

		// Map unpacked values to non-indexed inputs
		unpackedIdx := 0
		for _, idx := range nonIndexedInputs {
			if unpackedIdx < len(unpacked) {
				result[e.Inputs[idx].Name] = unpacked[unpackedIdx]
				unpackedIdx++
			}
		}
	}

	return result, nil
}

// DecodeEventIntoStruct decodes event data into the provided struct.
func (a *ABI) DecodeEventIntoStruct(name string, topics []common.Hash, data []byte, output any) error {
	_, ok := a.gethABI.Events[name]
	if !ok {
		return fmt.Errorf("event %q not found in ABI", name)
	}

	// Use go-ethereum's built-in unpacking
	return a.gethABI.UnpackIntoInterface(output, name, data)
}

// decodeIndexedTopic decodes an indexed topic value based on its type.
func decodeIndexedTopic(typ gethABI.Type, topic common.Hash) any {
	switch typ.T {
	case gethABI.AddressTy:
		return common.BytesToAddress(topic.Bytes())
	case gethABI.BoolTy:
		return topic[31] != 0
	case gethABI.IntTy, gethABI.UintTy:
		return new(big.Int).SetBytes(topic.Bytes())
	case gethABI.FixedBytesTy:
		// Return the topic as-is for fixed bytes
		return topic
	default:
		// For other types, return the hash
		return topic
	}
}

// DecodeError decodes an error return value.
func (a *ABI) DecodeError(data []byte) (string, []any, error) {
	if len(data) < 4 {
		return "", nil, fmt.Errorf("error data too short: expected at least 4 bytes, got %d", len(data))
	}

	// Check for standard Error(string) selector: 0x08c379a0
	if data[0] == 0x08 && data[1] == 0xc3 && data[2] == 0x79 && data[3] == 0xa0 {
		// Standard revert with error message
		if len(data) > 4 {
			// Decode the string
			reason, err := decodeString(data[4:])
			if err != nil {
				return "Error", nil, nil
			}
			return "Error", []any{reason}, nil
		}
		return "Error", nil, nil
	}

	// Check for Panic(uint256) selector: 0x4e487b71
	if data[0] == 0x4e && data[1] == 0x48 && data[2] == 0x7b && data[3] == 0x71 {
		if len(data) >= 36 {
			code := new(big.Int).SetBytes(data[4:36])
			return "Panic", []any{code}, nil
		}
		return "Panic", nil, nil
	}

	// Try to match against custom errors in ABI
	var selector [4]byte
	copy(selector[:], data[:4])

	for _, e := range a.gethABI.Errors {
		var errSelector [4]byte
		copy(errSelector[:], e.ID[:4])
		if errSelector == selector {
			unpacked, err := e.Inputs.Unpack(data[4:])
			if err != nil {
				return e.Name, nil, fmt.Errorf("failed to decode error %q: %w", e.Name, err)
			}
			return e.Name, unpacked, nil
		}
	}

	return "", nil, fmt.Errorf("unknown error selector: %x", selector)
}

// decodeString decodes an ABI-encoded string.
func decodeString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", fmt.Errorf("data too short for string")
	}

	// First 32 bytes is the offset (should be 32 for a single string)
	// Next 32 bytes is the length
	length := new(big.Int).SetBytes(data[32:64]).Uint64()

	if uint64(len(data)) < 64+length {
		return "", fmt.Errorf("data too short for string content")
	}

	return string(data[64 : 64+length]), nil
}

// DecodeCalldata decodes calldata and returns the method name and arguments.
func (a *ABI) DecodeCalldata(data []byte) (string, []any, error) {
	if len(data) < 4 {
		return "", nil, fmt.Errorf("calldata too short: expected at least 4 bytes, got %d", len(data))
	}

	var selector [4]byte
	copy(selector[:], data[:4])

	// Find matching method
	for _, m := range a.gethABI.Methods {
		var methodSelector [4]byte
		copy(methodSelector[:], m.ID)
		if methodSelector == selector {
			args, err := m.Inputs.Unpack(data[4:])
			if err != nil {
				return m.Name, nil, fmt.Errorf("failed to decode args for method %q: %w", m.Name, err)
			}
			return m.Name, args, nil
		}
	}

	return "", nil, fmt.Errorf("unknown function selector: %x", selector)
}

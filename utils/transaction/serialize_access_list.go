package transaction

import (
	"fmt"
	"strings"
)

// SerializeAccessList serializes an EIP-2930 access list for RLP encoding.
//
// Example:
//
//	serialized, err := SerializeAccessList(AccessList{
//		{
//			Address: "0x1234567890123456789012345678901234567890",
//			StorageKeys: []string{
//				"0x0000000000000000000000000000000000000000000000000000000000000001",
//			},
//		},
//	})
func SerializeAccessList(accessList AccessList) ([]any, error) {
	if len(accessList) == 0 {
		return []any{}, nil
	}

	result := make([]any, 0, len(accessList))

	for _, item := range accessList {
		// Validate address
		if !isValidAddress(item.Address) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidAddress, item.Address)
		}

		// Validate storage keys (must be 32 bytes = 64 hex chars + 0x prefix)
		for _, key := range item.StorageKeys {
			keyClean := strings.TrimPrefix(key, "0x")
			keyClean = strings.TrimPrefix(keyClean, "0X")
			if len(keyClean) != 64 {
				return nil, fmt.Errorf("%w: %s (size: %d)", ErrInvalidStorageKeySize, key, len(keyClean)/2)
			}
		}

		// Convert to RLP-compatible format: [address, [storageKey1, storageKey2, ...]]
		storageKeys := make([]any, len(item.StorageKeys))
		for i, key := range item.StorageKeys {
			storageKeys[i] = key
		}

		result = append(result, []any{item.Address, storageKeys})
	}

	return result, nil
}

// ParseAccessList parses an RLP-decoded access list.
func ParseAccessList(data []any) (AccessList, error) {
	if len(data) == 0 {
		return nil, nil
	}

	result := make(AccessList, 0, len(data))

	for _, item := range data {
		itemSlice, ok := item.([]any)
		if !ok || len(itemSlice) != 2 {
			continue
		}

		address, ok := itemSlice[0].(string)
		if !ok {
			continue
		}

		// Validate address
		if !isValidAddress(address) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidAddress, address)
		}

		storageKeysRaw, ok := itemSlice[1].([]any)
		if !ok {
			continue
		}

		storageKeys := make([]string, 0, len(storageKeysRaw))
		for _, keyRaw := range storageKeysRaw {
			key, ok := keyRaw.(string)
			if !ok {
				continue
			}
			// Normalize the key (trim if needed, but keep 32 bytes)
			storageKeys = append(storageKeys, normalizeStorageKey(key))
		}

		result = append(result, AccessListItem{
			Address:     address,
			StorageKeys: storageKeys,
		})
	}

	return result, nil
}

// normalizeStorageKey normalizes a storage key to 32 bytes.
func normalizeStorageKey(key string) string {
	keyClean := strings.TrimPrefix(key, "0x")
	keyClean = strings.TrimPrefix(keyClean, "0X")

	// If it's already 64 chars (32 bytes), return as-is
	if len(keyClean) == 64 {
		return "0x" + keyClean
	}

	// Pad left with zeros if shorter
	if len(keyClean) < 64 {
		return "0x" + strings.Repeat("0", 64-len(keyClean)) + keyClean
	}

	// Trim from left if longer (shouldn't happen for valid keys)
	return "0x" + keyClean[len(keyClean)-64:]
}

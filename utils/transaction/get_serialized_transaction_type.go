package transaction

import (
	"fmt"
	"strings"
)

// GetSerializedTransactionType determines the transaction type from a serialized transaction.
//
// Example:
//
//	txType, err := GetSerializedTransactionType("0x02...")
//	// txType = "eip1559"
func GetSerializedTransactionType(serializedTx string) (TransactionType, error) {
	if len(serializedTx) < 4 {
		return "", fmt.Errorf("%w: transaction too short", ErrInvalidSerializedTransactionType)
	}

	// Get the first byte (after 0x prefix)
	serializedType := strings.ToLower(serializedTx[0:4])

	switch serializedType {
	case "0x04":
		return TransactionTypeEIP7702, nil
	case "0x03":
		return TransactionTypeEIP4844, nil
	case "0x02":
		return TransactionTypeEIP1559, nil
	case "0x01":
		return TransactionTypeEIP2930, nil
	}

	// Check for legacy transaction (RLP encoded, starts with 0xc0-0xff)
	if len(serializedTx) >= 4 {
		firstByte := hexToInt(serializedTx[2:4])
		if firstByte >= 0xc0 {
			return TransactionTypeLegacy, nil
		}
	}

	return "", fmt.Errorf("%w: %s", ErrInvalidSerializedTransactionType, serializedType)
}

// hexToInt converts a 2-character hex string to an integer.
func hexToInt(s string) int {
	var result int
	for _, c := range strings.ToLower(s) {
		result <<= 4
		if c >= '0' && c <= '9' {
			result += int(c - '0')
		} else if c >= 'a' && c <= 'f' {
			result += int(c - 'a' + 10)
		}
	}
	return result
}

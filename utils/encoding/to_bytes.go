package encoding

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

// ToBytesConverter provides a fluent API for converting values to bytes.
type ToBytesConverter struct {
	value any
	opts  ToBytesOpts
}

// ToBytesOpts configures byte conversion behavior.
type ToBytesOpts struct {
	// Size pads/validates the output to this size.
	Size int
	// Signed treats the number as signed for integer conversions.
	Signed bool
}

// ToBytes creates a converter that will convert the value to bytes.
func ToBytes(value any) *ToBytesConverter {
	return &ToBytesConverter{value: value}
}

// WithSize sets the output size (pads if necessary).
func (c *ToBytesConverter) WithSize(size int) *ToBytesConverter {
	c.opts.Size = size
	return c
}

// WithSigned treats numbers as signed for conversion.
func (c *ToBytesConverter) WithSigned() *ToBytesConverter {
	c.opts.Signed = true
	return c
}

// Bytes performs the conversion.
func (c *ToBytesConverter) Bytes() ([]byte, error) {
	switch v := c.value.(type) {
	case []byte:
		return c.applySize(v)
	case string:
		if IsHex(v) {
			return c.hexToBytes(v)
		}
		return c.stringToBytes(v)
	case bool:
		return c.boolToBytes(v)
	case int:
		return c.numberToBytes(big.NewInt(int64(v)))
	case int8:
		return c.numberToBytes(big.NewInt(int64(v)))
	case int16:
		return c.numberToBytes(big.NewInt(int64(v)))
	case int32:
		return c.numberToBytes(big.NewInt(int64(v)))
	case int64:
		return c.numberToBytes(big.NewInt(v))
	case uint:
		return c.numberToBytes(new(big.Int).SetUint64(uint64(v)))
	case uint8:
		return c.numberToBytes(new(big.Int).SetUint64(uint64(v)))
	case uint16:
		return c.numberToBytes(new(big.Int).SetUint64(uint64(v)))
	case uint32:
		return c.numberToBytes(new(big.Int).SetUint64(uint64(v)))
	case uint64:
		return c.numberToBytes(new(big.Int).SetUint64(v))
	case *big.Int:
		return c.numberToBytes(v)
	default:
		return nil, errors.New("unsupported type for ToBytes")
	}
}

func (c *ToBytesConverter) applySize(b []byte) ([]byte, error) {
	if c.opts.Size > 0 {
		if len(b) > c.opts.Size {
			return nil, errors.New("size overflow")
		}
		return PadRight(b, c.opts.Size), nil
	}
	return b, nil
}

func (c *ToBytesConverter) hexToBytes(s string) ([]byte, error) {
	b, err := HexToBytes(s)
	if err != nil {
		return nil, err
	}
	return c.applySize(b)
}

func (c *ToBytesConverter) stringToBytes(s string) ([]byte, error) {
	b := []byte(s)
	return c.applySize(b)
}

func (c *ToBytesConverter) boolToBytes(v bool) ([]byte, error) {
	var b []byte
	if v {
		b = []byte{1}
	} else {
		b = []byte{0}
	}

	if c.opts.Size > 0 {
		return PadLeft(b, c.opts.Size), nil
	}
	return b, nil
}

func (c *ToBytesConverter) numberToBytes(n *big.Int) ([]byte, error) {
	if n == nil {
		return nil, errors.New("nil big.Int")
	}

	// Handle negative numbers with two's complement
	if c.opts.Signed && n.Sign() < 0 {
		if c.opts.Size == 0 {
			return nil, errors.New("signed conversion requires size")
		}
		// Two's complement: 2^(size*8) + value
		max := new(big.Int).Lsh(big.NewInt(1), uint(c.opts.Size*8))
		n = new(big.Int).Add(max, n)
	}

	b := n.Bytes()
	if len(b) == 0 {
		b = []byte{0}
	}

	if c.opts.Size > 0 {
		if len(b) > c.opts.Size {
			return nil, errors.New("integer overflow for size")
		}
		return PadLeft(b, c.opts.Size), nil
	}
	return b, nil
}

// Standalone conversion functions

// HexToBytes converts a hex string to bytes.
func HexToBytes(s string) ([]byte, error) {
	s = Strip0x(s)
	// Pad odd-length hex strings
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// StringToBytes converts a UTF-8 string to bytes.
func StringToBytes(s string) []byte {
	return []byte(s)
}

// BoolToBytes converts a boolean to bytes.
func BoolToBytes(v bool) []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

// NumberToBytes converts a number to bytes.
func NumberToBytes(n *big.Int) []byte {
	if n == nil {
		return []byte{0}
	}
	b := n.Bytes()
	if len(b) == 0 {
		return []byte{0}
	}
	return b
}

// NumberToBytesWithSize converts a number to bytes with padding.
func NumberToBytesWithSize(n *big.Int, size int, signed bool) ([]byte, error) {
	return ToBytes(n).WithSize(size).Bytes()
}

// Helper functions

// Strip0x removes the 0x prefix from a hex string.
func Strip0x(s string) string {
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		return s[2:]
	}
	return s
}

// IsHex checks if a string is a valid hex string.
func IsHex(s string) bool {
	if len(s) < 2 {
		return false
	}
	if s[0:2] != "0x" && s[0:2] != "0X" {
		return false
	}
	s = s[2:]
	if len(s) == 0 {
		return true // "0x" is valid
	}
	for _, c := range strings.ToLower(s) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// PadLeft pads bytes on the left to reach the target size.
func PadLeft(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}

// PadRight pads bytes on the right to reach the target size.
func PadRight(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded, b)
	return padded
}

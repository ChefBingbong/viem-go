package encoding

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

// ToHexConverter provides a fluent API for converting values to hex strings.
type ToHexConverter struct {
	value any
	opts  ToHexOpts
}

// ToHexOpts configures hex conversion behavior.
type ToHexOpts struct {
	// Size pads the output to this byte size.
	Size int
	// Signed treats the number as signed for integer conversions.
	Signed bool
}

// ToHex creates a converter that will convert the value to hex.
func ToHex(value any) *ToHexConverter {
	return &ToHexConverter{value: value}
}

// WithSize sets the output size in bytes (pads if necessary).
func (c *ToHexConverter) WithSize(size int) *ToHexConverter {
	c.opts.Size = size
	return c
}

// WithSigned treats numbers as signed for conversion.
func (c *ToHexConverter) WithSigned() *ToHexConverter {
	c.opts.Signed = true
	return c
}

// Hex performs the conversion.
func (c *ToHexConverter) Hex() (string, error) {
	switch v := c.value.(type) {
	case []byte:
		return c.bytesToHex(v)
	case string:
		return c.stringToHex(v)
	case bool:
		return c.boolToHex(v)
	case int:
		return c.numberToHex(big.NewInt(int64(v)))
	case int8:
		return c.numberToHex(big.NewInt(int64(v)))
	case int16:
		return c.numberToHex(big.NewInt(int64(v)))
	case int32:
		return c.numberToHex(big.NewInt(int64(v)))
	case int64:
		return c.numberToHex(big.NewInt(v))
	case uint:
		return c.numberToHex(new(big.Int).SetUint64(uint64(v)))
	case uint8:
		return c.numberToHex(new(big.Int).SetUint64(uint64(v)))
	case uint16:
		return c.numberToHex(new(big.Int).SetUint64(uint64(v)))
	case uint32:
		return c.numberToHex(new(big.Int).SetUint64(uint64(v)))
	case uint64:
		return c.numberToHex(new(big.Int).SetUint64(v))
	case *big.Int:
		return c.numberToHex(v)
	default:
		return "", errors.New("unsupported type for ToHex")
	}
}

func (c *ToHexConverter) bytesToHex(b []byte) (string, error) {
	hexStr := "0x" + hex.EncodeToString(b)
	if c.opts.Size > 0 {
		if len(b) > c.opts.Size {
			return "", errors.New("size overflow")
		}
		return PadHexRight(hexStr, c.opts.Size), nil
	}
	return hexStr, nil
}

func (c *ToHexConverter) stringToHex(s string) (string, error) {
	b := []byte(s)
	return c.bytesToHex(b)
}

func (c *ToHexConverter) boolToHex(v bool) (string, error) {
	var hexStr string
	if v {
		hexStr = "0x1"
	} else {
		hexStr = "0x0"
	}

	if c.opts.Size > 0 {
		return PadHexLeft(hexStr, c.opts.Size), nil
	}
	return hexStr, nil
}

func (c *ToHexConverter) numberToHex(n *big.Int) (string, error) {
	if n == nil {
		return "", errors.New("nil big.Int")
	}

	// Handle negative numbers with two's complement
	if c.opts.Signed && n.Sign() < 0 {
		if c.opts.Size == 0 {
			return "", errors.New("signed conversion requires size")
		}
		// Two's complement: 2^(size*8) + value
		max := new(big.Int).Lsh(big.NewInt(1), uint(c.opts.Size*8))
		n = new(big.Int).Add(max, n)
	}

	hexStr := "0x" + n.Text(16)

	if c.opts.Size > 0 {
		// Check for overflow
		hexLen := len(Strip0x(hexStr))
		if hexLen%2 != 0 {
			hexLen++
		}
		if hexLen/2 > c.opts.Size {
			return "", errors.New("integer overflow for size")
		}
		return PadHexLeft(hexStr, c.opts.Size), nil
	}

	return hexStr, nil
}

// Standalone conversion functions

// BoolToHex converts a boolean to a hex string.
func BoolToHex(v bool) string {
	if v {
		return "0x1"
	}
	return "0x0"
}

// NumberToHex converts a number to a hex string.
func NumberToHex(n *big.Int) string {
	if n == nil {
		return "0x0"
	}
	return "0x" + n.Text(16)
}

// NumberToHexWithSize converts a number to a hex string with padding.
func NumberToHexWithSize(n *big.Int, size int, signed bool) (string, error) {
	c := ToHex(n).WithSize(size)
	if signed {
		c = c.WithSigned()
	}
	return c.Hex()
}

// StringToHex converts a UTF-8 string to a hex string.
func StringToHex(s string) string {
	return "0x" + hex.EncodeToString([]byte(s))
}

// Helper functions

// PadHexLeft pads a hex string on the left to reach the target byte size.
func PadHexLeft(s string, size int) string {
	s = Strip0x(s)
	targetLen := size * 2
	if len(s) >= targetLen {
		return "0x" + s
	}
	return "0x" + strings.Repeat("0", targetLen-len(s)) + s
}

// PadHexRight pads a hex string on the right to reach the target byte size.
func PadHexRight(s string, size int) string {
	s = Strip0x(s)
	targetLen := size * 2
	if len(s) >= targetLen {
		return "0x" + s
	}
	return "0x" + s + strings.Repeat("0", targetLen-len(s))
}

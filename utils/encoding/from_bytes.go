package encoding

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
)

// ByteConverter provides a fluent API for converting bytes to other types.
type ByteConverter struct {
	data []byte
	opts ByteConverterOpts
}

// ByteConverterOpts configures byte conversion behavior.
type ByteConverterOpts struct {
	// Size validates that the byte array is exactly this size.
	Size int
	// Signed treats the bytes as a signed integer (for BigInt/Number conversion).
	Signed bool
}

// FromBytes creates a new ByteConverter from a byte slice.
func FromBytes(b []byte) *ByteConverter {
	return &ByteConverter{data: b}
}

// WithSize sets the expected size for validation.
func (c *ByteConverter) WithSize(size int) *ByteConverter {
	c.opts.Size = size
	return c
}

// WithSigned treats the bytes as signed for integer conversions.
func (c *ByteConverter) WithSigned() *ByteConverter {
	c.opts.Signed = true
	return c
}

// ToHex converts the bytes to a hex string with 0x prefix.
func (c *ByteConverter) ToHex() (string, error) {
	if c.opts.Size > 0 {
		if err := assertSize(c.data, c.opts.Size); err != nil {
			return "", err
		}
	}
	return BytesToHex(c.data), nil
}

// ToHexUnprefixed converts the bytes to a hex string without 0x prefix.
func (c *ByteConverter) ToHexUnprefixed() (string, error) {
	if c.opts.Size > 0 {
		if err := assertSize(c.data, c.opts.Size); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(c.data), nil
}

// ToBigInt converts the bytes to a *big.Int.
func (c *ByteConverter) ToBigInt() (*big.Int, error) {
	if c.opts.Size > 0 {
		if err := assertSize(c.data, c.opts.Size); err != nil {
			return nil, err
		}
	}
	return BytesToBigInt(c.data, c.opts.Signed), nil
}

// ToNumber converts the bytes to an int64.
func (c *ByteConverter) ToNumber() (int64, error) {
	if c.opts.Size > 0 {
		if err := assertSize(c.data, c.opts.Size); err != nil {
			return 0, err
		}
	}
	bi := BytesToBigInt(c.data, c.opts.Signed)
	if !bi.IsInt64() {
		return 0, errors.New("integer out of int64 range")
	}
	return bi.Int64(), nil
}

// ToUint converts the bytes to a uint64.
func (c *ByteConverter) ToUint() (uint64, error) {
	if c.opts.Size > 0 {
		if err := assertSize(c.data, c.opts.Size); err != nil {
			return 0, err
		}
	}
	bi := BytesToBigInt(c.data, false)
	if !bi.IsUint64() {
		return 0, errors.New("integer out of uint64 range")
	}
	return bi.Uint64(), nil
}

// ToBool converts the bytes to a boolean.
func (c *ByteConverter) ToBool() (bool, error) {
	data := c.data
	if c.opts.Size > 0 {
		if err := assertSize(data, c.opts.Size); err != nil {
			return false, err
		}
		data = TrimLeft(data)
	}

	if len(data) > 1 || (len(data) == 1 && data[0] > 1) {
		return false, errors.New("invalid bytes boolean")
	}
	if len(data) == 0 {
		return false, nil
	}
	return data[0] == 1, nil
}

// ToString converts the bytes to a UTF-8 string.
func (c *ByteConverter) ToString() (string, error) {
	data := c.data
	if c.opts.Size > 0 {
		if err := assertSize(data, c.opts.Size); err != nil {
			return "", err
		}
		data = TrimRight(data)
	}
	return string(data), nil
}

// ToBytes returns the underlying byte slice.
func (c *ByteConverter) ToBytes() []byte {
	return c.data
}

// Standalone conversion functions

// BytesToHex converts bytes to a hex string with 0x prefix.
func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// BytesToBigInt converts bytes to a *big.Int.
// If signed is true, treats the bytes as two's complement signed integer.
func BytesToBigInt(b []byte, signed bool) *big.Int {
	if len(b) == 0 {
		return big.NewInt(0)
	}

	result := new(big.Int).SetBytes(b)

	if signed && len(b) > 0 && b[0]&0x80 != 0 {
		// Negative number in two's complement
		// Calculate: result - 2^(len*8)
		max := new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8))
		result.Sub(result, max)
	}

	return result
}

// BytesToNumber converts bytes to an int64.
func BytesToNumber(b []byte, signed bool) (int64, error) {
	bi := BytesToBigInt(b, signed)
	if !bi.IsInt64() {
		return 0, errors.New("integer out of int64 range")
	}
	return bi.Int64(), nil
}

// BytesToUint converts bytes to a uint64.
func BytesToUint(b []byte) (uint64, error) {
	if len(b) > 8 {
		return 0, errors.New("integer out of uint64 range")
	}
	if len(b) < 8 {
		padded := make([]byte, 8)
		copy(padded[8-len(b):], b)
		b = padded
	}
	return binary.BigEndian.Uint64(b), nil
}

// BytesToBool converts bytes to a boolean.
func BytesToBool(b []byte) (bool, error) {
	trimmed := TrimLeft(b)
	if len(trimmed) > 1 || (len(trimmed) == 1 && trimmed[0] > 1) {
		return false, errors.New("invalid bytes boolean")
	}
	if len(trimmed) == 0 {
		return false, nil
	}
	return trimmed[0] == 1, nil
}

// BytesToString converts bytes to a UTF-8 string.
func BytesToString(b []byte) string {
	return string(TrimRight(b))
}

// Helper functions

func assertSize(data []byte, size int) error {
	if len(data) > size {
		return errors.New("size overflow")
	}
	return nil
}

// TrimLeft removes leading zero bytes.
func TrimLeft(b []byte) []byte {
	return bytes.TrimLeft(b, "\x00")
}

// TrimRight removes trailing zero bytes.
func TrimRight(b []byte) []byte {
	return bytes.TrimRight(b, "\x00")
}

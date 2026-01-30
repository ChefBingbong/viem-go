package encoding

import (
	"errors"
	"math/big"
	"strconv"
)

// HexConverter provides a fluent API for converting hex strings to other types.
type HexConverter struct {
	data string
	opts HexConverterOpts
}

// HexConverterOpts configures hex conversion behavior.
type HexConverterOpts struct {
	// Size validates that the hex represents exactly this many bytes.
	Size int
	// Signed treats the hex as a signed integer (for BigInt/Number conversion).
	Signed bool
}

// FromHex creates a new HexConverter from a hex string.
func FromHex(s string) *HexConverter {
	return &HexConverter{data: s}
}

// WithSize sets the expected size for validation.
func (c *HexConverter) WithSize(size int) *HexConverter {
	c.opts.Size = size
	return c
}

// WithSigned treats the hex as signed for integer conversions.
func (c *HexConverter) WithSigned() *HexConverter {
	c.opts.Signed = true
	return c
}

// ToBytes converts the hex string to bytes.
func (c *HexConverter) ToBytes() ([]byte, error) {
	b, err := HexToBytes(c.data)
	if err != nil {
		return nil, err
	}
	if c.opts.Size > 0 {
		if len(b) > c.opts.Size {
			return nil, errors.New("size overflow")
		}
		return PadRight(b, c.opts.Size), nil
	}
	return b, nil
}

// ToBigInt converts the hex string to a *big.Int.
func (c *HexConverter) ToBigInt() (*big.Int, error) {
	if c.opts.Size > 0 {
		if err := assertHexSize(c.data, c.opts.Size); err != nil {
			return nil, err
		}
	}
	return HexToBigInt(c.data, c.opts.Signed)
}

// ToNumber converts the hex string to an int64.
func (c *HexConverter) ToNumber() (int64, error) {
	if c.opts.Size > 0 {
		if err := assertHexSize(c.data, c.opts.Size); err != nil {
			return 0, err
		}
	}
	return HexToNumber(c.data, c.opts.Signed)
}

// ToUint converts the hex string to a uint64.
func (c *HexConverter) ToUint() (uint64, error) {
	if c.opts.Size > 0 {
		if err := assertHexSize(c.data, c.opts.Size); err != nil {
			return 0, err
		}
	}
	return HexToUint(c.data)
}

// ToBool converts the hex string to a boolean.
func (c *HexConverter) ToBool() (bool, error) {
	if c.opts.Size > 0 {
		if err := assertHexSize(c.data, c.opts.Size); err != nil {
			return false, err
		}
	}
	return HexToBool(c.data)
}

// ToString converts the hex string to a UTF-8 string.
func (c *HexConverter) ToString() (string, error) {
	b, err := HexToBytes(c.data)
	if err != nil {
		return "", err
	}
	if c.opts.Size > 0 {
		if len(b) > c.opts.Size {
			return "", errors.New("size overflow")
		}
	}
	return BytesToString(TrimRight(b)), nil
}

// String returns the original hex string.
func (c *HexConverter) String() string {
	return c.data
}

// Standalone conversion functions

// HexToBigInt converts a hex string to a *big.Int.
// If signed is true, treats the hex as two's complement signed integer.
func HexToBigInt(s string, signed bool) (*big.Int, error) {
	s = Strip0x(s)
	if len(s) == 0 {
		return big.NewInt(0), nil
	}

	n := new(big.Int)
	_, ok := n.SetString(s, 16)
	if !ok {
		return nil, errors.New("invalid hex string")
	}

	if signed {
		size := (len(s) + 1) / 2 // bytes
		max := new(big.Int).Lsh(big.NewInt(1), uint(size*8-1))
		max.Sub(max, big.NewInt(1))

		if n.Cmp(max) > 0 {
			// Negative number
			fullMax := new(big.Int).Lsh(big.NewInt(1), uint(size*8))
			n.Sub(n, fullMax)
		}
	}

	return n, nil
}

// HexToNumber converts a hex string to an int64.
func HexToNumber(s string, signed bool) (int64, error) {
	bi, err := HexToBigInt(s, signed)
	if err != nil {
		return 0, err
	}
	if !bi.IsInt64() {
		return 0, errors.New("integer out of int64 range")
	}
	return bi.Int64(), nil
}

// HexToUint converts a hex string to a uint64.
func HexToUint(s string) (uint64, error) {
	s = Strip0x(s)
	return strconv.ParseUint(s, 16, 64)
}

// HexToBool converts a hex string to a boolean.
func HexToBool(s string) (bool, error) {
	s = Strip0x(s)
	// Trim leading zeros
	trimmed := ""
	foundNonZero := false
	for _, c := range s {
		if c != '0' {
			foundNonZero = true
		}
		if foundNonZero {
			trimmed += string(c)
		}
	}
	if trimmed == "" || trimmed == "0" {
		return false, nil
	}
	if trimmed == "1" {
		return true, nil
	}
	return false, errors.New("invalid hex boolean")
}

// HexToString converts a hex string to a UTF-8 string.
func HexToString(s string) (string, error) {
	b, err := HexToBytes(s)
	if err != nil {
		return "", err
	}
	return BytesToString(TrimRight(b)), nil
}

// Helper functions

func assertHexSize(s string, size int) error {
	s = Strip0x(s)
	// Pad odd-length
	if len(s)%2 != 0 {
		s = "0" + s
	}
	if len(s)/2 > size {
		return errors.New("size overflow")
	}
	return nil
}

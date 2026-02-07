package unit

import (
	"errors"
	"math/big"
	"strconv"
	"sync"
)

// ErrInvalidDecimalNumber is returned when the value is not a valid decimal number.
var ErrInvalidDecimalNumber = errors.New("invalid decimal number")

// Cached powers of 10 for common decimal values.
// Avoids repeated big.Int allocations for the most common units (6, 8, 9, 18).
var powersOf10 [78]*big.Int

func init() {
	powersOf10[0] = big.NewInt(1)
	for i := 1; i < len(powersOf10); i++ {
		powersOf10[i] = new(big.Int).Mul(powersOf10[i-1], big.NewInt(10))
	}
}

// Pool for reusing big.Int allocations.
var bigIntPool = sync.Pool{
	New: func() any { return new(big.Int) },
}

func getBigInt() *big.Int {
	return bigIntPool.Get().(*big.Int)
}

// Pre-computed zero-padding strings to avoid strings.Repeat allocations.
// Index i contains a string of i zeros. Covers up to 78 decimals.
var zeroPad [79]string

func init() {
	var buf [79]byte
	for i := range buf {
		buf[i] = '0'
	}
	for i := 0; i < len(zeroPad); i++ {
		zeroPad[i] = string(buf[:i])
	}
}

// parseAndSplit does validation + splitting in a single pass over the string.
// Returns integer part, fraction part, whether the value is negative, and whether it's valid.
func parseAndSplit(s string) (integer, fraction string, negative, valid bool) {
	if len(s) == 0 {
		return "", "", false, false
	}

	start := 0
	if s[0] == '-' {
		negative = true
		start = 1
	}

	dot := -1
	hasDigit := false

	for i := start; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if c == '.' && dot < 0 {
			dot = i
		} else {
			return "", "", false, false
		}
	}

	if !hasDigit {
		return "", "", false, false
	}

	if dot < 0 {
		integer = s[start:]
		fraction = ""
	} else {
		integer = s[start:dot]
		fraction = s[dot+1:]
	}

	valid = true
	return
}

// ParseUnits multiplies a string representation of a number by a given exponent
// of base 10 (10^decimals).
//
// Example:
//
//	ParseUnits("420", 9)
//	// big.Int representing 420000000000
//
//	ParseUnits("1", 18)
//	// big.Int representing 1000000000000000000
//
//	ParseUnits("1.5", 18)
//	// big.Int representing 1500000000000000000
func ParseUnits(value string, decimals int) (*big.Int, error) {
	// Single-pass validation + splitting
	integer, fraction, negative, valid := parseAndSplit(value)
	if !valid {
		return nil, ErrInvalidDecimalNumber
	}

	// Trim trailing zeros from fraction
	fraction = trimRight(fraction, '0')

	// Handle rounding when fraction is larger than decimals
	if decimals == 0 {
		if len(fraction) > 0 && fraction[0] >= '5' {
			integer = incrementString(integer)
		}
		fraction = ""
	} else if len(fraction) > decimals {
		left := fraction[:decimals-1]
		unit := fraction[decimals-1]
		right := fraction[decimals:]

		roundVal, _ := strconv.ParseFloat(string(unit)+"."+right, 64)
		rounded := int(roundVal + 0.5)

		if rounded > 9 {
			fraction = incrementString(left) + "0"
			for len(fraction) < len(left)+1 {
				fraction = "0" + fraction
			}
		} else {
			fraction = left + strconv.Itoa(rounded)
		}

		if len(fraction) > decimals {
			fraction = fraction[1:]
			integer = incrementString(integer)
		}

		fraction = fraction[:decimals]
	} else {
		// Pad fraction with trailing zeros — use pre-computed pad to avoid allocation
		if pad := decimals - len(fraction); pad > 0 {
			if pad < len(zeroPad) {
				fraction = fraction + zeroPad[pad]
			} else {
				buf := make([]byte, len(fraction)+pad)
				copy(buf, fraction)
				for i := len(fraction); i < len(buf); i++ {
					buf[i] = '0'
				}
				fraction = string(buf)
			}
		}
	}

	if integer == "" {
		integer = "0"
	}

	intLen := len(integer)
	fracLen := len(fraction)
	combinedLen := intLen + fracLen

	// Fast path: combined fits in uint64 (up to 19 digits)
	if !negative && combinedLen <= 19 {
		val := uint64(0)
		for i := 0; i < intLen; i++ {
			val = val*10 + uint64(integer[i]-'0')
		}
		for i := 0; i < fracLen; i++ {
			val = val*10 + uint64(fraction[i]-'0')
		}
		return new(big.Int).SetUint64(val), nil
	}

	// Arithmetic path: if integer and fraction each fit in uint64, compute
	// result = integer * 10^decimals + fraction using big.Int arithmetic.
	// This avoids big.Int.SetString which is the main bottleneck for large values.
	// Covers the common case: ParseEther("123456789.123456789012345678")
	// where integer="123456789" (9 digits) and fraction has 18 digits.
	if intLen <= 19 && fracLen <= 19 && decimals < len(powersOf10) {
		intVal := uint64(0)
		for i := 0; i < intLen; i++ {
			intVal = intVal*10 + uint64(integer[i]-'0')
		}
		fracVal := uint64(0)
		for i := 0; i < fracLen; i++ {
			fracVal = fracVal*10 + uint64(fraction[i]-'0')
		}

		// result = intVal * 10^decimals + fracVal
		result := getBigInt()
		result.SetUint64(intVal)
		result.Mul(result, powersOf10[decimals])
		result.Add(result, new(big.Int).SetUint64(fracVal))

		if negative {
			result.Neg(result)
		}
		return result, nil
	}

	// Slow path: concatenate and use big.Int.SetString
	if combinedLen <= 96 {
		var buf [96]byte
		copy(buf[:intLen], integer)
		copy(buf[intLen:], fraction)
		result := getBigInt()
		_, ok := result.SetString(string(buf[:combinedLen]), 10)
		if !ok {
			bigIntPool.Put(result)
			return nil, ErrInvalidDecimalNumber
		}
		if negative {
			result.Neg(result)
		}
		return result, nil
	}

	combined := integer + fraction
	result, ok := new(big.Int).SetString(combined, 10)
	if !ok {
		return nil, ErrInvalidDecimalNumber
	}
	if negative {
		result.Neg(result)
	}
	return result, nil
}

// trimRight trims trailing bytes from a string without allocating if nothing to trim.
func trimRight(s string, b byte) string {
	i := len(s)
	for i > 0 && s[i-1] == b {
		i--
	}
	return s[:i]
}

// incrementString increments a decimal number string by 1.
// Avoids big.Int allocation for simple carry operations.
func incrementString(s string) string {
	if s == "" || s == "0" {
		return "1"
	}

	// Fast path: try uint64
	if len(s) <= 18 {
		val, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			return strconv.FormatUint(val+1, 10)
		}
	}

	// Fallback to big.Int for huge numbers
	v := getBigInt()
	_, ok := v.SetString(s, 10)
	if !ok {
		bigIntPool.Put(v)
		return "1"
	}
	v.Add(v, big.NewInt(1))
	result := v.String()
	bigIntPool.Put(v)
	return result
}

// MustParseUnits is like ParseUnits but panics on error.
func MustParseUnits(value string, decimals int) *big.Int {
	result, err := ParseUnits(value, decimals)
	if err != nil {
		panic(err)
	}
	return result
}

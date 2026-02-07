package address

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/sha3"
)

var (
	// ErrInvalidAddress is returned when an address is not valid
	ErrInvalidAddress = errors.New("invalid address")
)

// checksumAddressCache is an LRU cache for checksummed addresses.
// Mirrors viem's checksumAddressCache with 8192 entries.
var checksumAddressCache = newAddressStringCache(8192)

// Reusable hex lookup for lowercase conversion without allocation.
var hexLower [256]byte

func init() {
	for i := range hexLower {
		hexLower[i] = byte(i)
	}
	for c := byte('A'); c <= byte('F'); c++ {
		hexLower[c] = c + 32 // A->a, B->b, etc.
	}
}

// ChecksumAddress converts an address to EIP-55 checksum format.
// Optionally supports EIP-1191 chain-specific checksums (not recommended for general use).
//
// Warning: EIP-1191 checksum addresses are generally not backwards compatible with the
// wider Ethereum ecosystem, meaning it will break when validated against an application/tool
// that relies on EIP-55 checksum encoding (checksum without chainId).
//
// Example:
//
//	checksumAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac")
//	// "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
func ChecksumAddress(address string, chainId ...int64) Address {
	// Check cache
	cacheKey := address
	if len(chainId) > 0 && chainId[0] > 0 {
		cacheKey = address + "." + strconv.FormatInt(chainId[0], 10)
	}
	if cached, ok := checksumAddressCache.Get(cacheKey); ok {
		return Address(cached)
	}

	result := checksumAddressCore(address, chainId...)
	checksumAddressCache.Set(cacheKey, string(result))
	return result
}

// checksumAddressCore is the uncached checksum computation.
func checksumAddressCore(address string, chainId ...int64) Address {
	// Stack-allocated buffer for the lowercase 40-char address
	var addrBuf [40]byte

	// Strip "0x"/"0X" prefix and lowercase in one pass — no allocation
	src := address
	if len(src) >= 2 && (src[0] == '0' && (src[1] == 'x' || src[1] == 'X')) {
		src = src[2:]
	}
	if len(src) != 40 {
		return Address(address) // invalid, return as-is
	}
	for i := 0; i < 40; i++ {
		addrBuf[i] = hexLower[src[i]]
	}

	// Keccak256 hash — reuse stack buffer for output
	var hashBuf [32]byte
	hasher := sha3.NewLegacyKeccak256()

	if len(chainId) > 0 && chainId[0] > 0 {
		// EIP-1191: include chain ID prefix
		prefix := strconv.AppendInt(nil, chainId[0], 10)
		hasher.Write(prefix)
		hasher.Write([]byte("0x"))
	}
	hasher.Write(addrBuf[:])
	hasher.Sum(hashBuf[:0]) // write into stack buffer, no alloc

	// Apply checksum into a stack-allocated result buffer
	var result [42]byte
	result[0] = '0'
	result[1] = 'x'
	for i := 0; i < 40; i++ {
		c := addrBuf[i]

		var nibble byte
		if i%2 == 0 {
			nibble = hashBuf[i/2] >> 4
		} else {
			nibble = hashBuf[i/2] & 0x0f
		}

		if nibble >= 8 && c >= 'a' && c <= 'f' {
			result[i+2] = c - 32
		} else {
			result[i+2] = c
		}
	}

	return Address(result[:])
}

// GetAddress validates an address and returns it in checksummed format.
// Returns an error if the address is invalid.
//
// Example:
//
//	getAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac")
//	// "0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"
func GetAddress(address string, chainId ...int64) (Address, error) {
	if !IsAddress(address, IsAddressOptions{Strict: false}) {
		return "", fmt.Errorf("%w: %s", ErrInvalidAddress, address)
	}

	if len(chainId) > 0 {
		return ChecksumAddress(address, chainId[0]), nil
	}
	return ChecksumAddress(address), nil
}

// keccak256 computes the Keccak-256 hash of input data.
func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

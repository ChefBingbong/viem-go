package hash

import (
	"regexp"
)

var hashRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{64}$`)

// IsHash checks if a string is a valid 32-byte hash (64 hex characters with 0x prefix).
//
// Example:
//
//	isHash("0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad") // true
//	isHash("0x1234") // false (too short)
func IsHash(hash string) bool {
	return hashRegex.MatchString(hash)
}

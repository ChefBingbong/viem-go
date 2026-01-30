package hash

import (
	"errors"
	"strings"
)

// ErrUnableToNormalizeSignature is returned when signature normalization fails.
var ErrUnableToNormalizeSignature = errors.New("unable to normalize signature")

// NormalizeSignature normalizes a function or event signature by removing
// parameter names and unnecessary whitespace, keeping only the types.
//
// Example:
//
//	sig, _ := NormalizeSignature("function transfer(address to, uint256 amount)")
//	// "transfer(address,uint256)"
//
//	sig, _ := NormalizeSignature("event Transfer(address indexed from, address indexed to, uint256 amount)")
//	// "Transfer(address,address,uint256)"
func NormalizeSignature(signature string) (string, error) {
	active := true
	current := ""
	level := 0
	result := ""
	valid := false

	for i := 0; i < len(signature); i++ {
		char := signature[i]

		// If the character is a separator, we want to reactivate.
		if char == '(' || char == ')' || char == ',' {
			active = true
		}

		// If the character is a "level" token, we want to increment/decrement.
		if char == '(' {
			level++
		}
		if char == ')' {
			level--
		}

		// If we aren't active, we don't want to mutate the result.
		if !active {
			continue
		}

		// If level === 0, we are at the definition level.
		if level == 0 {
			if char == ' ' && (result == "event" || result == "function" || result == "") {
				result = ""
			} else {
				result += string(char)

				// If we are at the end of the definition, we must be finished.
				if char == ')' {
					valid = true
					break
				}
			}

			continue
		}

		// Ignore spaces
		if char == ' ' {
			// If the previous character is a separator, and the current section isn't empty, we want to deactivate.
			if i > 0 && signature[i-1] != ',' && current != "," && !strings.HasSuffix(current, ",(") {
				current = ""
				active = false
			}
			continue
		}

		result += string(char)
		current += string(char)
	}

	if !valid {
		return "", ErrUnableToNormalizeSignature
	}

	return result, nil
}

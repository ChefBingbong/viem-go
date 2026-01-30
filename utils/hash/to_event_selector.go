package hash

// ToEventSelector returns the event selector (topic0) for a given event definition.
// The selector is the full keccak256 hash of the event signature.
//
// Example:
//
//	selector, _ := ToEventSelector("event Transfer(address indexed from, address indexed to, uint256 amount)")
//	// "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
func ToEventSelector(event string) (string, error) {
	return ToSignatureHash(event)
}

// ToEventSelectorBytes returns the selector as raw bytes.
func ToEventSelectorBytes(event string) ([]byte, error) {
	return ToSignatureHashBytes(event)
}

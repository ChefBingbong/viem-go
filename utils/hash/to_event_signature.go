package hash

// ToEventSignature is an alias for ToSignature, specifically for event definitions.
//
// Example:
//
//	sig, _ := ToEventSignature("event Transfer(address indexed from, address indexed to, uint256 amount)")
//	// "Transfer(address,address,uint256)"
func ToEventSignature(event string) (string, error) {
	return ToSignature(event)
}

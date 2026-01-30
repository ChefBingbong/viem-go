package hash

// ToFunctionSignature is an alias for ToSignature, specifically for function definitions.
//
// Example:
//
//	sig, _ := ToFunctionSignature("function transfer(address to, uint256 amount)")
//	// "transfer(address,uint256)"
func ToFunctionSignature(fn string) (string, error) {
	return ToSignature(fn)
}

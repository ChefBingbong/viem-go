package signature

// RecoverMessageAddress recovers the Ethereum address that signed a message.
//
// Example:
//
//	address, err := RecoverMessageAddress(
//		NewSignableMessage("hello world"),
//		"0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c",
//	)
//	// "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
func RecoverMessageAddress(message SignableMessage, signature any) (string, error) {
	hash := HashMessage(message)
	return RecoverAddress(hash, signature)
}

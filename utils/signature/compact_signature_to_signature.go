package signature

// CompactSignatureToSignature converts an EIP-2098 compact signature to a standard signature.
// https://eips.ethereum.org/EIPS/eip-2098
//
// Example:
//
//	sig, err := CompactSignatureToSignature(&CompactSignature{
//		R: "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90",
//		YParityAndS: "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064",
//	})
//	// sig.R = "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90"
//	// sig.S = "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064"
//	// sig.YParity = 0
func CompactSignatureToSignature(compact *CompactSignature) (*Signature, error) {
	if compact == nil {
		return nil, ErrInvalidSignatureLength
	}

	// Convert yParityAndS to bytes
	yParityAndSBytes := hexToBytes(compact.YParityAndS)
	if len(yParityAndSBytes) == 0 {
		yParityAndSBytes = make([]byte, 32)
	}

	// Pad to 32 bytes if needed
	if len(yParityAndSBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(yParityAndSBytes):], yParityAndSBytes)
		yParityAndSBytes = padded
	}

	// Extract yParity from the top bit of the first byte
	yParity := 0
	if yParityAndSBytes[0]&0x80 != 0 {
		yParity = 1
	}

	// Clear the top bit to get s
	sBytes := make([]byte, 32)
	copy(sBytes, yParityAndSBytes)
	if yParity == 1 {
		sBytes[0] &= 0x7f
	}

	return &Signature{
		R:       compact.R,
		S:       bytesToHex(sBytes),
		YParity: yParity,
	}, nil
}

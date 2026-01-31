package formatters

// FormatProof formats an RPC proof into a Proof struct.
//
// Example:
//
//	rpcProof := RpcProof{
//		Address:      "0x...",
//		Balance:      "0xde0b6b3a7640000",
//		Nonce:        "0x1",
//		StorageProof: []RpcStorageProof{...},
//	}
//	proof := FormatProof(rpcProof)
func FormatProof(proof RpcProof) Proof {
	result := Proof{
		Address:      proof.Address,
		AccountProof: proof.AccountProof,
		CodeHash:     proof.CodeHash,
		StorageHash:  proof.StorageHash,
	}

	// Balance
	if proof.Balance != "" {
		result.Balance = hexToBigInt(proof.Balance)
	}

	// Nonce
	if proof.Nonce != "" {
		nonce := hexToInt(proof.Nonce)
		result.Nonce = &nonce
	}

	// Storage proof
	if len(proof.StorageProof) > 0 {
		result.StorageProof = formatStorageProof(proof.StorageProof)
	}

	return result
}

// formatStorageProof formats storage proofs.
func formatStorageProof(storageProof []RpcStorageProof) []StorageProof {
	result := make([]StorageProof, len(storageProof))
	for i, sp := range storageProof {
		result[i] = StorageProof{
			Key:   sp.Key,
			Proof: sp.Proof,
			Value: hexToBigInt(sp.Value),
		}
	}
	return result
}

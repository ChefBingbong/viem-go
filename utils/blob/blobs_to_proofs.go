package blob

import (
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// BlobsToProofs computes KZG proofs for a list of blobs and their commitments.
//
// Example:
//
//	blobs, _ := ToBlobs(data)
//	commitments, _ := BlobsToCommitments(blobs, kzgImpl)
//	proofs, err := BlobsToProofs(blobs, commitments, kzgImpl)
func BlobsToProofs(blobs [][]byte, commitments [][]byte, kzgImpl kzg.Kzg) ([][]byte, error) {
	if len(blobs) != len(commitments) {
		return nil, kzg.ErrInvalidBlobSize
	}

	proofs := make([][]byte, len(blobs))

	for i, blob := range blobs {
		proof, err := kzgImpl.ComputeBlobKzgProof(blob, commitments[i])
		if err != nil {
			return nil, err
		}
		proofs[i] = proof
	}

	return proofs, nil
}

// BlobsToProofsHex computes KZG proofs and returns them as hex strings.
func BlobsToProofsHex(blobs [][]byte, commitments [][]byte, kzgImpl kzg.Kzg) ([]string, error) {
	proofs, err := BlobsToProofs(blobs, commitments, kzgImpl)
	if err != nil {
		return nil, err
	}

	hexProofs := make([]string, len(proofs))
	for i, proof := range proofs {
		hexProofs[i] = bytesToHex(proof)
	}

	return hexProofs, nil
}

// BlobsHexToProofs computes proofs from hex-encoded blobs and commitments.
func BlobsHexToProofs(hexBlobs []string, hexCommitments []string, kzgImpl kzg.Kzg) ([][]byte, error) {
	blobs := make([][]byte, len(hexBlobs))
	for i, hexBlob := range hexBlobs {
		blob, err := hexToBytes(hexBlob)
		if err != nil {
			return nil, err
		}
		blobs[i] = blob
	}

	commitments := make([][]byte, len(hexCommitments))
	for i, hexCommitment := range hexCommitments {
		commitment, err := hexToBytes(hexCommitment)
		if err != nil {
			return nil, err
		}
		commitments[i] = commitment
	}

	return BlobsToProofs(blobs, commitments, kzgImpl)
}

// BlobsHexToProofsHex computes hex proofs from hex blobs and commitments.
func BlobsHexToProofsHex(hexBlobs []string, hexCommitments []string, kzgImpl kzg.Kzg) ([]string, error) {
	proofs, err := BlobsHexToProofs(hexBlobs, hexCommitments, kzgImpl)
	if err != nil {
		return nil, err
	}

	hexProofs := make([]string, len(proofs))
	for i, proof := range proofs {
		hexProofs[i] = bytesToHex(proof)
	}

	return hexProofs, nil
}

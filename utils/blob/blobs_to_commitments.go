package blob

import (
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// BlobsToCommitments computes KZG commitments for a list of blobs.
//
// Example:
//
//	blobs, _ := ToBlobs(data)
//	commitments, err := BlobsToCommitments(blobs, kzgImpl)
func BlobsToCommitments(blobs [][]byte, kzgImpl kzg.Kzg) ([][]byte, error) {
	commitments := make([][]byte, len(blobs))

	for i, blob := range blobs {
		commitment, err := kzgImpl.BlobToKzgCommitment(blob)
		if err != nil {
			return nil, err
		}
		commitments[i] = commitment
	}

	return commitments, nil
}

// BlobsToCommitmentsHex computes KZG commitments and returns them as hex strings.
func BlobsToCommitmentsHex(blobs [][]byte, kzgImpl kzg.Kzg) ([]string, error) {
	commitments, err := BlobsToCommitments(blobs, kzgImpl)
	if err != nil {
		return nil, err
	}

	hexCommitments := make([]string, len(commitments))
	for i, commitment := range commitments {
		hexCommitments[i] = bytesToHex(commitment)
	}

	return hexCommitments, nil
}

// BlobsHexToCommitments computes KZG commitments from hex-encoded blobs.
func BlobsHexToCommitments(hexBlobs []string, kzgImpl kzg.Kzg) ([][]byte, error) {
	blobs := make([][]byte, len(hexBlobs))
	for i, hexBlob := range hexBlobs {
		blob, err := hexToBytes(hexBlob)
		if err != nil {
			return nil, err
		}
		blobs[i] = blob
	}

	return BlobsToCommitments(blobs, kzgImpl)
}

// BlobsHexToCommitmentsHex computes hex commitments from hex blobs.
func BlobsHexToCommitmentsHex(hexBlobs []string, kzgImpl kzg.Kzg) ([]string, error) {
	blobs := make([][]byte, len(hexBlobs))
	for i, hexBlob := range hexBlobs {
		blob, err := hexToBytes(hexBlob)
		if err != nil {
			return nil, err
		}
		blobs[i] = blob
	}

	return BlobsToCommitmentsHex(blobs, kzgImpl)
}

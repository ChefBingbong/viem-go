package blob

import (
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// ToBlobSidecarsParams contains parameters for creating blob sidecars.
type ToBlobSidecarsParams struct {
	// Data to transform into blobs (mutually exclusive with Blobs)
	Data []byte
	// Pre-computed blobs (mutually exclusive with Data)
	Blobs [][]byte
	// Pre-computed commitments (required if Blobs is set without Kzg)
	Commitments [][]byte
	// Pre-computed proofs (required if Blobs is set without Kzg)
	Proofs [][]byte
	// KZG implementation (required if Data is set or if Commitments/Proofs need to be computed)
	Kzg kzg.Kzg
}

// ToBlobSidecars creates blob sidecars from data or pre-computed components.
//
// Example with raw data:
//
//	sidecars, err := ToBlobSidecars(ToBlobSidecarsParams{
//		Data: []byte("hello world"),
//		Kzg:  kzgImpl,
//	})
//
// Example with pre-computed components:
//
//	sidecars, err := ToBlobSidecars(ToBlobSidecarsParams{
//		Blobs:       blobs,
//		Commitments: commitments,
//		Proofs:      proofs,
//	})
func ToBlobSidecars(params ToBlobSidecarsParams) ([]BlobSidecar, error) {
	var blobs [][]byte
	var commitments [][]byte
	var proofs [][]byte
	var err error

	// Get or create blobs
	if params.Blobs != nil {
		blobs = params.Blobs
	} else if params.Data != nil {
		blobs, err = ToBlobs(params.Data)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, kzg.ErrEmptyBlob
	}

	// Get or compute commitments
	if params.Commitments != nil {
		commitments = params.Commitments
	} else if params.Kzg != nil {
		commitments, err = BlobsToCommitments(blobs, params.Kzg)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, kzg.ErrKzgNotInitialized
	}

	// Get or compute proofs
	if params.Proofs != nil {
		proofs = params.Proofs
	} else if params.Kzg != nil {
		proofs, err = BlobsToProofs(blobs, commitments, params.Kzg)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, kzg.ErrKzgNotInitialized
	}

	// Create sidecars
	sidecars := make([]BlobSidecar, len(blobs))
	for i := range blobs {
		sidecars[i] = BlobSidecar{
			Blob:       blobs[i],
			Commitment: commitments[i],
			Proof:      proofs[i],
		}
	}

	return sidecars, nil
}

// ToBlobSidecarsHex creates hex-encoded blob sidecars.
func ToBlobSidecarsHex(params ToBlobSidecarsParams) ([]BlobSidecarHex, error) {
	sidecars, err := ToBlobSidecars(params)
	if err != nil {
		return nil, err
	}

	hexSidecars := make([]BlobSidecarHex, len(sidecars))
	for i, sidecar := range sidecars {
		hexSidecars[i] = BlobSidecarHex{
			Blob:       bytesToHex(sidecar.Blob),
			Commitment: bytesToHex(sidecar.Commitment),
			Proof:      bytesToHex(sidecar.Proof),
		}
	}

	return hexSidecars, nil
}

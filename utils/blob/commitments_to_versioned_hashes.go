package blob

import (
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// CommitmentsToVersionedHashes transforms a list of commitments to their versioned hashes.
//
// Example:
//
//	commitments, _ := BlobsToCommitments(blobs, kzgImpl)
//	versionedHashes := CommitmentsToVersionedHashes(commitments, kzg.VersionedHashVersionKzg)
func CommitmentsToVersionedHashes(commitments [][]byte, version byte) [][]byte {
	hashes := make([][]byte, len(commitments))
	for i, commitment := range commitments {
		hashes[i] = CommitmentToVersionedHash(commitment, version)
	}
	return hashes
}

// CommitmentsToVersionedHashesDefault uses the default KZG version (0x01).
func CommitmentsToVersionedHashesDefault(commitments [][]byte) [][]byte {
	return CommitmentsToVersionedHashes(commitments, kzg.VersionedHashVersionKzg)
}

// CommitmentsToVersionedHashesHex returns the versioned hashes as hex strings.
func CommitmentsToVersionedHashesHex(commitments [][]byte, version byte) []string {
	hashes := make([]string, len(commitments))
	for i, commitment := range commitments {
		hashes[i] = CommitmentToVersionedHashHex(commitment, version)
	}
	return hashes
}

// CommitmentsHexToVersionedHashes computes versioned hashes from hex commitments.
func CommitmentsHexToVersionedHashes(hexCommitments []string, version byte) ([][]byte, error) {
	hashes := make([][]byte, len(hexCommitments))
	for i, hexCommitment := range hexCommitments {
		hash, err := CommitmentHexToVersionedHash(hexCommitment, version)
		if err != nil {
			return nil, err
		}
		hashes[i] = hash
	}
	return hashes, nil
}

// CommitmentsHexToVersionedHashesHex returns hex hashes from hex commitments.
func CommitmentsHexToVersionedHashesHex(hexCommitments []string, version byte) ([]string, error) {
	hashes := make([]string, len(hexCommitments))
	for i, hexCommitment := range hexCommitments {
		hash, err := CommitmentHexToVersionedHashHex(hexCommitment, version)
		if err != nil {
			return nil, err
		}
		hashes[i] = hash
	}
	return hashes, nil
}

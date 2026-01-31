package blob

import (
	"crypto/sha256"

	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// CommitmentToVersionedHash transforms a KZG commitment to its versioned hash.
// The versioned hash is SHA256(commitment) with the first byte replaced by the version.
//
// Example:
//
//	commitment, _ := kzgImpl.BlobToKzgCommitment(blob)
//	versionedHash := CommitmentToVersionedHash(commitment, kzg.VersionedHashVersionKzg)
func CommitmentToVersionedHash(commitment []byte, version byte) []byte {
	hash := sha256.Sum256(commitment)
	hash[0] = version
	return hash[:]
}

// CommitmentToVersionedHashDefault uses the default KZG version (0x01).
func CommitmentToVersionedHashDefault(commitment []byte) []byte {
	return CommitmentToVersionedHash(commitment, kzg.VersionedHashVersionKzg)
}

// CommitmentToVersionedHashHex returns the versioned hash as a hex string.
func CommitmentToVersionedHashHex(commitment []byte, version byte) string {
	return bytesToHex(CommitmentToVersionedHash(commitment, version))
}

// CommitmentHexToVersionedHash computes versioned hash from hex commitment.
func CommitmentHexToVersionedHash(hexCommitment string, version byte) ([]byte, error) {
	commitment, err := hexToBytes(hexCommitment)
	if err != nil {
		return nil, err
	}
	return CommitmentToVersionedHash(commitment, version), nil
}

// CommitmentHexToVersionedHashHex returns hex versioned hash from hex commitment.
func CommitmentHexToVersionedHashHex(hexCommitment string, version byte) (string, error) {
	hash, err := CommitmentHexToVersionedHash(hexCommitment, version)
	if err != nil {
		return "", err
	}
	return bytesToHex(hash), nil
}

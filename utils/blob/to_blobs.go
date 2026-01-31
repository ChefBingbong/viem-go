package blob

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// ToBlobs transforms arbitrary data into blobs for EIP-4844 transactions.
// Each blob is exactly BytesPerBlob (131072) bytes.
//
// The encoding follows EIP-4844 requirements:
// - Each field element (32 bytes) starts with a 0x00 byte to prevent overflow
// - The actual data is stored in the remaining 31 bytes
// - A terminator byte (0x80) marks the end of the data
//
// Example:
//
//	blobs, err := ToBlobs([]byte("hello world"))
//	// Returns one or more 131072-byte blobs
func ToBlobs(data []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, kzg.ErrEmptyBlob
	}

	if len(data) > kzg.MaxBytesPerTransaction {
		return nil, fmt.Errorf("%w: %d bytes exceeds max %d bytes",
			kzg.ErrBlobSizeTooLarge, len(data), kzg.MaxBytesPerTransaction)
	}

	var blobs [][]byte
	position := 0
	active := true

	for active {
		blob := make([]byte, BytesPerBlob)
		blobPos := 0
		fieldElements := 0

		for fieldElements < FieldElementsPerBlob {
			// Push a zero byte so the field element doesn't overflow the BLS modulus
			blob[blobPos] = 0x00
			blobPos++

			// Calculate how many bytes to copy (max 31 per field element)
			remaining := len(data) - position
			copyLen := 31
			if remaining < 31 {
				copyLen = remaining
			}

			// Copy data bytes
			if copyLen > 0 {
				copy(blob[blobPos:], data[position:position+copyLen])
				blobPos += copyLen
				position += copyLen
			}

			// If we've consumed all data or copied less than 31 bytes, add terminator
			if copyLen < 31 {
				blob[blobPos] = 0x80
				active = false
				break
			}

			fieldElements++
		}

		blobs = append(blobs, blob)
	}

	return blobs, nil
}

// ToBlobsHex transforms arbitrary data into hex-encoded blobs.
func ToBlobsHex(data []byte) ([]string, error) {
	blobs, err := ToBlobs(data)
	if err != nil {
		return nil, err
	}

	hexBlobs := make([]string, len(blobs))
	for i, blob := range blobs {
		hexBlobs[i] = "0x" + hex.EncodeToString(blob)
	}

	return hexBlobs, nil
}

// ToBlobsFromHex transforms hex-encoded data into blobs.
func ToBlobsFromHex(hexData string) ([][]byte, error) {
	data, err := hexToBytes(hexData)
	if err != nil {
		return nil, err
	}
	return ToBlobs(data)
}

// hexToBytes converts a hex string to bytes.
func hexToBytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// bytesToHex converts bytes to a hex string.
func bytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

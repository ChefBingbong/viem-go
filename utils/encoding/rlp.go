package encoding

import (
	"errors"
)

// RecursiveBytes represents a recursive array of bytes for RLP encoding.
type RecursiveBytes []any // can contain []byte or RecursiveBytes

// RlpEncoder provides RLP encoding functionality.
type RlpEncoder struct {
	data any // []byte, string (hex), or RecursiveBytes
}

// ToRlp creates a new RLP encoder.
func ToRlp(data any) *RlpEncoder {
	return &RlpEncoder{data: data}
}

// Bytes encodes to RLP and returns bytes.
func (e *RlpEncoder) Bytes() ([]byte, error) {
	return rlpEncode(e.data)
}

// Hex encodes to RLP and returns hex string.
func (e *RlpEncoder) Hex() (string, error) {
	b, err := e.Bytes()
	if err != nil {
		return "", err
	}
	return BytesToHex(b), nil
}

// RlpDecoder provides RLP decoding functionality.
type RlpDecoder struct {
	data []byte
}

// FromRlp creates a new RLP decoder.
func FromRlp(data any) (*RlpDecoder, error) {
	var b []byte
	switch v := data.(type) {
	case []byte:
		b = v
	case string:
		var err error
		b, err = HexToBytes(v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid RLP input type")
	}
	return &RlpDecoder{data: b}, nil
}

// Bytes decodes RLP and returns the result as recursive bytes.
func (d *RlpDecoder) Bytes() (any, error) {
	if len(d.data) == 0 {
		return []byte{}, nil
	}
	result, _, err := rlpDecode(d.data, 0)
	return result, err
}

// Hex decodes RLP and returns the result as recursive hex strings.
func (d *RlpDecoder) Hex() (any, error) {
	result, err := d.Bytes()
	if err != nil {
		return nil, err
	}
	return convertToHex(result), nil
}

// Standalone functions

// RlpEncode encodes data to RLP bytes.
func RlpEncode(data any) ([]byte, error) {
	return rlpEncode(data)
}

// RlpEncodeToHex encodes data to RLP hex string.
func RlpEncodeToHex(data any) (string, error) {
	b, err := RlpEncode(data)
	if err != nil {
		return "", err
	}
	return BytesToHex(b), nil
}

// RlpDecode decodes RLP bytes.
func RlpDecode(data []byte) (any, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}
	result, _, err := rlpDecode(data, 0)
	return result, err
}

// RlpDecodeHex decodes RLP hex string.
func RlpDecodeHex(s string) (any, error) {
	b, err := HexToBytes(s)
	if err != nil {
		return nil, err
	}
	return RlpDecode(b)
}

// Internal RLP encoding

func rlpEncode(data any) ([]byte, error) {
	switch v := data.(type) {
	case []byte:
		return rlpEncodeBytes(v), nil
	case string:
		// Assume hex string
		b, err := HexToBytes(v)
		if err != nil {
			return nil, err
		}
		return rlpEncodeBytes(b), nil
	case []any:
		return rlpEncodeList(v)
	case RecursiveBytes:
		return rlpEncodeList(v)
	default:
		return nil, errors.New("unsupported RLP type")
	}
}

func rlpEncodeBytes(b []byte) []byte {
	if len(b) == 1 && b[0] < 0x80 {
		// Single byte < 0x80 encodes as itself
		return b
	}
	if len(b) <= 55 {
		// Short string: 0x80 + len, then bytes
		return append([]byte{byte(0x80 + len(b))}, b...)
	}
	// Long string: 0xb7 + len of len, len, then bytes
	lenBytes := encodeLength(len(b))
	prefix := byte(0xb7 + len(lenBytes))
	return append(append([]byte{prefix}, lenBytes...), b...)
}

func rlpEncodeList(items []any) ([]byte, error) {
	var encoded []byte
	for _, item := range items {
		b, err := rlpEncode(item)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, b...)
	}

	if len(encoded) <= 55 {
		// Short list: 0xc0 + len, then items
		return append([]byte{byte(0xc0 + len(encoded))}, encoded...), nil
	}
	// Long list: 0xf7 + len of len, len, then items
	lenBytes := encodeLength(len(encoded))
	prefix := byte(0xf7 + len(lenBytes))
	return append(append([]byte{prefix}, lenBytes...), encoded...), nil
}

func encodeLength(length int) []byte {
	if length < 256 {
		return []byte{byte(length)}
	}
	if length < 65536 {
		return []byte{byte(length >> 8), byte(length)}
	}
	if length < 16777216 {
		return []byte{byte(length >> 16), byte(length >> 8), byte(length)}
	}
	return []byte{byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)}
}

// Internal RLP decoding

func rlpDecode(data []byte, offset int) (any, int, error) {
	if offset >= len(data) {
		return nil, offset, errors.New("unexpected end of RLP data")
	}

	prefix := data[offset]

	// Single byte
	if prefix < 0x80 {
		return []byte{prefix}, offset + 1, nil
	}

	// Short string (0-55 bytes)
	if prefix <= 0xb7 {
		length := int(prefix - 0x80)
		if offset+1+length > len(data) {
			return nil, 0, errors.New("invalid RLP string length")
		}
		return data[offset+1 : offset+1+length], offset + 1 + length, nil
	}

	// Long string
	if prefix <= 0xbf {
		lenOfLen := int(prefix - 0xb7)
		if offset+1+lenOfLen > len(data) {
			return nil, 0, errors.New("invalid RLP string length")
		}
		length := decodeLength(data[offset+1 : offset+1+lenOfLen])
		if offset+1+lenOfLen+length > len(data) {
			return nil, 0, errors.New("invalid RLP string length")
		}
		return data[offset+1+lenOfLen : offset+1+lenOfLen+length], offset + 1 + lenOfLen + length, nil
	}

	// Short list (0-55 bytes total)
	if prefix <= 0xf7 {
		length := int(prefix - 0xc0)
		return decodeList(data, offset+1, length)
	}

	// Long list
	lenOfLen := int(prefix - 0xf7)
	if offset+1+lenOfLen > len(data) {
		return nil, 0, errors.New("invalid RLP list length")
	}
	length := decodeLength(data[offset+1 : offset+1+lenOfLen])
	return decodeList(data, offset+1+lenOfLen, length)
}

func decodeLength(data []byte) int {
	var length int
	for _, b := range data {
		length = (length << 8) | int(b)
	}
	return length
}

func decodeList(data []byte, start, length int) ([]any, int, error) {
	var result []any
	end := start + length
	offset := start

	for offset < end {
		item, newOffset, err := rlpDecode(data, offset)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, item)
		offset = newOffset
	}

	return result, end, nil
}

func convertToHex(data any) any {
	switch v := data.(type) {
	case []byte:
		return BytesToHex(v)
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = convertToHex(item)
		}
		return result
	default:
		return v
	}
}

package encoding_test

import (
	"math/big"
	"testing"

	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// FromBytes tests

func TestBytesToHex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"empty", []byte{}, "0x"},
		{"hello world", []byte("Hello World!"), "0x48656c6c6f20576f726c6421"},
		{"single byte", []byte{0x01}, "0x01"},
		{"multiple bytes", []byte{0x01, 0xa4}, "0x01a4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encoding.BytesToHex(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToHex(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytesToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		signed   bool
		expected int64
	}{
		{"420 unsigned", []byte{0x01, 0xa4}, false, 420},
		{"1 unsigned", []byte{0x01}, false, 1},
		{"0 unsigned", []byte{0x00}, false, 0},
		{"positive signed", []byte{0x01, 0xa4}, true, 420},
		{"negative signed", []byte{0xff}, true, -1},
		{"negative signed 2", []byte{0x80}, true, -128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encoding.BytesToBigInt(tt.input, tt.signed)
			if result.Int64() != tt.expected {
				t.Errorf("BytesToBigInt(%v, %v) = %d, want %d", tt.input, tt.signed, result.Int64(), tt.expected)
			}
		})
	}
}

func TestBytesToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
		hasError bool
	}{
		{"true", []byte{0x01}, true, false},
		{"false", []byte{0x00}, false, false},
		{"empty", []byte{}, false, false},
		{"padded true", []byte{0x00, 0x00, 0x01}, true, false},
		{"invalid", []byte{0x02}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.BytesToBool(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("BytesToBool(%v) expected error", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("BytesToBool(%v) unexpected error: %v", tt.input, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("BytesToBool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"hello world", []byte("Hello World!"), "Hello World!"},
		{"with null padding", append([]byte("Hello"), 0x00, 0x00), "Hello"},
		{"empty", []byte{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encoding.BytesToString(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToString(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFromBytesBuilder(t *testing.T) {
	t.Run("to hex", func(t *testing.T) {
		result, err := encoding.FromBytes([]byte{0x01, 0xa4}).ToHex()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "0x01a4" {
			t.Errorf("got %s, want 0x01a4", result)
		}
	})

	t.Run("to bigint signed", func(t *testing.T) {
		result, err := encoding.FromBytes([]byte{0xff}).WithSigned().ToBigInt()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Int64() != -1 {
			t.Errorf("got %d, want -1", result.Int64())
		}
	})

	t.Run("with size validation", func(t *testing.T) {
		_, err := encoding.FromBytes([]byte{0x01, 0x02, 0x03}).WithSize(2).ToHex()
		if err == nil {
			t.Error("expected size overflow error")
		}
	})
}

// ToBytes tests

func TestHexToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{"with prefix", "0x48656c6c6f20576f726c6421", []byte("Hello World!")},
		{"without prefix", "48656c6c6f", []byte("Hello")},
		{"odd length", "0x1a4", []byte{0x01, 0xa4}},
		{"empty", "0x", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.HexToBytes(tt.input)
			if err != nil {
				t.Errorf("HexToBytes(%s) error: %v", tt.input, err)
			}
			if string(result) != string(tt.expected) {
				t.Errorf("HexToBytes(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	result := encoding.StringToBytes("Hello World!")
	expected := []byte("Hello World!")
	if string(result) != string(expected) {
		t.Errorf("StringToBytes = %v, want %v", result, expected)
	}
}

func TestBoolToBytes(t *testing.T) {
	if encoding.BoolToBytes(true)[0] != 0x01 {
		t.Error("BoolToBytes(true) should be [0x01]")
	}
	if encoding.BoolToBytes(false)[0] != 0x00 {
		t.Error("BoolToBytes(false) should be [0x00]")
	}
}

func TestNumberToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected []byte
	}{
		{"420", big.NewInt(420), []byte{0x01, 0xa4}},
		{"0", big.NewInt(0), []byte{0x00}},
		{"1", big.NewInt(1), []byte{0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encoding.NumberToBytes(tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("NumberToBytes(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToBytesBuilder(t *testing.T) {
	t.Run("from number with size", func(t *testing.T) {
		result, err := encoding.ToBytes(big.NewInt(420)).WithSize(4).Bytes()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := []byte{0x00, 0x00, 0x01, 0xa4}
		if string(result) != string(expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})

	t.Run("from hex string", func(t *testing.T) {
		result, err := encoding.ToBytes("0x1a4").Bytes()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := []byte{0x01, 0xa4}
		if string(result) != string(expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})

	t.Run("from bool", func(t *testing.T) {
		result, err := encoding.ToBytes(true).WithSize(32).Bytes()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 32 {
			t.Errorf("expected length 32, got %d", len(result))
		}
		if result[31] != 0x01 {
			t.Error("expected last byte to be 0x01")
		}
	})

	t.Run("signed negative", func(t *testing.T) {
		result, err := encoding.ToBytes(big.NewInt(-1)).WithSize(1).WithSigned().Bytes()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result[0] != 0xff {
			t.Errorf("got %v, want [0xff]", result)
		}
	})
}

// FromHex tests

func TestHexToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		signed   bool
		expected int64
	}{
		{"420 unsigned", "0x1a4", false, 420},
		{"1 unsigned", "0x1", false, 1},
		{"0 unsigned", "0x0", false, 0},
		{"positive signed", "0x1a4", true, 420},
		{"negative signed", "0xff", true, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.HexToBigInt(tt.input, tt.signed)
			if err != nil {
				t.Errorf("HexToBigInt(%s, %v) error: %v", tt.input, tt.signed, err)
			}
			if result.Int64() != tt.expected {
				t.Errorf("HexToBigInt(%s, %v) = %d, want %d", tt.input, tt.signed, result.Int64(), tt.expected)
			}
		})
	}
}

func TestHexToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"true", "0x1", true, false},
		{"false", "0x0", false, false},
		{"padded true", "0x0000000000000000000000000000000000000000000000000000000000000001", true, false},
		{"padded false", "0x0000000000000000000000000000000000000000000000000000000000000000", false, false},
		{"invalid", "0x2", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.HexToBool(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("HexToBool(%s) expected error", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("HexToBool(%s) unexpected error: %v", tt.input, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("HexToBool(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHexToString(t *testing.T) {
	result, err := encoding.HexToString("0x48656c6c6f20576f726c6421")
	if err != nil {
		t.Errorf("HexToString error: %v", err)
	}
	if result != "Hello World!" {
		t.Errorf("HexToString = %s, want Hello World!", result)
	}
}

func TestFromHexBuilder(t *testing.T) {
	t.Run("to bigint signed", func(t *testing.T) {
		result, err := encoding.FromHex("0xff").WithSigned().ToBigInt()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Int64() != -1 {
			t.Errorf("got %d, want -1", result.Int64())
		}
	})

	t.Run("to bytes", func(t *testing.T) {
		result, err := encoding.FromHex("0x1a4").ToBytes()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 2 || result[0] != 0x01 || result[1] != 0xa4 {
			t.Errorf("got %v, want [0x01, 0xa4]", result)
		}
	})
}

// ToHex tests

func TestBoolToHex(t *testing.T) {
	if encoding.BoolToHex(true) != "0x1" {
		t.Error("BoolToHex(true) should be 0x1")
	}
	if encoding.BoolToHex(false) != "0x0" {
		t.Error("BoolToHex(false) should be 0x0")
	}
}

func TestNumberToHex(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{"420", big.NewInt(420), "0x1a4"},
		{"0", big.NewInt(0), "0x0"},
		{"1", big.NewInt(1), "0x1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encoding.NumberToHex(tt.input)
			if result != tt.expected {
				t.Errorf("NumberToHex(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStringToHex(t *testing.T) {
	result := encoding.StringToHex("Hello World!")
	if result != "0x48656c6c6f20576f726c6421" {
		t.Errorf("StringToHex = %s, want 0x48656c6c6f20576f726c6421", result)
	}
}

func TestToHexBuilder(t *testing.T) {
	t.Run("number with size", func(t *testing.T) {
		result, err := encoding.ToHex(big.NewInt(420)).WithSize(32).Hex()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x00000000000000000000000000000000000000000000000000000000000001a4"
		if result != expected {
			t.Errorf("got %s, want %s", result, expected)
		}
	})

	t.Run("bool with size", func(t *testing.T) {
		result, err := encoding.ToHex(true).WithSize(32).Hex()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected := "0x0000000000000000000000000000000000000000000000000000000000000001"
		if result != expected {
			t.Errorf("got %s, want %s", result, expected)
		}
	})

	t.Run("signed negative", func(t *testing.T) {
		result, err := encoding.ToHex(big.NewInt(-1)).WithSize(1).WithSigned().Hex()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "0xff" {
			t.Errorf("got %s, want 0xff", result)
		}
	})
}

// RLP tests

func TestRlpEncodeBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"empty", []byte{}, "0x80"},
		{"single byte < 0x80", []byte{0x7f}, "0x7f"},
		{"single byte = 0x80", []byte{0x80}, "0x8180"},
		{"short string", []byte("dog"), "0x83646f67"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.ToRlp(tt.input).Hex()
			if err != nil {
				t.Errorf("RlpEncode error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("RlpEncode(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRlpEncodeList(t *testing.T) {
	tests := []struct {
		name     string
		input    []any
		expected string
	}{
		{"empty list", []any{}, "0xc0"},
		{"list with strings", []any{[]byte("cat"), []byte("dog")}, "0xc88363617483646f67"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoding.ToRlp(tt.input).Hex()
			if err != nil {
				t.Errorf("RlpEncode error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("RlpEncode(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRlpDecode(t *testing.T) {
	t.Run("decode bytes", func(t *testing.T) {
		result, err := encoding.RlpDecode([]byte{0x83, 0x64, 0x6f, 0x67})
		if err != nil {
			t.Errorf("RlpDecode error: %v", err)
		}
		b, ok := result.([]byte)
		if !ok {
			t.Errorf("expected []byte, got %T", result)
		}
		if string(b) != "dog" {
			t.Errorf("got %s, want dog", string(b))
		}
	})

	t.Run("decode list", func(t *testing.T) {
		decoder, err := encoding.FromRlp("0xc88363617483646f67")
		if err != nil {
			t.Errorf("FromRlp error: %v", err)
		}
		result, err := decoder.Bytes()
		if err != nil {
			t.Errorf("Decode error: %v", err)
		}
		list, ok := result.([]any)
		if !ok {
			t.Errorf("expected []any, got %T", result)
		}
		if len(list) != 2 {
			t.Errorf("expected 2 items, got %d", len(list))
		}
	})
}

func TestRlpRoundtrip(t *testing.T) {
	original := []any{[]byte("hello"), []byte("world")}

	encoded, err := encoding.RlpEncode(original)
	if err != nil {
		t.Errorf("RlpEncode error: %v", err)
	}

	decoded, err := encoding.RlpDecode(encoded)
	if err != nil {
		t.Errorf("RlpDecode error: %v", err)
	}

	list, ok := decoded.([]any)
	if !ok {
		t.Errorf("expected []any, got %T", decoded)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
	if string(list[0].([]byte)) != "hello" {
		t.Errorf("first item = %s, want hello", string(list[0].([]byte)))
	}
	if string(list[1].([]byte)) != "world" {
		t.Errorf("second item = %s, want world", string(list[1].([]byte)))
	}
}

// Helper function tests

func TestIsHex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0x1a4", true},
		{"0x", true},
		{"0xdeadbeef", true},
		{"1a4", false},
		{"0xgg", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := encoding.IsHex(tt.input)
			if result != tt.expected {
				t.Errorf("IsHex(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	result := encoding.PadLeft([]byte{0x01, 0xa4}, 4)
	expected := []byte{0x00, 0x00, 0x01, 0xa4}
	if string(result) != string(expected) {
		t.Errorf("PadLeft = %v, want %v", result, expected)
	}
}

func TestPadRight(t *testing.T) {
	result := encoding.PadRight([]byte{0x01, 0xa4}, 4)
	expected := []byte{0x01, 0xa4, 0x00, 0x00}
	if string(result) != string(expected) {
		t.Errorf("PadRight = %v, want %v", result, expected)
	}
}

func TestPadHexLeft(t *testing.T) {
	result := encoding.PadHexLeft("0x1a4", 4)
	if result != "0x000001a4" {
		t.Errorf("PadHexLeft = %s, want 0x000001a4", result)
	}
}

func TestPadHexRight(t *testing.T) {
	result := encoding.PadHexRight("0x1a4", 4)
	if result != "0x1a400000" {
		t.Errorf("PadHexRight = %s, want 0x1a400000", result)
	}
}
